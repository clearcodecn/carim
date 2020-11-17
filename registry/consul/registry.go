package consul

import (
	"context"
	"fmt"
	reg "github.com/clearcodecn/carim/registry"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"
)

type registry struct {
	client  *api.Client
	options Options

	registration *api.AgentServiceRegistration
	stopChan     chan error
}

func (r *registry) Init(ctx context.Context) error {
	var options = DefaultOption
	opts := OptionFromContext(ctx)
	for _, o := range opts {
		o(&options)
	}

	address, err := url.Parse(options.Address)
	if err != nil {
		return errors.Wrap(err, "failed to parse address")
	}

	config := api.DefaultConfig()
	if address != nil {
		config.Address = address.Host
		config.Scheme = address.Scheme
	}

	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	r.client = client
	r.options = options

	return nil
}

func (r *registry) Register(service *reg.Service) error {
	if err := r.doRegister(service); err != nil {
		return err
	}
	if r.options.TTL == 0 {
		return nil
	}
	go func() {
		for {
			after := time.After(r.options.TTL)
			select {
			case <-r.stopChan:
				return
			case <-after:
				if err := r.passTTL(); err != nil {
					log.Println("[ERROR] failed register service", err)
				}
			}
		}
	}()
	return nil
}

func (r *registry) doRegister(service *reg.Service) error {
	var registration *api.AgentServiceRegistration
	if r.registration != nil {
		registration = r.registration
	} else {

		host, port, err := net.SplitHostPort(service.Endpoint)
		if err != nil {
			return err
		}
		if port == "" {
			return errors.New("invalid port")
		}
		if host == "" {
			return errors.New("invalid ip")
		}
		pi, err := strconv.Atoi(port)
		if err != nil {
			return err
		}

		registration = &api.AgentServiceRegistration{
			ID:      service.Id,
			Name:    service.Name,
			Address: host,
			Port:    pi,
			Meta: map[string]string{
				"version": service.Version,
			},
		}
		var check api.AgentServiceCheck
		if r.options.TTL != 0 {
			check = api.AgentServiceCheck{
				CheckID:                        fmt.Sprintf("service:%s:%s", service.Name, service.Id),
				Name:                           fmt.Sprintf("check for service: %s:%s", service.Name, service.Id),
				TTL:                            r.options.TTL.String(),
				TLSSkipVerify:                  true,
				DeregisterCriticalServiceAfter: "1m",
			}
		}
		registration.Check = &check
	}
	if err := r.client.Agent().ServiceRegister(registration); err != nil {
		return err
	}
	r.registration = registration
	return nil
}

func (r *registry) passTTL() error {
	if r.registration == nil || r.registration.Check == nil {
		return errors.New("registration or check is nil")
	}
	if err := r.client.Agent().PassTTL(r.registration.Check.CheckID, "ok"); err != nil {
		return err
	}
	return nil
}

func (r *registry) DeRegister(service *reg.Service) error {
	return r.client.Agent().ServiceDeregister(service.Id)
}

func (r *registry) ServiceByName(name string) (*reg.Service, error) {
	services, err := r.client.Agent().Services()
	if err != nil {
		return nil, err
	}
	for _, v := range services {
		if v.Service == name {
			return &reg.Service{
				Id:       v.ID,
				Name:     name,
				Endpoint: v.Address,
				Version:  v.Meta["version"],
			}, nil
		}
	}

	return nil, errors.New("not found")
}

func (r *registry) ServiceList() ([]*reg.Service, error) {

	services, err := r.client.Agent().Services()
	if err != nil {
		return nil, err
	}

	var regServices []*reg.Service
	for _, v := range services {
		regServices = append(regServices, &reg.Service{
			Id:       v.ID,
			Name:     v.Service,
			Endpoint: v.Address,
			Version:  v.Meta["version"],
		})
	}

	return regServices, nil
}
