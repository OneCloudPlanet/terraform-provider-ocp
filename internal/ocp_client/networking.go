package ocp_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const NetworkingUri = "cluster/networking"

type Networking struct {
	ID      string `json:"id"`
	Name    string `json:"network_name"`
	Version string `json:"version"`
}

func (c *Client) Networking(ctx context.Context) ([]Networking, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodGet, NetworkingUri, nil)
	if err != nil {
		return nil, err
	}
	var result []Networking
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []Networking{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}
	return result, nil
}
