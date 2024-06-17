package onecloud

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ClusterCreateOptions struct {
	ClusterName    string                  `json:"cluster_name"`
	ClusterVersion string                  `json:"cluster_version"`
	MasterFlavorId string                  `json:"master_flavor_id"`
	MasterCount    int                     `json:"master_count"`
	NodePool       []NodePoolCreateOptions `json:"node_pools"`
	Image          string                  `json:"image"`
	Networking     string                  `json:"networking"`
	RestrictionApi bool                    `json:"restriction_api"`
	RestrictionIps []string                `json:"restriction_ips"`
	Addons         []Addon                 `json:"addons"`
}

type NodePoolCreateOptions struct {
	Name      string  `json:"name"`
	FlavorId  string  `json:"flavor_id"`
	Count     int     `json:"count"`
	Autoscale bool    `json:"autoscale"`
	MaxCount  int     `json:"max_count"`
	IsDefault bool    `json:"is_default"`
	Labels    []Label `json:"labels"`
	Taints    []Taint `json:"taints"`
	Cluster   string  `json:"cluster"`
}

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Taint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"`
}

type Addon struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *ClusterCreateOptions) toJson() ([]byte, diag.Diagnostics) {
	res, err := json.Marshal(c)
	if err != nil {
		return nil, diag.Errorf("Can't convert to json")
	}

	return res, nil
}

func GetClusterCreateOptions(d *schema.ResourceData) *ClusterCreateOptions {
	return &ClusterCreateOptions{
		ClusterName:    d.Get("cluster_name").(string),
		ClusterVersion: d.Get("cluster_version").(string),
		MasterFlavorId: d.Get("master_flavor_id").(string),
		MasterCount:    d.Get("master_count").(int),
		Image:          d.Get("image").(string),
		Networking:     d.Get("networking").(string),
		RestrictionApi: d.Get("restriction_api").(bool),
		RestrictionIps: getListOfString(d.Get("restriction_ips").([]interface{})),
		NodePool:       getNodePool(d.Get("node_pool").([]interface{})),
		Addons:         getAddons(d.Get("addons").([]interface{})),
	}
}

func GetNodePoolCreateOptions(d *schema.ResourceData) *NodePoolCreateOptions {
	return &NodePoolCreateOptions{
		Name:      d.Get("name").(string),
		FlavorId:  d.Get("flavor_id").(string),
		Count:     d.Get("node_count").(int),
		Autoscale: d.Get("autoscale").(bool),
		MaxCount:  d.Get("max_count").(int),
		IsDefault: false,
		Labels:    getLabels(d.Get("labels").([]interface{})),
		Taints:    getTaints(d.Get("taints").([]interface{})),
		Cluster:   d.Get("cluster").(string),
	}
}

func getAddons(data []interface{}) []Addon {
	addons := make([]Addon, len(data))
	for i, item := range data {
		addonMap := item.(map[string]interface{})
		addons[i] = Addon{
			Name:    addonMap["name"].(string),
			Version: addonMap["version"].(string),
		}
	}
	return addons
}

func getLabels(data []interface{}) []Label {
	labels := make([]Label, len(data))
	for i, item := range data {
		labelMap := item.(map[string]interface{})
		labels[i] = Label{
			Key:   labelMap["key"].(string),
			Value: labelMap["value"].(string),
		}
	}
	return labels
}

func getTaints(data []interface{}) []Taint {
	taints := make([]Taint, len(data))
	for i, item := range data {
		taintMap := item.(map[string]interface{})
		taints[i] = Taint{
			Key:    taintMap["key"].(string),
			Value:  taintMap["value"].(string),
			Effect: taintMap["effect"].(string),
		}
	}
	return taints
}

func getNodePool(data []interface{}) []NodePoolCreateOptions {
	NodepoolMap := data[0].(map[string]interface{})
	return []NodePoolCreateOptions{{
		Name:      NodepoolMap["name"].(string),
		FlavorId:  NodepoolMap["flavor_id"].(string),
		Count:     NodepoolMap["node_count"].(int),
		Autoscale: NodepoolMap["autoscale"].(bool),
		MaxCount:  NodepoolMap["max_count"].(int),
		IsDefault: true,
		Labels:    make([]Label, 0),
		Taints:    make([]Taint, 0),
	}}
}

func getListOfString(data []interface{}) []string {
	list := make([]string, len(data))
	for i, v := range data {
		list[i] = v.(string)
	}
	return list
}
