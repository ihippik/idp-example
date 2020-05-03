package provider

import (
	"github.com/ory/hydra-client-go/client"
)

type Service struct {
	hydra *client.OryHydra
}

// NewService create new service instance.
func NewService(hydra *client.OryHydra) *Service {
	return &Service{hydra: hydra}
}
