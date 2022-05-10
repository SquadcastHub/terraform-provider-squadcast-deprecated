package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type TeamRef struct {
	ID   string `json:"id" tf:"id"`
	Name string `json:"name" tf:"name"`
}

func (tr *TeamRef) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(tr)
}

type UserRef struct {
	ID   string `json:"id" tf:"id"`
	Name string `json:"name" tf:"name"`
}

func (ur *UserRef) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(ur)
}

type Squad struct {
	ID      string     `json:"id" tf:"id"`
	Name    string     `json:"name" tf:"name"`
	Slug    string     `json:"slug" tf:"slug"`
	Team    TeamRef    `json:"team" tf:"-"`
	Members []*UserRef `json:"members" tf:"-"`
}

func (s *Squad) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(s)
	if err != nil {
		return nil, err
	}

	m["team_id"] = s.Team.ID

	team, err := tfutils.Encode(s.Team)
	if err != nil {
		return nil, err
	}
	m["team"] = team

	members, err := tfutils.EncodeSlice(s.Members)
	if err != nil {
		return nil, err
	}
	m["members"] = members

	return m, nil
}

func (client *Client) GetSquadById(ctx context.Context, teamID string, id string) (*Squad, error) {
	path := fmt.Sprintf("/teams/%s/squads/%s", teamID, id)

	return Request[any, Squad](http.MethodGet, path, client, ctx, nil)
}

func (client *Client) ListSquads(ctx context.Context, teamID string) ([]*Squad, error) {
	path := fmt.Sprintf("/teams/%s/squads", teamID)

	return RequestSlice[any, Squad](http.MethodGet, path, client, ctx, nil)
}

type CreateSquadReq struct {
	Name      string   `json:"name"`
	MemberIDs []string `json:"memberIds"`
}

func (client *Client) CreateSquad(ctx context.Context, teamID string, req *CreateSquadReq) (*Squad, error) {
	path := fmt.Sprintf("/teams/%s/squads", teamID)
	return Request[CreateSquadReq, Squad](http.MethodPost, path, client, ctx, req)
}
