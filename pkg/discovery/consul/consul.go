package consul

import (
	"context"
	"errors"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"movie-app/pkg/discovery"
	"strconv"
	"strings"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(config consul.Config) (*Registry, error) {
	client, err := consul.NewClient(&config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in the form of <host>:<port>, example: localhost:8080")
	}

	port, err := strconv.Atoi(parts[1])
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
