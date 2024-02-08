package consul

import (
	"context"
	"errors"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"golang.org/x/exp/rand"
	"movie-app/pkg/discovery"
	"strconv"
	"strings"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(config *consul.Config) (*Registry, error) {
	if config == nil {
		config = consul.DefaultConfig()
	}
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in the form of <host>:<port>/api/<version>, example: localhost:8080/api/v1")
	}

	apiParts := strings.Split(parts[1], "/")
	if len(apiParts) != 3 {
		return errors.New("hostPort must be in the form of <host>:<port>/api/<version>, example: localhost:8080/api/v1")
	}

	port, err := strconv.Atoi(apiParts[0])
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: parts[0],
		Port:    port,
		Check: &consul.AgentServiceCheck{
			CheckID: instanceID,
			TTL:     "5s",
		},
	})
}

func (r *Registry) Deregister(ctx context.Context, instanceID string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}

	var addrs []string
	for _, entry := range entries {
		addrs = append(addrs, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}
	return addrs, nil
}

func (r *Registry) ReportHealthyState(instanceID string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}

func (r *Registry) GetRoundRobinAddress(ctx context.Context, serviceName string) (string, error) {
	addrs, err := r.ServiceAddresses(ctx, serviceName)
	if err != nil {
		return "", err
	} else if len(addrs) == 0 {
		return "", discovery.ErrNotFound
	}
	return addrs[rand.Intn(len(addrs))], nil
}
