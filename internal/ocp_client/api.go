package ocp_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const userAgent = "ocp-terraform/0.1"

type API struct {
	HTTPClient *http.Client
	Token      string
	Endpoint   string
	UserAgent  string
}

func (api *API) makeRequest(ctx context.Context, method, uri string, params interface{}) ([]byte, int, error) {
	jsonBody, err := HandleParams(params)
	if err != nil {
		return nil, 0, err
	}

	var resp *http.Response
	var respErr error
	var reqBody io.Reader
	var respBody []byte
	if jsonBody != nil {
		reqBody = bytes.NewReader(jsonBody)
	}

	resp, respErr = api.request(ctx, method, uri, reqBody)
	if respErr != nil || resp.StatusCode >= http.StatusInternalServerError {
		if respErr == nil {
			respBody, err = io.ReadAll(resp.Body)
			resp.Body.Close() //nolint

			respErr = fmt.Errorf("could not read response body, %w", err)
			fmt.Printf("Request: %s %s got an error response %d", method, uri, resp.StatusCode)
		} else {
			fmt.Printf("Error performing request: %s %s : %s \n", method, uri, respErr.Error())
		}
	} else {
		respBody, err = io.ReadAll(resp.Body)
		defer resp.Body.Close() //nolint
		if err != nil {
			return nil, resp.StatusCode, fmt.Errorf("could not read response body, %w", err)
		}
	}
	if respErr != nil {
		if resp != nil {
			return nil, resp.StatusCode, respErr
		}
		return nil, 0, respErr
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, resp.StatusCode, handleStatusCode(resp.StatusCode, respBody, uri)
	}

	return respBody, resp.StatusCode, nil
}

func (api *API) request(ctx context.Context, method, uri string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, api.Endpoint+uri, body)
	if err != nil {
		return nil, fmt.Errorf("HTTP request creation failed, %w", err)
	}

	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Set("Authorization", "OpenAPIToken "+api.Token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed, %w", err)
	}

	return resp, nil
}

func HandleParams(params interface{}) ([]byte, error) {
	var jsonBody []byte
	var err error

	if params == nil {
		return nil, nil
	}

	if paramBytes, ok := params.([]byte); ok {
		jsonBody = paramBytes
	} else {
		jsonBody, err = json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("error marshalling params to JSON, %w", err)
		}
	}
	return jsonBody, nil
}

func handleStatusCode(statusCode int, body []byte, uri string) error {
	if statusCode >= http.StatusInternalServerError {
		return fmt.Errorf("http status %d: service failed.\n%v\n%v", statusCode, body, uri)
	}
	errBody := errors.New(fmt.Sprintf("Bad requset: %s", body))
	return errBody
}
