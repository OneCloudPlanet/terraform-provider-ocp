package ocp_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const ClusterAddonsUri = "cluster/addons"

type ClusterAddon struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Releases    []Release `json:"releases"`
}

type Release struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

func (c *Client) ClusterAddons(ctx context.Context) ([]ClusterAddon, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodGet, ClusterAddonsUri, nil)
	if err != nil {
		return nil, err
	}
	var result []ClusterAddon
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []ClusterAddon{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result, nil
}
