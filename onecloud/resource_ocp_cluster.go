package onecloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"strings"
	"time"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOCPClusterCreate,
		ReadContext:   resourceOCPClusterRead,
		UpdateContext: resourceOCPClusterUpdate,
		DeleteContext: resourceOCPClusterDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"master_flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"master_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"networking": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"restriction_api": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"restriction_ips": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"node_pool": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
								return strings.EqualFold(old, new)
							},
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_count": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},
						"autoscale": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: false,
						},
						"max_count": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: false,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nodes": {
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
									"ready": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flavor": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_pool": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"control_plane": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"api_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_nodes": {
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
						"ready": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_pool": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_plane": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOCPClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	cluster, err := client.GetCluster(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	fetchErr := fetchClusterState(cluster, d)
	if fetchErr != nil {
		return fetchErr
	}

	return nil
}

func resourceOCPClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterData := GetClusterCreateOptions(d)

	resp, err := client.CreateCluster(ctx, clusterData)
	if err != nil {
		return diag.FromErr(err)
	}
	operationId := resp["operation_id"].(string)

	result, waitErr := waitForOperationSuccess(ctx, *client, operationId, d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		return diag.FromErr(waitErr)
	}

	clusterId := result.(map[string]interface{})["primary_object_id"].(string)
	d.SetId(clusterId)

	cluster, err := client.GetCluster(ctx, clusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	fetchErr := fetchClusterState(cluster, d)
	if fetchErr != nil {
		return fetchErr
	}

	return nil
}

func resourceOCPClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	var clusterUpdate bool = false
	clusterUpdateData := make(map[string]interface{})
	if d.HasChange("cluster_version") {
		clusterUpdate = true
		clusterUpdateData["cluster_version"] = d.Get("cluster_version")
	}
	if d.HasChange("restriction_api") || d.HasChange("restriction_ips") {
		clusterUpdate = true
		clusterUpdateData["restriction_api"] = d.Get("restriction_api")
		clusterUpdateData["restriction_ips"] = d.Get("restriction_ips")
	}
	if clusterUpdate {
		resp, err := client.UpdateCluster(ctx, d.Id(), clusterUpdateData)
		if err != nil {
			return diag.FromErr(err)
		}
		for field, value := range resp {
			err := d.Set(field, value)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("node_pool") {
		oldPool, newPool := d.GetChange("node_pool")
		oldSet, newSet := oldPool.([]interface{}), newPool.([]interface{})

		for _, newNp := range newSet {
			var oldNodePool map[string]interface{}
			newNodePool := newNp.(map[string]interface{})

			for _, oldVal := range oldSet {
				if oldVal.(map[string]interface{})["id"] == newNodePool["id"] {
					oldNodePool = oldVal.(map[string]interface{})
					break
				}
			}

			if !reflect.DeepEqual(oldNodePool["node_count"], newNodePool["node_count"]) ||
				!reflect.DeepEqual(oldNodePool["autoscale"], newNodePool["autoscale"]) ||
				!reflect.DeepEqual(oldNodePool["max_count"], newNodePool["max_count"]) {

				data := map[string]interface{}{
					"count":     newNodePool["node_count"],
					"autoscale": newNodePool["autoscale"],
					"max_count": newNodePool["max_count"],
				}
				resp, err := client.UpdateNodePool(ctx, oldNodePool["id"].(string), data)
				if err != nil {
					return diag.FromErr(err)
				}

				newNodePool["node_count"] = resp["count"]
				newNodePool["autoscale"] = resp["autoscale"]
				newNodePool["max_count"] = resp["max_count"]

				err = d.Set("node_pool", []map[string]interface{}{newNodePool})
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}
	return nil
}

func resourceOCPClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.DeleteCluster(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	operationId := resp["operation_id"].(string)
	_, waitErr := waitForOperationSuccess(ctx, *client, operationId, d.Timeout(schema.TimeoutDelete))
	if waitErr != nil {
		return diag.FromErr(waitErr)
	}

	return nil
}

func fetchClusterState(cluster map[string]interface{}, d *schema.ResourceData) diag.Diagnostics {
	nodePoolsMap := make(map[string]map[string]interface{})
	clusterNodePools := cluster["node_pools"].([]interface{})
	nodePools := d.Get("node_pool").([]interface{})[0].(map[string]interface{})

	for _, np := range clusterNodePools {
		nodePoolsMap[np.(map[string]interface{})["name"].(string)] = np.(map[string]interface{})
	}

	mappedNp := nodePoolsMap[nodePools["name"].(string)]
	nodePools["id"] = mappedNp["id"].(string)
	nodePools["flavor"] = mappedNp["flavor"].(string)
	nodePools["is_default"] = mappedNp["is_default"].(bool)
	nodePools["status"] = mappedNp["status"].(string)
	nodePools["nodes"] = mappedNp["nodes"].([]interface{})

	err := d.Set("cluster_name", cluster["cluster_name"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("cluster_version", cluster["cluster_version"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("restriction_api", cluster["restriction_api"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("restriction_ips", cluster["restriction_ips"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("node_pool", []map[string]interface{}{nodePools})
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("api_address", cluster["api_address"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("control_nodes", cluster["control_nodes"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("status", cluster["status"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("status_reason", cluster["status_reason"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("created_at", cluster["created_at"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("updated_at", cluster["updated_at"])
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
