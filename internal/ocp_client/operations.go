package ocp_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const OperationsUri = "operations/"

func (c *Client) GetOperation(ctx context.Context, id string) (map[string]interface{}, int, error) {
	resp, code, err := c.API.makeRequest(ctx, http.MethodGet, OperationsUri+id, nil)
	if err != nil {
		return nil, code, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return result, code, fmt.Errorf("error when decoding json response, %w", err)
	}
	return result, code, nil
}
