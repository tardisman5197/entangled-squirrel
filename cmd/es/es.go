package main

import (
	"context"
	"es/internal/consts"
	"es/internal/squirrel"
	"fmt"
	"os"
	"time"
)

func main() {
	s := squirrel.NewSquirrel(
		"http://entangled-squirrel-0.duckdns.org",
		1526,
		"http://entangled-squirrel-1.duckdns.org:1526",
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*consts.SetupTimout,
	)
	err := s.Setup(ctx)
	cancel()
	if err != nil {
		fmt.Printf("could not setup squirrel, got %v\n", err)
		os.Exit(10)
	}

	serverErrors := s.StartServer()

	ctx, cancelPress := context.WithCancel(context.Background())
	pressErrors := s.ListenForPress(ctx)

	ctx, cancelDiscover := context.WithCancel(context.Background())
	discoverErrors := s.DiscoverLoop(ctx)

	errorCode := 0

	select {
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

	ctx, cancel = context.WithTimeout(
		context.Background(),
		time.Millisecond*consts.TearDownTimeout,
	)
	err = s.TearDown(ctx)
	if err != nil {
		fmt.Printf("Could not tear down squirrel, got %v\n", err)
	}
	cancel()

	os.Exit(errorCode)
}
