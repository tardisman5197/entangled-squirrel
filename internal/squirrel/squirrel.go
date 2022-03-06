package squirrel

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/NebulousLabs/go-upnp"
)

type Squirrel struct {
	// entangledSquirrelURLs contains a list of squirrels
	// which can be communicated with.
	entangledSquirrelURLs map[string]bool

	port uint16

	server *http.Server
	igd    *upnp.IGD
}

func NewSquirrel(connectedSquirrelURL string) *Squirrel {
	return &Squirrel{
		entangledSquirrelURLs: map[string]bool{connectedSquirrelURL: true},
	}
}

func (s *Squirrel) Setup(ctx context.Context) error {
	igd, err := upnp.DiscoverCtx(ctx)
	if err != nil {
		return fmt.Errorf("could not discover router, got %v", err)
	}
	s.igd = igd

	err = s.igd.Forward(s.port, "an entangled squirrel")
	if err != nil {
		return fmt.Errorf("could not forward a port on the router, got %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/update", s.update).Methods("POST")
	r.HandleFunc("/squirrels/add", s.addSquirrel).Methods("POST")
	r.HandleFunc("/squirrels/known", s.knownSquirrels).Methods("GET")

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

func (s *Squirrel) Start() {
	go func() {

	}()

}
