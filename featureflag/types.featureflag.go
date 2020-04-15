package featureflag

import (
	"net/http"
	"sync"

	consulapi "github.com/hashicorp/consul/api"
	c "github.com/robfig/cron"
)

type (

	// Instance type
	Instance int

	// Status type
	Status int

	// Config struct of the feature flag
	Config struct {
		ConsulURL    string
		CronInterval int
		HTTPClient   *http.Client
		Instances    []Instance
		Prefix       string
	}

	// main module
	flagModule struct {
		config Config
		kv     *consulapi.KV
		flags  sync.Map
		cron   *c.Cron
	}

	// private flag struct
	flag struct {
		Status Status
		Users  []int64
		Value  map[int]string
	}

	// FeatureFlag interface provide function for feature flagging
	FeatureFlag interface {
		HasAccess(flagKey string, instance Instance, userID int64, isQA bool) bool
		GetActiveValue(flagKey string, instance Instance, userID int64, isQA bool) (bool, map[int]string)
		HasAccessInPercentageUser(flagKey string, instance Instance, userID int64, isQA bool) bool
	}
)
