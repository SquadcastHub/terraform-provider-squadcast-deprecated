package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mvmaasakkers/go-problemdetails"
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
		Status  int    `json:"status_code"`
		Message string `json:"error_message,omitempty"`
	} `json:"meta,omitempty"`
}

// func Request[T interface{}, U interface{}](method string, path string) func(*Client, context.Context, *T) (*U, error) {
// 	return func(client *Client, ctx context.Context, payload *T) (*U, error) {
// 		var buf *bytes.Buffer

// 		if payload != nil {
// 			body, err := json.Marshal(payload)
// 			if err != nil {
// 				return nil, err
// 			}
// 			buf = bytes.NewBuffer(body)
// 		}

// 		req, err := http.NewRequestWithContext(ctx, method, client.BaseURL+path, buf)
// 		if err != nil {
// 			return nil, err
// 		}
// 		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.AccessToken))
// 		req.Header.Set("User-Agent", client.UserAgent)

// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			return nil, err
// 		}

// 		var p problemdetails.ProblemDetails
// 		var d U

// 		defer resp.Body.Close()
// 		bytes, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if resp.StatusCode > 299 {
// 			if err := json.Unmarshal(bytes, &p); err != nil {
// 				return nil, err
// 			}
// 			return nil, &p
// 		}

// 		if err := json.Unmarshal(bytes, &d); err != nil {
// 			return nil, err
// 		}

// 		return &d, nil
// 	}
// }

func Request[T interface{}, U interface{}](method string, path string, client *Client, ctx context.Context, payload *T) (*U, error) {
	var buf *bytes.Buffer

	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, client.BaseURL+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.AccessToken))
	req.Header.Set("User-Agent", client.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var prob problemdetails.ProblemDetails
	var data U

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		if err := json.Unmarshal(bytes, &prob); err != nil {
			return nil, err
		}
		return nil, &prob
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func RequestSlice[T interface{}, U interface{}](method string, path string, client *Client, ctx context.Context, payload *T) ([]*U, error) {
	data, err := Request[T, []*U](method, path, client, ctx, payload)
	if err != nil {
		return nil, err
	}

	return *data, nil
}
