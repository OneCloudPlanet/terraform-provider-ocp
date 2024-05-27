package onecloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"time"
)

func resourceNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOCPNodePoolCreate,
		ReadContext:   resourceOCPNodePoolRead,
		UpdateContext: resourceOCPNodePoolUpdate,
		DeleteContext: resourceOCPNodePoolDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MinItems: 0,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MinItems: 0,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
	}
}

func resourceOCPNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.GetNodePool(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	errFetch := fetchNodePoolState(res, d)
	if errFetch != nil {
		return errFetch
	}
	return nil
}

func resourceOCPNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	nodePoolData := GetNodePoolCreateOptions(d)
	res, err := client.CreateNodePool(ctx, nodePoolData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(res["id"].(string))

	return nil
}

func resourceOCPNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("node_count") || d.HasChange("autoscale") || d.HasChange("max_count") {
		updateData := make(map[string]interface{})
		updateData["count"] = d.Get("node_count")
		updateData["autoscale"] = d.Get("autoscale")
		updateData["max_count"] = d.Get("max_count")

		res, err := client.UpdateNodePool(ctx, d.Id(), updateData)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set("node_count", res["count"])
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("autoscale", res["autoscale"])
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("max_count", res["max_count"])
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceOCPNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getOCPClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.DeleteNodePool(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func fetchNodePoolState(nodePool map[string]interface{}, d *schema.ResourceData) diag.Diagnostics {
	err := d.Set("name", nodePool["name"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("flavor", nodePool["flavor"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("node_count", nodePool["count"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("autoscale", nodePool["autoscale"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("max_count", nodePool["max_count"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("is_default", nodePool["is_default"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("status", nodePool["status"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("labels", nodePool["labels"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("taints", nodePool["taints"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("nodes", nodePool["nodes"])
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
