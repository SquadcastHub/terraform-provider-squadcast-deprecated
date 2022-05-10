package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type AccessToken struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	IssuedAt     int64  `json:"issued_at"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

func (client *Client) GetAccessToken(ctx context.Context) error {
	path := "/oauth/access-token"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.AuthBaseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Refresh-Token", client.RefreshToken)
	req.Header.Set("User-Agent", client.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var response struct {
		Data AccessToken `json:"data"`
		*Meta
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &response); err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(response.Meta.Meta.Message)
	}

	client.AccessToken = response.Data.AccessToken
	fmt.Printf("\nAccess token: %#v\n", client.AccessToken)
	return nil
}
