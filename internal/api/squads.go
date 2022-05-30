package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type OwnerRef struct {
	ID   string `json:"id" tf:"id"`
	Type string `json:"type" tf:"type"`
}

type Squad struct {
	ID        string   `json:"id" tf:"id"`
	Name      string   `json:"name" tf:"name"`
	Slug      string   `json:"slug" tf:"-"`
	Owner     OwnerRef `json:"owner" tf:"-"`
	MemberIDs []string `json:"members" tf:"member_ids"`
}

func (s *Squad) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(s)
	if err != nil {
		return nil, err
	}

	m["team_id"] = s.Owner.ID

	return m, nil
}

func (client *Client) GetSquadById(ctx context.Context, teamID string, id string) (*Squad, error) {
	path := fmt.Sprintf("/squads/%s?owner_id=%s", teamID, id)

	return Request[any, Squad](http.MethodGet, path, client, ctx, nil)
}

func (client *Client) ListSquads(ctx context.Context, teamID string) ([]*Squad, error) {
	path := fmt.Sprintf("/squads?owner_id=%s", teamID)

	return RequestSlice[any, Squad](http.MethodGet, path, client, ctx, nil)
}

type CreateSquadReq struct {
	Name      string   `json:"name"`
	TeamID    string   `json:"owner_id"`
	MemberIDs []string `json:"members"`
}

func (client *Client) CreateSquad(ctx context.Context, req *CreateSquadReq) (*Squad, error) {
	path := fmt.Sprintf("/squads")
	return Request[CreateSquadReq, Squad](http.MethodPost, path, client, ctx, req)
}

func (client *Client) UpdateSquad(ctx context.Context, req *Squad) (*Squad, error) {
	path := fmt.Sprintf("/squads/%s", req.ID)
	return Request[Squad, Squad](http.MethodPut, path, client, ctx, req)
}

func (client *Client) DeleteSquad(ctx context.Context, id string) (*any, error) {
	path := fmt.Sprintf("/squads/%s", id)
	return Request[any, any](http.MethodDelete, path, client, ctx, nil)
}
