package onecloud

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-onecloud/internal/ocp_client"
)

func dataSourceClusterNetworking() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingRead,
		Schema: map[string]*schema.Schema{
			"networking": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

type networkingSearchFilter struct {
	name    string
	version string
}

func dataSourceNetworkingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	networking, err := client.Networking(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	networkingIds := []string{}
	for _, n := range networking {
		networkingIds = append(networkingIds, n.ID)
	}

	filter := getNetworkingFilterMap(d)
	networking = filterNetworking(networking, filter)

	networkingObj, err := serializeNetworking(networking)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("networking", networkingObj); err != nil {
		return diag.FromErr(err)
	}
	checksum, err := stringListChecksum(networkingIds)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(checksum)
	return nil
}

func getNetworkingFilterMap(d *schema.ResourceData) networkingSearchFilter {
	var filterMap map[string]interface{}
	filter := networkingSearchFilter{}

	if filterSet, ok := d.GetOk("filter"); ok {
		if filterSet.(*schema.Set).Len() == 0 {
			return filter
		}
		filterMap = filterSet.(*schema.Set).List()[0].(map[string]interface{})
	}

	version, ok := filterMap["version"]
	if ok {
		filter.version = version.(string)
	}

	name, ok := filterMap["network_name"]
	if ok {
		filter.name = name.(string)
	}

	return filter
}

func filterNetworking(networking []ocp_client.Networking, filter networkingSearchFilter) []ocp_client.Networking {
	var filteredNetworking []ocp_client.Networking

	if filter.version == "" && filter.name == "" {
		return networking
	}

	for _, n := range networking {
		if (filter.version == "" || n.Version == filter.version) &&
			(filter.name == "" || n.Name == filter.name) {
			filteredNetworking = append(filteredNetworking, n)
		}
	}

	return filteredNetworking
}

func serializeNetworking(networking []ocp_client.Networking) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	for _, obj := range networking {
		jsonData, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		var sObj map[string]interface{}
		err = json.Unmarshal(jsonData, &sObj)
		if err != nil {
			return nil, err
		}
		result = append(result, sObj)
	}
	return result, nil
}
