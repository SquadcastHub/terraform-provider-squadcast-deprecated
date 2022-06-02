package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type Team struct {
	ID          string        `json:"id" tf:"id"`
	Name        string        `json:"name" tf:"name"`
	Description string        `json:"description" tf:"description"`
	Default     bool          `json:"default" tf:"default"`
	Members     []*TeamMember `json:"members" tf:"-"`
	Roles       []*TeamRole   `json:"roles" tf:"-"`
}

func (t *Team) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(t)
	if err != nil {
		return nil, err
	}

	members, err := tfutils.EncodeSlice(t.Members)
	if err != nil {
		return nil, err
	}
	m["members"] = members

	roles, err := tfutils.EncodeSlice(t.Roles)
	if err != nil {
		return nil, err
	}
	m["roles"] = roles

	return m, nil
}

type TeamMember struct {
	UserID  string   `json:"user_id" tf:"user_id"`
	RoleIDs []string `json:"role_ids" tf:"role_ids"`
}

func (tm *TeamMember) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(tm)
}

type TeamRole struct {
	ID        string                 `json:"id" tf:"id"`
	Name      string                 `json:"name" tf:"name"`
	Slug      string                 `json:"slug" tf:"-"`
	Default   bool                   `json:"default" tf:"default"`
	Abilities RBACEntityAbilitiesMap `json:"abilities" tf:"-"`
}

func (tr *TeamRole) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(tr)
	if err != nil {
		return nil, err
	}

	abilities := make([]string, 0, 100)
	for _, kv := range tr.Abilities {
		for k := range kv {
			abilities = append(abilities, k)
		}
	}

	sort.Strings(abilities)
	m["abilities"] = abilities

	return m, nil
}

type RBACAbilityMap map[string]bool
type RBACEntityAbilitiesMap map[string]RBACAbilityMap

// func (client *Client) GetTeamById(ctx context.Context, id string) (*Team, error) {
// 	path := fmt.Sprintf("/teams/%s", id)

// 	return Get[Team](client, ctx, path)
// }

func (client *Client) GetTeamByName(ctx context.Context, name string) (*Team, error) {
	url := fmt.Sprintf("%s/teams/by-name?name=%s", client.BaseURLV3, url.QueryEscape(name))

	return Request[any, Team](http.MethodGet, url, client, ctx, nil)
}

// type TeamCreateReq struct {
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	// MembersIDs  string `json:"members_ids"`
// }

// func (client *Client) CreateTeam(ctx context.Context, req TeamCreateReq) (*Team, error) {
// 	path := fmt.Sprintf("/teams")

// 	body, err := json.Marshal(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	team, err := Post[Team](client, ctx, path, bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return team, nil
// }

func Get[TRes interface{}](client *Client, ctx context.Context, path string) (*TRes, error) {
	return Request[interface{}, TRes](http.MethodGet, path, client, ctx, nil)
}

func Post[TRes interface{}](client *Client, ctx context.Context, path string, payload interface{}) (*TRes, error) {
	return Request[interface{}, TRes](http.MethodPost, path, client, ctx, &payload)
}
