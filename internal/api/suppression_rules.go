package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type SuppressionRuleCondition struct {
	LHS string `json:"lhs" tf:"lhs"`
	Op  string `json:"op" tf:"op"`
	RHS string `json:"rhs" tf:"rhs"`
}

func (c *SuppressionRuleCondition) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(c)
}

type SuppressionRule struct {
	IsBasic         bool                        `json:"is_basic" tf:"is_basic"`
	Description     string                      `json:"description" tf:"description"`
	Expression      string                      `json:"expression" tf:"expression"`
	BasicExpression []*SuppressionRuleCondition `json:"basic_expression" tf:"basic_expression"`
}

func (r *SuppressionRule) Encode() (map[string]interface{}, error) {
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

type SuppressionRules struct {
	ID        string             `json:"id" tf:"id"`
	ServiceID string             `json:"service_id" tf:"service_id"`
	Rules     []*SuppressionRule `json:"rules" tf:"-"`
}

func (s *SuppressionRules) Encode() (map[string]interface{}, error) {
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

func (client *Client) GetSuppressionRules(ctx context.Context, serviceID, teamID string) (*SuppressionRules, error) {
	path := fmt.Sprintf("/services/%s/suppression-rules", serviceID)

	return Request[any, SuppressionRules](http.MethodGet, path, client, ctx, nil)
}

type UpdateSuppressionRulesReq struct {
	Rules []SuppressionRule `json:"rules"`
}

func (client *Client) UpdateSuppressionRules(ctx context.Context, serviceID, teamID string, req *UpdateSuppressionRulesReq) (*SuppressionRules, error) {
	path := fmt.Sprintf("/services/%s/suppression-rules", serviceID)
	return Request[UpdateSuppressionRulesReq, SuppressionRules](http.MethodPost, path, client, ctx, req)
}
