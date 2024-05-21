package ocp_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const ClusterUri = "cluster/"

func (c *Client) CreateCluster(ctx context.Context, data interface{}) (map[string]interface{}, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodPost, ClusterUri, data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return result, fmt.Errorf("error when decoding json response, %w", err)
	}
	return result, nil
}

func (c *Client) GetCluster(ctx context.Context, clusterId string) (map[string]interface{}, error) {
	resp, code, err := c.API.makeRequest(ctx, http.MethodGet, ClusterUri+clusterId+"/", nil)
	if err != nil {
		if code == 404 {
			return nil, nil
		}
		return nil, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return result, fmt.Errorf("error when decoding json response, %w", err)
	}
	return result, nil
}

func (c *Client) UpdateCluster(ctx context.Context, clusterId string, data interface{}) (map[string]interface{}, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodPatch, ClusterUri+clusterId+"/", data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return result, fmt.Errorf("error when decoding json response, %w", err)
	}
	return result, nil
}

func (c *Client) DeleteCluster(ctx context.Context, clusterId string) (map[string]interface{}, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodDelete, ClusterUri+clusterId+"/", nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return result, fmt.Errorf("error when decoding json response, %w", err)
	}
	return result, nil
}
