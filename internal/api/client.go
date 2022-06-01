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

	RefreshToken   string
	AccessToken    string
	OrganizationID string

	UserAgent   string
	BaseURLV2   string
	BaseURLV3   string
	AuthBaseURL string
}

type AppError struct {
	Status  int    `json:"status"`
	Message string `json:"error_message,omitempty"`
}

// Meta holds the status of the request informations
type Meta struct {
	Meta AppError `json:"meta,omitempty"`
}

func Request[TReq interface{}, TRes interface{}](method string, url string, client *Client, ctx context.Context, payload *TReq) (*TRes, error) {
	var req *http.Request
	var resData *TRes
	var err error

	if method == "GET" {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	} else {
		buf := &bytes.Buffer{}
		if payload != nil {
			body, err := json.Marshal(payload)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(body)
		}
		req, err = http.NewRequestWithContext(ctx, method, url, buf)
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
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
			return nil, fmt.Errorf("%s %s returned an unexpected error with no body", method, url)
		} else {
			return nil, nil
		}
	}

	err = json.Unmarshal(bytes, &response)
	if err == nil {
		resData = response.Data
	} else {
		// its an array, so unmarshalling is different
		var responseWithArray struct {
			Data []*TRes `json:"data"`
			*Meta
		}
		err1 := json.Unmarshal(bytes, &responseWithArray)
		if err1 != nil {
			return nil, err1
		}
		resData = responseWithArray.Data[0]
	}

	if resp.StatusCode > 299 {
		if response.Meta != nil {
			return nil, fmt.Errorf("%s %s returned %d: %s", method, url, response.Meta.Meta.Status, response.Meta.Meta.Message)
		} else {
			return nil, fmt.Errorf("%s %s returned %d with an unexpected error: %#v", method, url, resp.StatusCode, response)
		}
	}

	return resData, nil
}

// func RequestNew[TReq interface{}, TRes []interface{}](method string, url string, client *Client, ctx context.Context, payload *TReq) (TRes, error) {
// 	var req *http.Request
// 	var err error

// 	if method == "GET" {
// 		req, err = http.NewRequestWithContext(ctx, method, url, nil)
// 	} else {
// 		buf := &bytes.Buffer{}
// 		if payload != nil {
// 			body, err := json.Marshal(payload)
// 			if err != nil {
// 				return nil, err
// 			}
// 			buf = bytes.NewBuffer(body)
// 		}
// 		req, err = http.NewRequestWithContext(ctx, method, url, buf)
// 		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
// 	}

// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.AccessToken))
// 	req.Header.Set("User-Agent", client.UserAgent)

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var response struct {
// 		Data *TRes `json:"data"`
// 		*Meta
// 	}

// 	defer resp.Body.Close()
// 	bytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(bytes) == 0 {
// 		if resp.StatusCode > 299 {
// 			return nil, fmt.Errorf("%s %s returned an unexpected error with no body", method, url)
// 		} else {
// 			return nil, nil
// 		}
// 	}

// 	if err := json.Unmarshal(bytes, &response); err != nil {
// 		return nil, err
// 		if err != nil {
// 			var response struct {
// 				Data *[]TRes `json:"data"`
// 				*Meta
// 			}
// 			err := json.Unmarshal(bytes, &response)
// 			if err != nil{

// 			}

// 		}
// 	}

// 	if resp.StatusCode > 299 {
// 		if response.Meta != nil {
// 			return nil, fmt.Errorf("%s %s returned %d: %s", method, url, response.Meta.Meta.Status, response.Meta.Meta.Message)
// 		} else {
// 			return nil, fmt.Errorf("%s %s returned %d with an unexpected error: %#v", method, url, resp.StatusCode, response)
// 		}
// 	}

// 	return response.Data, nil
// }

func RequestSlice[TReq interface{}, TRes interface{}](method string, url string, client *Client, ctx context.Context, payload *TReq) ([]*TRes, error) {
	data, err := Request[TReq, []*TRes](method, url, client, ctx, payload)
	if err != nil {
		return nil, err
	}

	return *data, nil
}

func Get[TRes interface{}](client *Client, ctx context.Context, path string) (*TRes, error) {
	return Request[interface{}, TRes](http.MethodGet, path, client, ctx, nil)
}

func Post[TRes interface{}](client *Client, ctx context.Context, path string, payload interface{}) (*TRes, error) {
	return Request[interface{}, TRes](http.MethodPost, path, client, ctx, &payload)
}
