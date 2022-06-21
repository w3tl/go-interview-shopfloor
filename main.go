package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"shopfloor/client"
	"shopfloor/resource"
	"shopfloor/resourcehub"
	"time"

	"syscall"
)

func main() {
	end, ctx := gracefulShutdown()

	childCtx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()

	hub := resourcehub.NewServer()

	resources := []*resource.Resource{
		resource.New("machine_1"),
		resource.New("machine_2"),
		resource.New("machine_3"),
	}

	for _, res := range resources {
		err := hub.AddResource(res)
		if err != nil {
			fmt.Printf("Couldn't add a resource %s: %v\n", res.Name(), err)
		}
	}

	srv := hub.Listen(3030)

	go func() {
		err := client.Run(childCtx, "client 1", "http://localhost:3030", "machine_1")
		if err != nil {
			fmt.Printf("Client 1 faced with an error: %v", err)
		}
	}()

	go func() {
		err := client.Run(childCtx, "client 2", "http://localhost:3030", "machine_1")
		if err != nil {
			fmt.Printf("Client 2 faced with an error: %v", err)
		}
	}()

	select {
	case <-end:
	case <-childCtx.Done():
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}

func gracefulShutdown() (<-chan struct{}, context.Context) {
	end := make(chan struct{})
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		fmt.Print("Shutting down gracefully")
		cancel()
		close(end)
	}()
	return end, ctx
}
