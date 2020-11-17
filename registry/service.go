package registry

import (
	"context"
	"fmt"
)

type Service struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Version  string `json:"version"`
}

func (s *Service) String() string {
	return fmt.Sprintf("name=%s endpoint=%s version=`%s`", s.Name, s.Endpoint, s.Version)
}

type Registry interface {
	Init(ctx context.Context)
	Register(service *Service) error
	DeRegister(service *Service) error
	ServiceByName(name string) (*Service, error)
	ServiceList() ([]*Service, error)
}
