package onecloud

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-onecloud/internal/ocp_client"
)

func dataSourceClusterAddons() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterAddonsRead,
		Schema: map[string]*schema.Schema{
			"addons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"releases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
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
					},
				},
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

type addonSearchFilter struct {
	name string
}

func dataSourceClusterAddonsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterAddons, err := client.ClusterAddons(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var clusterAddonsIds []string
	for _, version := range clusterAddons {
		clusterAddonsIds = append(clusterAddonsIds, version.ID)
	}

	filter := getAddonFilterMap(d)
	clusterAddons = filterClusterAddons(clusterAddons, filter)

	addonsObj, err := serializeClusterAddons(clusterAddons)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("addons", addonsObj); err != nil {
		return diag.FromErr(err)
	}
	checksum, err := stringListChecksum(clusterAddonsIds)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(checksum)
	return nil
}

func getAddonFilterMap(d *schema.ResourceData) addonSearchFilter {
	var filterMap map[string]interface{}
	filter := addonSearchFilter{}

	if filterSet, ok := d.GetOk("filter"); ok {
		if filterSet.(*schema.Set).Len() == 0 {
			return filter
		}
		filterMap = filterSet.(*schema.Set).List()[0].(map[string]interface{})
	}

	name, ok := filterMap["name"]
	if ok {
		filter.name = name.(string)
	}

	return filter
}

func filterClusterAddons(addons []ocp_client.ClusterAddon, filter addonSearchFilter) []ocp_client.ClusterAddon {
	var filteredAddons []ocp_client.ClusterAddon

	if filter.name == "" {
		return addons
	}

	for _, addon := range addons {
		if addon.Name == filter.name {
			filteredAddons = append(filteredAddons, addon)
		}
	}
	return filteredAddons
}

func serializeClusterAddons(addons []ocp_client.ClusterAddon) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	for _, obj := range addons {
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
