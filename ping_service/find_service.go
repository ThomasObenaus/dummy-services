package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	consul "github.com/hashicorp/consul/api"
)

//Client provides an interface for getting data out of Consul
type Client interface {
	// Get a Service from consul
	Service(string, string) ([]*consul.ServiceEntry, *consul.QueryMeta, error)
	// Register a service with local agent
	Register(string, string, int) error
	// Deregister a service with local agent
	DeRegister(string) error

	FindProvider(string) (string, error)
}

type client struct {
	consul *consul.Client
}

//NewConsul returns a Client interface for given consul address
func NewConsulClient(addr string) (Client, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	config := consul.DefaultConfig()
	config.Address = addr
	c, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &client{consul: c}, nil
}

// Register a service with consul local agent
func (c *client) Register(name string, address string, port int) error {
	rand.Seed(time.Now().UTC().UnixNano())

	reg := &consul.AgentServiceRegistration{
		ID:      name,
		Name:    name,
		Port:    port,
		Address: address,
	}
	return c.consul.Agent().ServiceRegister(reg)
}

// DeRegister a service with consul local agent
func (c *client) DeRegister(id string) error {
	return c.consul.Agent().ServiceDeregister(id)
}

// Service return a service
func (c *client) Service(service, tag string) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	passingOnly := true
	addrs, meta, err := c.consul.Health().Service(service, tag, passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return nil, nil, fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return nil, nil, err
	}
	return addrs, meta, nil
}

func (c *client) FindProvider(serviceName string) (string, error) {

	serviceEntries, _, err := c.Service(serviceName, "")
	if err != nil {
		return "", err
	}

	numServiceEntries := len(serviceEntries)
	selectedIdx := rand.Intn(numServiceEntries)

	serviceDef := *serviceEntries[selectedIdx].Service
	addr := serviceDef.Address + ":" + strconv.Itoa(serviceDef.Port)
	// log.Printf("%s %s\n", serviceDef.Service, addr)
	return addr, nil
}
