package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

// type DeduplicationRuleCondition struct {
// 	LHS string `json:"lhs" tf:"lhs"`
// 	Op  string `json:"op" tf:"op"`
// 	RHS string `json:"rhs" tf:"rhs"`
// }

type Data struct {
	Slo *Slo `json:"slo,omitempty"`
}

type Slo struct {
	ID               uint     `json:"id,omitempty" tf:"id"`
	Name             string   `json:"name" tf:"name"`
	Description      string   `json:"description,omitempty" tf:"description"`
	TimeIntervalType string   `json:"time_interval_type" tf:"time_interval_type"`
	ServiceIDs       []string `json:"service_ids" tf:"service_ids"`
	Slis             []string `json:"slis" tf:"slis"`
	TargetSlo        float64  `json:"target_slo" tf:"target_slo"`
	StartTime        string   `json:"start_time,omitempty" tf:"start_time"`
	EndTime          string   `json:"end_time,omitempty" tf:"end_time"`
	DurationInDays   int      `json:"duration_in_days,omitempty" tf:"duration_in_days"`
	// Tags                json.RawMessage       `json:"tags,omitempty" tf:"tags"`
	SloMonitoringChecks []*SloMonitoringCheck `json:"slo_monitoring_checks" tf:"rules"`
	// SloActions          []SloAction          `json:"slo_actions"`
	OwnerType string `json:"owner_type" tf:"owner_type"`
	OwnerID   string `json:"owner_id" tf:"owner_id"`
}

type SloMonitoringCheck struct {
	ID uint `json:"id,omitempty" tf:"id"`
	// SloID     uint   `json:"slo_id,omitempty"`
	Name      string `json:"name" tf:"name"`
	Threshold int    `json:"threshold" tf:"threshold"`
	OwnerType string `json:"owner_type" tf:"owner_type"`
	OwnerID   string `json:"owner_id" tf:"owner_id"`
	IsChecked bool   `json:"is_checked" tf:"is_checked"`
}

// type SloAction struct {
// 	ID        uint   `json:"id,omitempty"`
// 	SloID     uint   `json:"slo_id,omitempty"`
// 	Type      string `json:"type"`
// 	UserID    string `json:"user_id"`
// 	SquadID   string `json:"squad_id"`
// 	ServiceID string `json:"service_id"`
// 	OwnerType string `json:"owner_type"`
// 	OwnerID   string `json:"owner_id"`
// }

func (c *SloMonitoringCheck) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(c)
}

// type DeduplicationRule struct {
// 	IsBasic                 bool                          `json:"is_basic" tf:"is_basic"`
// 	Description             string                        `json:"description" tf:"description"`
// 	Expression              string                        `json:"expression" tf:"expression"`
// 	DependencyDeduplication bool                          `json:"dependency_deduplication" tf:"dependency_deduplication"`
// 	TimeUnit                string                        `json:"time_unit" tf:"time_unit"`
// 	TimeWindow              int                           `json:"time_window" tf:"time_window"`
// 	BasicExpression         []*DeduplicationRuleCondition `json:"basic_expression" tf:"basic_expression"`
// }

func (r *Slo) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(r)
	if err != nil {
		return nil, err
	}

	slo_monitoring_checks, err := tfutils.EncodeSlice(r.SloMonitoringChecks)
	if err != nil {
		return nil, err
	}
	m["rules"] = slo_monitoring_checks

	return m, nil
}

func (r *Data) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(r)
	if err != nil {
		return nil, err
	}

	slo, err := tfutils.Encode(r.Slo)
	if err != nil {
		return nil, err
	}
	m["slo"] = slo

	return m, nil
}

// type DeduplicationRules struct {
// 	ID        string               `json:"id" tf:"id"`
// 	ServiceID string               `json:"service_id" tf:"service_id"`
// 	Rules     []*DeduplicationRule `json:"rules" tf:"-"`
// }

// func (s *DeduplicationRules) Encode() (map[string]interface{}, error) {
// 	m, err := tfutils.Encode(s)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rules, err := tfutils.EncodeSlice(s.Rules)
// 	if err != nil {
// 		return nil, err
// 	}
// 	m["rules"] = rules

// 	return m, nil
// }

// func (client *Client) CreateService(ctx context.Context, req *CreateServiceReq) (*Service, error) {
// 	url := fmt.Sprintf("%s/services", client.BaseURLV3)
// 	return Request[CreateServiceReq, Service](http.MethodPost, url, client, ctx, req)
// }

func (client *Client) CreateSlo(ctx context.Context, orgID string, req *Slo) (*Slo, error) {
	url := fmt.Sprintf("%s/slo?owner_type=team&owner_id=611262fcd5b4ea846b534a8a", client.BaseURLV3)
	a, er := Request[Slo, Data](http.MethodPost, url, client, ctx, req)
	fmt.Println(er)
	return a.Slo, er
}

func (client *Client) GetSlo(ctx context.Context, orgID, sloID string) (*Slo, error) {
	url := fmt.Sprintf("%s/slo/%s?owner_type=team&owner_id=611262fcd5b4ea846b534a8a", client.BaseURLV3, sloID)
	a, er := Request[any, Data](http.MethodGet, url, client, ctx, nil)
	fmt.Println(a, er)
	return a.Slo, er
}

func (client *Client) UpdateSlo(ctx context.Context, orgID, sloID string, req *Slo) (*Slo, error) {
	url := fmt.Sprintf("%s/slo/%s?owner_type=team&owner_id=611262fcd5b4ea846b534a8a", client.BaseURLV3, sloID)
	a, er := Request[Slo, Data](http.MethodPut, url, client, ctx, req)
	return a.Slo, er
}

func (client *Client) DeleteSlo(ctx context.Context, orgID, sloID string) (*any, error) {
	url := fmt.Sprintf("%s/slo/%s?owner_type=team&owner_id=611262fcd5b4ea846b534a8a", client.BaseURLV3, sloID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
