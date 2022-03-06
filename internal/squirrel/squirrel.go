package squirrel

import (
	"context"
	"encoding/json"
	"es/internal/consts"
	"es/internal/utils"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/NebulousLabs/go-upnp"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type Squirrel struct {
	// entangledSquirrelURLs contains a list of squirrels
	// which can be communicated with.
	entangledSquirrelURLs map[string]int
	entangledSquirrelLock sync.RWMutex

	port   uint16
	url    string
	server *http.Server
	igd    *upnp.IGD

	button gpio.PinIO

	lightsLock sync.Mutex
	lights     gpio.PinIO
}

func NewSquirrel(url string, port uint16, knownSquirrelURL string) *Squirrel {
	return &Squirrel{
		url:                   url,
		port:                  port,
		entangledSquirrelURLs: map[string]int{knownSquirrelURL: 0},
	}
}

func (s *Squirrel) Setup(ctx context.Context) error {
	// Setup buttons and lights
	_, err := host.Init()
	if err != nil {
		return fmt.Errorf("could not setup hardware, got %v", err)
	}

	button := gpioreg.ByName(consts.ButtonName)
	if button == nil {
		return fmt.Errorf("could not find button")
	}
	err = button.In(gpio.PullDown, gpio.BothEdges)
	if err != nil {
		return fmt.Errorf("could not init button, got %v", err)
	}
	s.button = button

	lights := gpioreg.ByName(consts.LightsName)
	if button == nil {
		return fmt.Errorf("could not find lights")
	}
	err = lights.Out(gpio.Low)
	if err != nil {
		return fmt.Errorf("could not init lights, got %v", err)
	}
	s.lightsLock.Lock()
	s.lights = lights
	s.lightsLock.Unlock()

	// Setup networking
	igd, err := upnp.DiscoverCtx(ctx)
	if err != nil {
		return fmt.Errorf("could not discover router, got %v", err)
	}
	s.igd = igd

	err = s.igd.Forward(s.port, "an entangled squirrel")
	if err != nil {
		return fmt.Errorf("could not forward a port on the router, got %v", err)
	}

	// Setup server
	r := mux.NewRouter()
	r.HandleFunc("/update", s.flashReq).Methods("GET")
	r.HandleFunc("/squirrels/add", s.addSquirrelsReq).Methods("POST")
	r.HandleFunc("/squirrels/known", s.knownSquirrelsReq).Methods("GET")

	s.server = &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%v", s.port),
	}

	return nil
}

func (s *Squirrel) TearDown(ctx context.Context) error {
	s.server.Shutdown(ctx)

	err := s.igd.Clear(s.port)
	if err != nil {
		return fmt.Errorf("could not remove port forwarding, got %v", err)
	}
	return nil
}

func (s *Squirrel) StartServer() chan error {
	errors := make(chan error, 1)
	go func() {
		errors <- s.server.ListenAndServe()
	}()
	return errors
}

func (s *Squirrel) ListenForPress(ctx context.Context) chan error {
	errors := make(chan error, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				errors <- nil
				return
			default:
			}
			press := s.button.WaitForEdge(time.Millisecond * consts.ButtonTimout)
			if press {
				s.handlePress()
			}
		}
	}()
	return errors
}

func (s *Squirrel) DiscoverLoop(ctx context.Context) chan error {
	errors := make(chan error, 1)
	go func() {
		ticker := time.NewTicker(time.Millisecond * consts.DiscoverInterval)
		for {
			select {
			case <-ctx.Done():
				errors <- nil
				return
			case <-ticker.C:
				s.discoverSquirrels()
			}
		}
	}()
	return errors
}

func (s *Squirrel) handlePress() {
	s.entangledSquirrelLock.RLock()
	urls := utils.GetKeys(s.entangledSquirrelURLs)
	s.entangledSquirrelLock.RUnlock()

	for _, url := range urls {
		err := utils.SendFlash(url)
		if err != nil {
			s.entangledSquirrelLock.Lock()
			s.entangledSquirrelURLs[url] = s.entangledSquirrelURLs[url] + 1
			fmt.Printf("Could not send flash, got %v\n", err)
			if s.entangledSquirrelURLs[url] > consts.SquirrelErrorLimit {
				delete(s.entangledSquirrelURLs, url)
				fmt.Printf("Removed Squirrel from entangled list, too many failed connections\n")
			}
			s.entangledSquirrelLock.Unlock()
		}
	}
}

func (s *Squirrel) discoverSquirrels() {
	s.entangledSquirrelLock.RLock()
	urls := utils.GetKeys(s.entangledSquirrelURLs)
	s.entangledSquirrelLock.RUnlock()

	for _, url := range urls {
		resp, err := http.Get(fmt.Sprintf("%s/squirrels/known", url))
		if err != nil {
			fmt.Printf("Could not get list of known squirrels, got %v\n", err)
		}

		var squirrels []string
		err = json.NewDecoder(resp.Body).Decode(&squirrels)
		if err != nil {
			fmt.Printf("Could not read list of known squirrels, got %v\n", err)
		}

		s.addSquirrels(squirrels)
	}
}

func (s *Squirrel) addSquirrels(urls []string) {
	s.entangledSquirrelLock.Lock()
	for _, url := range urls {
		if _, found := s.entangledSquirrelURLs[url]; !found {
			if url != s.url {
				s.entangledSquirrelURLs[url] = 0
			}
		}
	}
	s.entangledSquirrelLock.Unlock()
}

func (s *Squirrel) flash() {
	s.lightsLock.Lock()
	for i := 0; i < consts.NoOfFlashes; i++ {
		s.lights.Out(gpio.High)
		time.Sleep(time.Millisecond * consts.FlashInterval)
	}
	s.lightsLock.Unlock()
}
