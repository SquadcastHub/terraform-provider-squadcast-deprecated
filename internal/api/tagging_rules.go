package api

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type TaggingRuleCondition struct {
	LHS string `json:"lhs" tf:"lhs"`
	Op  string `json:"op" tf:"op"`
	RHS string `json:"rhs" tf:"rhs"`
}

func (c *TaggingRuleCondition) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(c)
}

type TaggingRuleTagValue struct {
	Value string `json:"value" tf:"value"`
	Color string `json:"color" tf:"color"`
}

func (tv *TaggingRuleTagValue) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(tv)
}

type TaggingRule struct {
	IsBasic         bool                           `json:"is_basic" tf:"is_basic"`
	Expression      string                         `json:"expression" tf:"expression"`
	BasicExpression []*TaggingRuleCondition        `json:"basic_expression" tf:"basic_expressions"`
	Tags            map[string]TaggingRuleTagValue `json:"tags" tf:"-"`
}

func (r *TaggingRule) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(r)
	if err != nil {
		return nil, err
	}

	basicExpression, err := tfutils.EncodeSlice(r.BasicExpression)
	if err != nil {
		return nil, err
	}
	m["basic_expressions"] = basicExpression

	tags := make([]interface{}, 0, len(r.Tags))

	keys := make([]string, 0, len(r.Tags))
	for k := range r.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := r.Tags[k]
		mtag, err := v.Encode()
		if err != nil {
			return nil, err
		}
		mtag["key"] = k
		tags = append(tags, mtag)
	}
	m["tags"] = tags

	return m, nil
}

type TaggingRules struct {
	ID        string         `json:"id" tf:"id"`
	ServiceID string         `json:"service_id" tf:"service_id"`
	Rules     []*TaggingRule `json:"rules" tf:"-"`
}

func (s *TaggingRules) Encode() (map[string]interface{}, error) {
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

func (client *Client) GetTaggingRules(ctx context.Context, serviceID, teamID string) (*TaggingRules, error) {
	url := fmt.Sprintf("%s/services/%s/tagging-rules", client.BaseURLV3, serviceID)

	return Request[any, TaggingRules](http.MethodGet, url, client, ctx, nil)
}

type UpdateTaggingRulesReq struct {
	Rules []TaggingRule `json:"rules"`
}

func (client *Client) UpdateTaggingRules(ctx context.Context, serviceID, teamID string, req *UpdateTaggingRulesReq) (*TaggingRules, error) {
	url := fmt.Sprintf("%s/services/%s/tagging-rules", client.BaseURLV3, serviceID)

	return Request[UpdateTaggingRulesReq, TaggingRules](http.MethodPost, url, client, ctx, req)
}
