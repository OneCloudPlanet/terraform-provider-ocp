package onecloud

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-onecloud/internal/ocp_client"
	"sort"
	"strings"
)

func dataSourceFlavor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlavorRead,
		Schema: map[string]*schema.Schema{
			"flavors": {
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
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_mb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_gb": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"root_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"assigned_cluster_templates": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ephemeral_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"flavor_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"out_of_stock": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reseller_resources": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"swap": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used_by_resellers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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
						"vcpus": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"memory_mb": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"memory_gb": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"root_gb": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

type flavorSearchFilter struct {
	vcpus    int
	memoryGb float64
	memoryMb int
	rootGb   int
}

func dataSourceFlavorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	flavors, err := client.Flavors(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	fmt.Println(flavors)
	flavorsIds := []string{}
	for _, flavor := range flavors {
		flavorsIds = append(flavorsIds, flavor.ID)
	}

	filter := getFlavorFilterMap(d)
	flavors = filterFlavor(flavors, filter)

	flavorsObj, err := serializeFlavors(flavors)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("flavors", flavorsObj); err != nil {
		return diag.FromErr(err)
	}
	checksum, err := stringListChecksum(flavorsIds)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(checksum)
	return nil
}

func getFlavorFilterMap(d *schema.ResourceData) flavorSearchFilter {
	var filterMap map[string]interface{}
	filter := flavorSearchFilter{}

	if filterSet, ok := d.GetOk("filter"); ok {
		if filterSet.(*schema.Set).Len() == 0 {
			return filter
		}
		filterMap = filterSet.(*schema.Set).List()[0].(map[string]interface{})
	}

	vcpus, ok := filterMap["vcpus"]
	if ok {
		filter.vcpus = vcpus.(int)
	}

	memoryGb, ok := filterMap["memory_gb"]
	if ok {
		filter.memoryGb = memoryGb.(float64)
	}

	memoryMb, ok := filterMap["memory_mb"]
	if ok {
		filter.memoryMb = memoryMb.(int)
	}

	rootGb, ok := filterMap["root_gb"]
	if ok {
		filter.rootGb = rootGb.(int)
	}

	return filter
}

func filterFlavor(flavors []ocp_client.Flavor, filter flavorSearchFilter) []ocp_client.Flavor {
	var filteredFlavors []ocp_client.Flavor

	if filter.vcpus == 0 && filter.rootGb == 0 && filter.memoryMb == 0 && filter.memoryGb == 0 {
		return flavors
	}

	for _, f := range flavors {
		if (filter.vcpus == 0 || f.Vcpus == filter.vcpus) &&
			(filter.memoryGb == 0 || f.MemoryGb == filter.memoryGb) &&
			(filter.memoryMb == 0 || f.MemoryMb == filter.memoryMb) &&
			(filter.rootGb == 0 || f.RootGb == filter.rootGb) {
			filteredFlavors = append(filteredFlavors, f)
		}
	}

	return filteredFlavors
}

func serializeFlavors(flavors []ocp_client.Flavor) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	for _, obj := range flavors {
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

func stringChecksum(s string) (string, error) {
	h := md5.New() // #nosec
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs), nil
}

func stringListChecksum(s []string) (string, error) {
	sort.Strings(s)
	checksum, err := stringChecksum(strings.Join(s, ""))
	if err != nil {
		return "", err
	}

	return checksum, nil
}
