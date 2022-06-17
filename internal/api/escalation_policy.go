package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tf"
)

type EscalationPolicy struct {
	ID          string   `json:"id" tf:"id"`
	Name        string   `json:"name" tf:"name"`
	Description string   `json:"description" tf:"description"`
	Owner       OwnerRef `json:"owner" tf:"-"`
}

func (s *EscalationPolicy) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}

	m["team_id"] = s.Owner.ID

	return m, nil
}

func (client *Client) GetEscalationPolicyByName(ctx context.Context, teamID string, name string) ([]*EscalationPolicy, error) {
	url := fmt.Sprintf("%s/escalation-policies?name=%s&owner_id=%s", client.BaseURLV3, name, teamID)
	return RequestSlice[any, EscalationPolicy](http.MethodGet, url, client, ctx, nil)
}
