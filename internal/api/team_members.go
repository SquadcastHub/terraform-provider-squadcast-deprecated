package api

import (
	"context"
	"fmt"
	"net/http"
)

type UpdateTeamMembersReq struct {
	Members []*TeamMember `json:"members"`
}

func (client *Client) UpdateTeamMembers(ctx context.Context, teamID string, req *UpdateTeamMembersReq) (*Team, error) {
	url := fmt.Sprintf("%s/teams/%s/update-members", client.BaseURLV3, teamID)

	return Request[UpdateTeamMembersReq, Team](http.MethodPatch, url, client, ctx, req)
}
