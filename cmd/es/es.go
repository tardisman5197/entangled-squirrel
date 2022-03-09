package main

import (
	"context"
	"es/internal/consts"
	"es/internal/squirrel"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	noHardware := flag.Bool("no-hardware", false, "no hardware connected")
	noUPNP := flag.Bool("no-upnp", false, "do not port forward")
	port := flag.Int("port", 1526, "port to listen for other squirrels on")

	flag.Parse()

	fmt.Printf(
		"hardware: %v\nport: %v\nupnp: %v\n",
		!*noHardware,
		*port,
		!*noUPNP,
	)

	s := squirrel.NewSquirrel(
		"http://entangled-squirrel-0.duckdns.org",
		*port,
		"http://entangled-squirrel-1.duckdns.org:1526",
		!*noHardware,
		!*noUPNP,
	)

	err := s.Setup()
	if err != nil {
		fmt.Printf("could not setup squirrel, got %v\n", err)
		os.Exit(10)
	}

	serverErrors := s.StartServer()

	ctx, cancelPress := context.WithCancel(context.Background())
	pressErrors := s.ListenForPress(ctx)

	ctx, cancelDiscover := context.WithCancel(context.Background())
	discoverErrors := s.DiscoverLoop(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	errorCode := 0

	select {
	case <-sigs:
		fmt.Println("stop signal received")
	case err := <-serverErrors:
		if err != nil {
			fmt.Printf("error while running server, got %v\n", err)
			errorCode = 11
		}
	case err := <-discoverErrors:
		if err != nil {
			fmt.Printf("error while discovering squirrels, got %v\n", err)
			errorCode = 12
		}
	case err := <-pressErrors:
		if err != nil {
			fmt.Printf("error while watching button, got %v\n", err)
			errorCode = 13
		}
	}

	cancelPress()
	cancelDiscover()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*consts.TearDownTimeout,
	)
	err = s.TearDown(ctx)
	if err != nil {
		fmt.Printf("could not tear down squirrel, got %v\n", err)
	}
	cancel()

	fmt.Printf("squirrel stopped, with code %v\n", errorCode)

	os.Exit(errorCode)
}
