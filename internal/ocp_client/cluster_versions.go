package ocp_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const ClusterVersionsUri = "cluster/versions"

type ClusterVersion struct {
	ID      string  `json:"id"`
	Version string  `json:"version"`
	Images  []Image `json:"images"`
}

type Image struct {
	Name        string `json:"name"`
	ImageName   string `json:"image_name"`
	OpenstackId string `json:"openstack_id"`
	OsDistro    string `json:"os_distro"`
}

func (c *Client) ClusterVersions(ctx context.Context) ([]ClusterVersion, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodGet, ClusterVersionsUri, nil)
	if err != nil {
		return nil, err
	}
	var result []ClusterVersion
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []ClusterVersion{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result, nil
}
