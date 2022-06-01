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
	path := fmt.Sprintf("/services/%s?owner_id=%s", id, teamID)

	return Request[any, Service](http.MethodGet, path, client, ctx, nil)
}

func (client *Client) GetServiceByName(ctx context.Context, teamID string, name string) (*Service, error) {
	path := fmt.Sprintf("/services/by-name?name=%s&owner_id=%s", name, teamID)

	return Request[any, Service](http.MethodGet, path, client, ctx, nil)
}

func (client *Client) ListServices(ctx context.Context, teamID string) ([]*Service, error) {
	path := fmt.Sprintf("/services?owner_id=%s", teamID)

	return RequestSlice[any, Service](http.MethodGet, path, client, ctx, nil)
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

func (client *Client) CreateService(ctx context.Context, req *CreateServiceReq) (*Service, error) {
	path := fmt.Sprintf("/services")
	return Request[CreateServiceReq, Service](http.MethodPost, path, client, ctx, req)
}

func (client *Client) UpdateService(ctx context.Context, id string, req *UpdateServiceReq) (*Service, error) {
	path := fmt.Sprintf("/services/%s", id)
	return Request[UpdateServiceReq, Service](http.MethodPut, path, client, ctx, req)
}

func (client *Client) DeleteService(ctx context.Context, id string) (*any, error) {
	path := fmt.Sprintf("/services/%s", id)
	return Request[any, any](http.MethodDelete, path, client, ctx, nil)
}
