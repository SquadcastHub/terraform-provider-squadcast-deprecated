package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Host   string
	Region string

	RefreshToken string
	AccessToken  string

	UserAgent   string
	BaseURL     string
	AuthBaseURL string
}

// Meta holds the status of the request informations
type Meta struct {
	Meta struct {
		Status  int    `json:"status"`
		Message string `json:"error_message,omitempty"`
	} `json:"meta,omitempty"`
}

func Request[TReq interface{}, TRes interface{}](method string, path string, client *Client, ctx context.Context, payload *TReq) (*TRes, error) {
	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequestWithContext(ctx, method, client.BaseURL+path, nil)
	} else {
		buf := &bytes.Buffer{}
		if payload != nil {
			body, err := json.Marshal(payload)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(body)
		}
		req, err = http.NewRequestWithContext(ctx, method, client.BaseURL+path, buf)
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.AccessToken))
	req.Header.Set("User-Agent", client.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data *TRes `json:"data"`
		*Meta
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		if resp.StatusCode > 299 {
			return nil, fmt.Errorf("%s %s returned an unexpected error with no body", method, path)
		} else {
			return nil, nil
		}
	}

	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		if response.Meta != nil {
			return nil, fmt.Errorf("%s %s returned %d: %s", method, path, response.Meta.Meta.Status, response.Meta.Meta.Message)
		} else {
			return nil, fmt.Errorf("%s %s returned an unexpected error: %#v", method, path, response)
		}
	}

	return response.Data, nil
}

func RequestSlice[TReq interface{}, TRes interface{}](method string, path string, client *Client, ctx context.Context, payload *TReq) ([]*TRes, error) {
	data, err := Request[TReq, []*TRes](method, path, client, ctx, payload)
	if err != nil {
		return nil, err
	}

	return *data, nil
}
