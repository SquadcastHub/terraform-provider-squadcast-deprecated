package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type Team struct {
	ID          string        `json:"id" tf:"id"`
	CreatedAt   string        `json:"created_at" tf:"created_at"`
	UpdatedAt   string        `json:"updated_at" tf:"updated_at"`
	CreatedBy   string        `json:"created_by" tf:"created_by"`
	Name        string        `json:"name" tf:"name"`
	Description string        `json:"description" tf:"description"`
	Slug        string        `json:"slug" tf:"slug"`
	Default     bool          `json:"default" tf:"default"`
	Members     []*TeamMember `json:"members" tf:"-"`
	Roles       []*TeamRole   `json:"roles" tf:"-"`
}

func (t *Team) Encode() (map[string]interface{}, error) {
	m, err := tfutils.EncodeGeneric(t)
	if err != nil {
		return nil, err
	}

	membersmap, err := tfutils.EncodeSliceGeneric(t.Members)
	if err != nil {
		return nil, err
	}
	m["members"] = membersmap

	rolesmap, err := tfutils.EncodeSlice(t.Roles)
	if err != nil {
		return nil, err
	}
	m["roles"] = rolesmap

	return m, nil
}

type TeamMember struct {
	UserID  string   `json:"user_id" tf:"user_id"`
	RoleIDs []string `json:"role_ids" tf:"role_ids"`
}

type TeamRole struct {
	ID        string                 `json:"id" tf:"id"`
	Name      string                 `json:"name" tf:"name"`
	Slug      string                 `json:"slug" tf:"slug"`
	Default   bool                   `json:"default" tf:"default"`
	Abilities RBACEntityAbilitiesMap `json:"abilities" tf:"-"`
}

func (tr *TeamRole) Encode() (map[string]interface{}, error) {
	m, err := tfutils.EncodeGeneric(tr)
	if err != nil {
		return nil, err
	}

	abilities := make([]string, 0, 10)
	for _, kv := range tr.Abilities {
		for k := range kv {
			abilities = append(abilities, k)
		}
	}

	m["abilities"] = abilities

	return m, nil
}

type RBACAbilityMap map[string]bool
type RBACEntityAbilitiesMap map[string]RBACAbilityMap

func (client *Client) GetTeamById(ctx context.Context, id string) (*Team, error) {
	path := fmt.Sprintf("/teams/%s", id)

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
		Data Team `json:"data"`
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
