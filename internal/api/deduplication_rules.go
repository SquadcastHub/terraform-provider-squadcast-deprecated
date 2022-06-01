package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type DeduplicationRuleCondition struct {
	LHS string `json:"lhs" tf:"lhs"`
	Op  string `json:"op" tf:"op"`
	RHS string `json:"rhs" tf:"rhs"`
}

func (c *DeduplicationRuleCondition) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(c)
}

type DeduplicationRule struct {
	IsBasic                 bool                          `json:"is_basic" tf:"is_basic"`
	Description             string                        `json:"description" tf:"description"`
	Expression              string                        `json:"expression" tf:"expression"`
	DependencyDeduplication bool                          `json:"dependency_deduplication" tf:"dependency_deduplication"`
	TimeUnit                string                        `json:"time_unit" tf:"time_unit"`
	TimeWindow              int                           `json:"time_window" tf:"time_window"`
	BasicExpression         []*DeduplicationRuleCondition `json:"basic_expression" tf:"basic_expression"`
}

func (r *DeduplicationRule) Encode() (map[string]interface{}, error) {
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

type DeduplicationRules struct {
	ID        string               `json:"id" tf:"id"`
	ServiceID string               `json:"service_id" tf:"service_id"`
	Rules     []*DeduplicationRule `json:"rules" tf:"-"`
}

func (s *DeduplicationRules) Encode() (map[string]interface{}, error) {
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

func (client *Client) GetDeduplicationRules(ctx context.Context, serviceID, teamID string) (*DeduplicationRules, error) {
	url := fmt.Sprintf("%s/services/%s/deduplication-rules", client.BaseURLV3, serviceID)

	return Request[any, DeduplicationRules](http.MethodGet, url, client, ctx, nil)
}

type UpdateDeduplicationRulesReq struct {
	Rules []DeduplicationRule `json:"rules"`
}

func (client *Client) UpdateDeduplicationRules(ctx context.Context, serviceID, teamID string, req *UpdateDeduplicationRulesReq) (*DeduplicationRules, error) {
	url := fmt.Sprintf("%s/services/%s/deduplication-rules", client.BaseURLV3, serviceID)
	return Request[UpdateDeduplicationRulesReq, DeduplicationRules](http.MethodPost, url, client, ctx, req)
}
