package featureflag

import (
	"errors"

	consulapi "github.com/hashicorp/consul/api"
)

func TFuncPatch() {
	newConsulClient = func(config *consulapi.Config) (*consulapi.Client, error) {
		if config.Address == "error" {
			return nil, errors.New("foo")
		}

		return consulapi.NewClient(config)
	}
}
