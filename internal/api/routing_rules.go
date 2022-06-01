package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type RoutingRuleCondition struct {
	LHS string `json:"lhs" tf:"lhs"`
	RHS string `json:"rhs" tf:"rhs"`
}

func (c *RoutingRuleCondition) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(c)
}

type RouteTo struct {
	EntityID   string `json:"entity_id" tf:"route_to_id"`
	EntityType string `json:"entity_type" tf:"route_to_type"`
}

type RoutingRule struct {
	IsBasic         bool                    `json:"is_basic" tf:"is_basic"`
	Expression      string                  `json:"expression" tf:"expression"`
	BasicExpression []*RoutingRuleCondition `json:"basic_expression" tf:"basic_expression"`
	RouteTo         RouteTo                 `json:"route_to" tf:"route_to,squash"`
}

func (r *RoutingRule) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(r)
	if err != nil {
		return nil, err
	}

	basicExpression, err := tfutils.EncodeSlice(r.BasicExpression)
	if err != nil {
		return nil, err
	}
	m["basic_expression"] = basicExpression

	return m, nil
}

type RoutingRules struct {
	ID        string         `json:"id" tf:"id"`
	ServiceID string         `json:"service_id" tf:"service_id"`
	Rules     []*RoutingRule `json:"rules" tf:"-"`
}

func (s *RoutingRules) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(s)
	if err != nil {
		return nil, err
	}

	rules, err := tfutils.EncodeSlice(s.Rules)
	if err != nil {
		return nil, err
	}
	m["rules"] = rules

	return m, nil
}

func (client *Client) GetRoutingRules(ctx context.Context, serviceID, teamID string) (*RoutingRules, error) {
	url := fmt.Sprintf("%s/services/%s/routing-rules", client.BaseURLV3, serviceID)

	return Request[any, RoutingRules](http.MethodGet, url, client, ctx, nil)
}

type UpdateRoutingRulesReq struct {
	Rules []RoutingRule `json:"rules"`
}

func (client *Client) UpdateRoutingRules(ctx context.Context, serviceID, teamID string, req *UpdateRoutingRulesReq) (*RoutingRules, error) {
	url := fmt.Sprintf("%s/services/%s/routing-rules", client.BaseURLV3, serviceID)
	return Request[UpdateRoutingRulesReq, RoutingRules](http.MethodPost, url, client, ctx, req)
}
