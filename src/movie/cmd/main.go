package main

import (
	"context"
	"flag"
	"fmt"
	_consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	v1 "movie-app/movie/internal/v1"
	"movie-app/pkg/discovery"
	"movie-app/pkg/discovery/consul"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	service = "api"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8000, "API handler port")
	flag.Parse()
	log.Println("Starting api service on port", port)

	registry, deregister, err := registerService("v1", port, &_consul.Config{
		Address: "localhost:8500",
	})
	if err != nil {
		log.Fatalln("failed to register service:", err.Error())
	}

	go func() {
		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

		select {
		case <-sigCh:
			deregister()
			os.Exit(0)
		}
	}()

	initApp("grpc", fmt.Sprintf("%d", port), registry)
}

func registerService(version string, port int, config *_consul.Config) (*consul.Registry, func(), error) {
	registry, err := consul.NewRegistry(config)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	svcID := discovery.GenerateInstanceID(service)

	if err := registry.Register(ctx, svcID, service, fmt.Sprintf("localhost:%d/api/%s", port, version)); err != nil {
		return nil, nil, err
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(svcID); err != nil {
				log.Println("failed to report healthy state:", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	deregister := func() {
		if err := registry.Deregister(ctx, svcID); err != nil {
			log.Println("failed to deregister service:", err.Error())
		} else {
			log.Printf("deregistered service %s id %s\n", service, svcID)
		}
	}
	return registry, deregister, nil
}

func initApp(appType, port string, registry discovery.Registry) {
	var err error

	if appType == "gin" {
		err = v1.NewGinApp(port, &registry).Run()
	} else if appType == "grpc" {
		err = v1.NewGRPCApp(
			"localhost:"+port,
			&registry,
			map[string][]grpc.DialOption{
				"metadata": []grpc.DialOption{
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				},
				"rating": []grpc.DialOption{
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				},
			}).Run()
	}

	if err != nil {
		log.Fatalln("failed to start "+appType+" app:", err.Error())
	}
}
