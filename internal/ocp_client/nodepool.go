package ocp_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const NodePoolUri = "node-pool/"

func (c *Client) GetNodePool(ctx context.Context, nodepoolId string) (map[string]interface{}, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodGet, NodePoolUri+nodepoolId, nil)
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

func (c *Client) CreateNodePool(ctx context.Context, data interface{}) (map[string]interface{}, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodPost, NodePoolUri, data)
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

func (c *Client) UpdateNodePool(ctx context.Context, NodePoolId string, data interface{}) (map[string]interface{}, error) {
	resp, _, err := c.API.makeRequest(ctx, http.MethodPatch, NodePoolUri+NodePoolId+"/", data)
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

func (c *Client) DeleteNodePool(ctx context.Context, NodePoolId string) error {
	_, _, err := c.API.makeRequest(ctx, http.MethodDelete, NodePoolUri+NodePoolId+"/", nil)
	if err != nil {
		return err
	}

	return nil
}
