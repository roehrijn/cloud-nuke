package externalcreds

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var configsByRegion = map[string]aws.Config{}
var lock sync.Mutex

func Get(region string) (aws.Config, error) {
	lock.Lock()
	defer lock.Unlock()

	cfg, found := configsByRegion[region]
	if !found {
		var err error
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
		if err != nil {
			return aws.Config{}, err
		}
		configsByRegion[region] = cfg
	}

	return cfg, nil
}

func Set(region string, cfg aws.Config) {
	lock.Lock()
	defer lock.Unlock()
	configsByRegion[region] = cfg
}
