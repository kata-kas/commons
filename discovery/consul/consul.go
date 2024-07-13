package consul

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	consul "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(addr, serviceName string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Registry{client}, nil
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid host:port format. Eg: localhost:8080")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	// host := parts[0]

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: serviceName,
		Port:    port,
		Check: &consul.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func (r *Registry) Deregister(ctx context.Context, instanceID, serverName string) error {
	log.Printf("Deregistering service %s with ID %s", serverName, instanceID)
	return r.client.Agent().CheckDeregister(instanceID)
}

func (r *Registry) HealthCheck(instanceID, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	var hosts []string
	for _, entry := range entries {
		hosts = append(hosts, entry.Service.Address+":"+strconv.Itoa(entry.Service.Port))
	}

	return hosts, nil
}
