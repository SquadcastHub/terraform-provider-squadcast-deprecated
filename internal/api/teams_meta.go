package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type teamMetaRole struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type TeamMeta struct {
	ID          string          `json:"id" tf:"id"`
	Name        string          `json:"name" tf:"name"`
	Description string          `json:"description" tf:"description"`
	Default     bool            `json:"default" tf:"default"`
	Roles       []*teamMetaRole `json:"roles" tf:"-"`
}

func (t *TeamMeta) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(t)
	if err != nil {
		return nil, err
	}

	defaultRoleNames := map[string]string{
		"Manage Team": "manage_team",
		"Admin":       "admin",
		"User":        "user",
		"Observer":    "observer",
	}

	roles := map[string]interface{}{}

	for _, role := range t.Roles {
		key := defaultRoleNames[role.Name]
		if key != "" {
			roles[key] = role.ID
		}
	}
	m["default_role_ids"] = roles

	return m, nil
}

func (client *Client) GetTeamMetaById(ctx context.Context, id string) (*TeamMeta, error) {
	url := fmt.Sprintf("%s/teams/%s", client.BaseURLV3, id)

	return Request[any, TeamMeta](http.MethodGet, url, client, ctx, nil)
}

type CreateTeamReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (client *Client) CreateTeam(ctx context.Context, req *CreateTeamReq) (*TeamMeta, error) {
	url := fmt.Sprintf("%s/teams", client.BaseURLV3)

	return Request[CreateTeamReq, TeamMeta](http.MethodPost, url, client, ctx, req)
}

type UpdateTeamMetaReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (client *Client) UpdateTeamMeta(ctx context.Context, id string, req *UpdateTeamMetaReq) (*TeamMeta, error) {
	url := fmt.Sprintf("%s/teams/%s/meta", client.BaseURLV3, id)

	return Request[UpdateTeamMetaReq, TeamMeta](http.MethodPatch, url, client, ctx, req)
}

// type UpdateTeamMembersReq struct {
// 	Members []*TeamMember `json:"members"`
// }

// func (client *Client) UpdateTeamMembers(ctx context.Context, id string, req *UpdateTeamMembersReq) (*TeamMeta, error) {
// 	url := fmt.Sprintf("%s/teams/%s/members", client.BaseURLV3, id)

// 	return Request[UpdateTeamMembersReq, TeamMeta](http.MethodPatch, url, client, ctx, req)
// }

func (client *Client) DeleteTeam(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/teams/%s", client.BaseURLV3, id)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
