package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tf"
)

type Data struct {
	Slo *Slo `json:"slo,omitempty"`
}

type Slo struct {
	ID                  uint                  `json:"id,omitempty" tf:"id"`
	Name                string                `json:"name" tf:"name"`
	Description         string                `json:"description,omitempty" tf:"description"`
	TimeIntervalType    string                `json:"time_interval_type" tf:"time_interval_type"`
	ServiceIDs          []string              `json:"service_ids" tf:"service_ids"`
	Slis                []string              `json:"slis" tf:"slis"`
	TargetSlo           float64               `json:"target_slo" tf:"target_slo"`
	StartTime           string                `json:"start_time,omitempty" tf:"start_time"`
	EndTime             string                `json:"end_time,omitempty" tf:"end_time"`
	DurationInDays      int                   `json:"duration_in_days,omitempty" tf:"duration_in_days"`
	SloMonitoringChecks []*SloMonitoringCheck `json:"slo_monitoring_checks" tf:"rules"`
	SloActions          []*SloAction          `json:"slo_actions" tf:"notify"`
	OwnerType           string                `json:"owner_type" tf:"owner_type"`
	OwnerID             string                `json:"owner_id" tf:"owner_id"`
}

type SloMonitoringCheck struct {
	ID        uint   `json:"id,omitempty" tf:"id"`
	SloID     int64  `json:"slo_id,omitempty" tf:"slo_id"`
	Name      string `json:"name" tf:"name"`
	Threshold int    `json:"threshold" tf:"threshold"`
	OwnerType string `json:"owner_type" tf:"owner_type"`
	OwnerID   string `json:"owner_id" tf:"owner_id"`
	IsChecked bool   `json:"is_checked" tf:"is_checked"`
}

type SloAction struct {
	ID        uint   `json:"id,omitempty" tf:"id"`
	SloID     int64  `json:"slo_id,omitempty" tf:"slo_id"`
	Type      string `json:"type" tf:"type"`
	UserID    string `json:"user_id" tf:"user_id"`
	SquadID   string `json:"squad_id" tf:"squad_id"`
	ServiceID string `json:"service_id" tf:"service_id"`
	OwnerType string `json:"owner_type" tf:"owner_type"`
	OwnerID   string `json:"owner_id" tf:"owner_id"`
}

type SloNotify struct {
	ID        uint     `json:"id,omitempty" tf:"id"`
	SloID     int64    `json:"slo_id,omitempty" tf:"slo_id"`
	Users     []string `json:"users" tf:"users"`
	Squads    []string `json:"squads" tf:"squads"`
	Service   string   `json:"service" tf:"service"`
	OwnerType string   `json:"owner_type" tf:"owner_type"`
	OwnerID   string   `json:"owner_id" tf:"owner_id"`
}

func (c *SloMonitoringCheck) Encode() (map[string]interface{}, error) {
	return tf.Encode(c)
}

func (c *SloNotify) Encode() (map[string]interface{}, error) {
	return tf.Encode(c)
}

func (r *Slo) Encode() (map[string]interface{}, error) {
	notify := make([]*SloNotify, 0)
	notify = append(notify, &SloNotify{})

	m, err := tf.Encode(r)
	if err != nil {
		return nil, err
	}

	sloMonitoringChecks, err := tf.EncodeSlice(r.SloMonitoringChecks)
	if err != nil {
		return nil, err
	}
	m["rules"] = sloMonitoringChecks

	for _, n := range r.SloActions {
		if n.UserID != "" {
			notify[0].Users = append(notify[0].Users, n.UserID)
		}
		if n.SquadID != "" {
			notify[0].Squads = append(notify[0].Squads, n.SquadID)
		}
		if n.ServiceID != "" {
			notify[0].Service = n.ServiceID
		}
	}

	if len(r.SloActions) > 0 {
		notify[0].SloID = int64(r.ID)
		notify[0].OwnerID = r.OwnerID
		notify[0].OwnerType = r.OwnerType
	}

	notifyObj, err := tf.EncodeSlice(notify)
	if err != nil {
		fmt.Println(err)

	}
	m["notify"] = notifyObj

	return m, nil
}

func (r *Data) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(r)
	if err != nil {
		return nil, err
	}

	slo, err := tf.Encode(r.Slo)
	if err != nil {
		return nil, err
	}
	m["slo"] = slo

	return m, nil
}

func (client *Client) CreateSlo(ctx context.Context, orgID string, req *Slo) (*Slo, error) {
	url := fmt.Sprintf("%s/slo?owner_type=team&owner_id=611262fcd5b4ea846b534a8a", client.BaseURLV3)
	data, er := Request[Slo, Data](http.MethodPost, url, client, ctx, req)
	return data.Slo, er
}

func (client *Client) GetSlo(ctx context.Context, orgID, sloID string) (*Slo, error) {
	url := fmt.Sprintf("%s/slo/%s?owner_type=team&owner_id=611262fcd5b4ea846b534a8a", client.BaseURLV3, sloID)
	data, er := Request[any, Data](http.MethodGet, url, client, ctx, nil)
	if data != nil {
		return data.Slo, er
	} else {
		return nil, errors.New("Slo not found")
	}
}

func (client *Client) UpdateSlo(ctx context.Context, orgID, sloID string, req *Slo) (*Slo, error) {
	url := fmt.Sprintf("%s/slo/%s?owner_type=team&owner_id=611262fcd5b4ea846b534a8a", client.BaseURLV3, sloID)
	return Request[Slo, Slo](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteSlo(ctx context.Context, orgID, sloID string) (*any, error) {
	url := fmt.Sprintf("%s/slo/%s?owner_type=team&owner_id=611262fcd5b4ea846b534a8a", client.BaseURLV3, sloID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
