package ocp_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const FlavorsUri = "openstack/instances/create_options"

type Flavor struct {
	ID                       string   `json:"id"`
	Name                     string   `json:"name"`
	Description              string   `json:"description"`
	Vcpus                    int      `json:"vcpus"`
	MemoryMb                 int      `json:"memory_mb"`
	MemoryGb                 float64  `json:"memory_gb"`
	RootGb                   int      `json:"root_gb"`
	AssignedClusterTemplates []string `json:"assigned_cluster_templates"`
	EphemeralGb              int      `json:"ephemeral_gb"`
	FlavorGroup              *string  `json:"flavor_group"`
	OutOfStock               bool     `json:"out_of_stock"`
	Properties               string   `json:"properties"`
	Region                   string   `json:"region"`
	ResellerResources        *string  `json:"reseller_resources"`
	Swap                     int      `json:"swap"`
	UsedByResellers          []string `json:"used_by_resellers"`
}

func (c *Client) Flavors(ctx context.Context) ([]Flavor, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodGet, FlavorsUri, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Flavor []Flavor `json:"flavor"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []Flavor{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Flavor, nil
}
