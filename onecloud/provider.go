package onecloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCP_REGION", nil),
				Description: "VPC region to import resources associated with the specific region. 'ua' is used by default ",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCP_API_TOKEN", nil),
				Description: "Service user password",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ocp_flavor":             dataSourceFlavor(),
			"ocp_cluster_version":    dataSourceClusterVersion(),
			"ocp_cluster_networking": dataSourceClusterNetworking(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ocp_cluster":  resourceCluster(),
			"ocp_nodepool": resourceNodePool(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config, diagError := getConfig(d)
	if diagError != nil {
		return nil, diagError
	}

	return config, nil
}
