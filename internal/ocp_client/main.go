package ocp_client

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	UARegion = "ua"
	PLRegion = "pl"
)

var regions = map[string]string{
	UARegion: "core.ocplanet.cloud",
	PLRegion: "core-pl.ocplanet.cloud",
}

type Client struct {
	Region string
	API    *API
}

func GetClient(token, region string) (*Client, error) {
	endpoint, err := getEndpoint(region)
	if err != nil {
		return nil, err
	}
	client := &Client{
		Region: region,
		API: &API{
			HTTPClient: http.DefaultClient,
			Token:      token,
			Endpoint:   endpoint,
			UserAgent:  userAgent,
		},
	}
	return client, nil
}

func getEndpoint(region string) (string, error) {
	var baseRoute string
	if v, ok := regions[region]; ok {
		baseRoute = v
	} else {
		return "", errors.New(fmt.Sprintf("Region %s does not exist", region))
	}
	return fmt.Sprintf("https://%s/backend/api/", baseRoute), nil
}
