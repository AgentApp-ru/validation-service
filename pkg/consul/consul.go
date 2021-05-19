package consul

import (
	"encoding/json"
	"fmt"
	"strings"
	"validation_service/pkg/config"
	"validation_service/pkg/log"

	"github.com/hashicorp/consul/api"
)

type client struct {
	client *api.Client
	kv     *api.KV
}

var c *client

func Init() {
	var err error
	c = new(client)

	c.client, err = api.NewClient(&api.Config{
		Address: config.Settings.ConsulInfo.Address,
		Token:   config.Settings.ConsulInfo.Token,
	})
	if err != nil {
		log.Logger.Fatalf("Could not connect to consul: %s", err)
	}

	c.kv = c.client.KV()
}

func Get(object string) (interface{}, error) {
	var data interface{}

	pair, _, err := c.kv.Get(fmt.Sprintf("VALIDATIONS/%s", strings.ToUpper(object)), nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(pair.Value, &data)

	return data, err
}
