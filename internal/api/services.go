package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type Service struct {
	ID                 string   `json:"id" tf:"id"`
	Name               string   `json:"name" tf:"name"`
	APIKey             string   `json:"api_key" tf:"api_key"`
	Email              string   `json:"email" tf:"email"`
	Description        string   `json:"description" tf:"description"`
	EscalationPolicyID string   `json:"escalation_policy_id" tf:"escalation_policy_id"`
	OnMaintenance      bool     `json:"on_maintenance" tf:"-"`
	Owner              OwnerRef `json:"owner" tf:"-"`
	Dependencies       []string `json:"depends" tf:"dependencies"`
}

func (s *Service) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(s)
	if err != nil {
		return nil, err
	}

	m["team_id"] = s.Owner.ID

	return m, nil
}

func (client *Client) GetServiceById(ctx context.Context, teamID string, id string) (*Service, error) {
	url := fmt.Sprintf("%s/services/%s?owner_id=%s", client.BaseURLV3, id, teamID)

	return Request[any, Service](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetServiceByName(ctx context.Context, teamID string, name string) (*Service, error) {
	url := fmt.Sprintf("%s/services/by-name?name=%s&owner_id=%s", client.BaseURLV3, name, teamID)

	return Request[any, Service](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) ListServices(ctx context.Context, teamID string) ([]*Service, error) {
	url := fmt.Sprintf("%s/services?owner_id=%s", client.BaseURLV3, teamID)

	return RequestSlice[any, Service](http.MethodGet, url, client, ctx, nil)
}

type CreateServiceReq struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	TeamID             string `json:"owner_id"`
	EscalationPolicyID string `json:"escalation_policy_id"`
	EmailPrefix        string `json:"email_prefix"`
}

type UpdateServiceReq struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	EscalationPolicyID string `json:"escalation_policy_id"`
	EmailPrefix        string `json:"email_prefix"`
}

type UpdateServiceDependenciesReq struct {
	Data []string `json:"data"`
}

func (client *Client) CreateService(ctx context.Context, req *CreateServiceReq) (*Service, error) {
	url := fmt.Sprintf("%s/services", client.BaseURLV3)
	return Request[CreateServiceReq, Service](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateService(ctx context.Context, id string, req *UpdateServiceReq) (*Service, error) {
	url := fmt.Sprintf("%s/services/%s", client.BaseURLV3, id)
	return Request[UpdateServiceReq, Service](http.MethodPut, url, client, ctx, req)
}

func (client *Client) UpdateServiceDependencies(ctx context.Context, id string, req *UpdateServiceDependenciesReq) (*any, error) {
	url := fmt.Sprintf("%s/organizations/%s/services/%s/dependencies", client.BaseURLV2, client.OrganizationID, id)
	return Request[UpdateServiceDependenciesReq, any](http.MethodPost, url, client, ctx, req)
}

func (client *Client) DeleteService(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/services/%s", client.BaseURLV3, id)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
