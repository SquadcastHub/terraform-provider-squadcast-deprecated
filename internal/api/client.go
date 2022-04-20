package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Host   string
	Region string

	RefreshToken string
	AccessToken  string

	UserAgent string
	BaseURL   string
}

// Meta holds the status of the request informations
type Meta struct {
	Meta struct {
		Status  int    `json:"status_code"`
		Message string `json:"error_message,omitempty"`
	} `json:"meta,omitempty"`
}

func Get[T interface{}](client *Client, ctx context.Context, path string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.BaseURL+path, nil)
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
		Data T `json:"data"`
		*Meta
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(response.Meta.Meta.Message)
	}

	return &response.Data, nil
}

func Post[T interface{}](client *Client, ctx context.Context, path string, body io.Reader) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.BaseURL+path, body)
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
		Data T `json:"data"`
		*Meta
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(response.Meta.Meta.Message)
	}

	return &response.Data, nil
}
