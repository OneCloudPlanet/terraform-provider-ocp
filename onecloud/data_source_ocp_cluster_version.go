package onecloud

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-onecloud/internal/ocp_client"
)

func dataSourceClusterVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterVersionRead,
		Schema: map[string]*schema.Schema{
			"versions": {
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
						"images": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"openstack_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_distro": {
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
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"os_distro": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

type versionSearchFilter struct {
	version   string
	imageName string
	osDistro  string
}

func dataSourceClusterVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterVersions, err := client.ClusterVersions(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var clusterVersionsIds []string
	for _, version := range clusterVersions {
		clusterVersionsIds = append(clusterVersionsIds, version.ID)
	}

	filter := getVersionFilterMap(d)
	clusterVersions = filterClusterVersion(clusterVersions, filter)

	clusterObj, err := serializeClusterVersions(clusterVersions)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("versions", clusterObj); err != nil {
		return diag.FromErr(err)
	}
	checksum, err := stringListChecksum(clusterVersionsIds)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(checksum)
	return nil
}

func getVersionFilterMap(d *schema.ResourceData) versionSearchFilter {
	var filterMap map[string]interface{}
	filter := versionSearchFilter{}

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

	imageName, ok := filterMap["image_name"]
	if ok {
		filter.imageName = imageName.(string)
	}

	osDistro, ok := filterMap["os_distro"]
	if ok {
		filter.osDistro = osDistro.(string)
	}

	return filter
}

func filterClusterVersion(versions []ocp_client.ClusterVersion, filter versionSearchFilter) []ocp_client.ClusterVersion {
	var filteredVersions []ocp_client.ClusterVersion

	if filter.version == "" && filter.imageName == "" && filter.osDistro == "" {
		return versions
	}

	for _, version := range versions {
		if filter.version == "" || version.Version == filter.version {
			var filteredImages []ocp_client.Image
			for _, image := range version.Images {
				if (filter.imageName == "" || image.Name == filter.imageName) &&
					(filter.osDistro == "" || image.OsDistro == filter.osDistro) {
					filteredImages = append(filteredImages, image)
				}
			}
			version.Images = filteredImages
			filteredVersions = append(filteredVersions, version)
		}
	}

	return filteredVersions
}

func serializeClusterVersions(versions []ocp_client.ClusterVersion) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	for _, obj := range versions {
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
