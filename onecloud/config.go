package onecloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-onecloud/internal/ocp_client"
	"sync"
)

var (
	cfgSingletone *Config
	once          sync.Once
)

// Config contains all available configuration options.
type Config struct {
	ApiToken string
	Region   string
	Context  context.Context
	lock     sync.Mutex
}

func getConfig(d *schema.ResourceData) (*Config, diag.Diagnostics) {
	once.Do(func() {
		cfgSingletone = &Config{
			ApiToken: d.Get("api_token").(string),
		}
		if v, ok := d.GetOk("region"); ok {
			cfgSingletone.Region = v.(string)
		} else {
			cfgSingletone.Region = ocp_client.UARegion
		}
	})

	return cfgSingletone, nil
}

func getOCPClient(meta interface{}) (*ocp_client.Client, error) {
	config := meta.(*Config)
	client, err := ocp_client.GetClient(config.ApiToken, config.Region)
	if err != nil {
		return nil, err
	}
	return client, nil
}
