package consul

import (
	"errors"
	"fmt"
	"strings"
	"validation_service/pkg/config"
	"validation_service/pkg/log"

	"github.com/hashicorp/consul/api"
)

type consulStorage struct {
	client *api.Client
	kv     *api.KV
}

var Storage *consulStorage

func Init() {
	var err error
	Storage = new(consulStorage)

	Storage.client, err = api.NewClient(&api.Config{
		Address: config.Settings.ConsulInfo.Address,
		Token:   config.Settings.ConsulInfo.Token,
	})
	if err != nil {
		log.Logger.Fatalf("Could not connect to consul: %s", err)
	}

	Storage.kv = Storage.client.KV()
}

func (c *consulStorage) Get(object string) ([]byte, error) {
	pair, _, err := c.kv.Get(fmt.Sprintf("VALIDATIONS/%s/%s", strings.ToUpper(config.Settings.Env), strings.ToUpper(object)), nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return []byte{}, errors.New("not found")
	}

	return pair.Value, err
}
