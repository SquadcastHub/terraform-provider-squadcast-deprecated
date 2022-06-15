package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

func resourceSlo() *schema.Resource {
	return &schema.Resource{
		Description: "DeduplicationRules resource.",

		CreateContext: resourceSloCreate,
		ReadContext:   resourceSloRead,
		UpdateContext: resourceSloUpdate,
		DeleteContext: resourceSloDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Slo name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Slo description.",
				Type:        schema.TypeString,
				Default:     "Slo created from terraform provider",
				Optional:    true,
			},
			"target_slo": {
				Description: "Target Slo.",
				Type:        schema.TypeFloat,
				Required:    true,
			},
			"owner_type": {
				Description: "Owner type",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "team",
			},
			"owner_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tfutils.ValidateObjectID,
			},
			"org_id": {
				Description:  "Org id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tfutils.ValidateObjectID,
			},
			"service_ids": {
				Description: "Service ids.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"slis": {
				Description: "Slis.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			// TODO: add validation, support only 2 values
			"time_interval_type": {
				Description: "Slo type",
				Type:        schema.TypeString,
				Required:    true,
			},
			"duration_in_days": {
				Description: "Duration in days.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"start_time": {
				Description: "Start time.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"end_time": {
				Description: "End time.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"rules": {
				Description: "Slo monitoring checks.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "id.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"slo_id": {
							Description: "Slo id.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "Name of monitoring check.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"threshold": {
							Description: "Threshold.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"is_checked": {
							Description: "is checked?",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"owner_type": {
							Description: "Owner type",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "team",
						},
						"owner_id": {
							Description: "Team id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
				Optional: true,
			},
			"notify": {
				Description: "Notify.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "id.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"slo_id": {
							Description: "Slo id.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"users": {
							Description: "User ids..",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"squads": {
							Description: "Squad ids..",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"service": {
							Description: "Service id.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"owner_type": {
							Description: "Owner type",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "team",
						},
						"owner_id": {
							Description: "Team id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceSloCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	alerts := make([]*api.SloMonitoringCheck, 0)
	sloActions := make([]*api.SloAction, 0)
	notify := make([]*api.SloNotify, 0)

	err := Decode(d.Get("rules"), &alerts)
	if err != nil {
		return diag.FromErr(err)
	}

	err = Decode(d.Get("notify"), &notify)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, alert := range alerts {
		switch alert.Name {
		case "breached_error_budget":
			alert.Name = "is_breached_err_budget"
		case "unhealthy_slo":
			alert.Name = "is_unhealthy_slo"
		case "increased_false_positives":
			alert.Name = "increased_false_positives_threshold"
		case "remaining_error_budget":
			alert.Name = "remaining_err_budget_threshold"
		}
		alert.OwnerType = "team"
		alert.OwnerID = d.Get("owner_id").(string)
		alert.IsChecked = true
	}

	for _, userID := range notify[0].Users {
		user := &api.SloAction{
			Type:      "USER",
			UserID:    userID,
			OwnerID:   d.Get("owner_id").(string),
			OwnerType: d.Get("owner_type").(string),
		}
		sloActions = append(sloActions, user)
	}

	if notify[0].Service != "" {
		service := &api.SloAction{
			Type:      "SERVICE",
			UserID:    notify[0].Service,
			OwnerID:   d.Get("owner_id").(string),
			OwnerType: d.Get("owner_type").(string),
		}
		sloActions = append(sloActions, service)
	}

	for _, squadID := range notify[0].Squads {
		user := &api.SloAction{
			Type:      "SQUAD",
			UserID:    squadID,
			OwnerID:   d.Get("owner_id").(string),
			OwnerType: d.Get("owner_type").(string),
		}
		sloActions = append(sloActions, user)
	}

	tflog.Info(ctx, "Creating Slos", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	orgID := d.Get("org_id").(string)

	slo, err := client.CreateSlo(ctx, orgID, &api.Slo{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		TargetSlo:           d.Get("target_slo").(float64),
		ServiceIDs:          tfutils.ListToSlice[string](d.Get("service_ids")),
		Slis:                tfutils.ListToSlice[string](d.Get("slis")),
		TimeIntervalType:    d.Get("time_interval_type").(string),
		DurationInDays:      d.Get("duration_in_days").(int),
		StartTime:           d.Get("start_time").(string),
		EndTime:             d.Get("end_time").(string),
		SloMonitoringChecks: alerts,
		SloActions:          sloActions,
		OwnerType:           d.Get("owner_type").(string),
		OwnerID:             d.Get("owner_id").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	idStr := strconv.FormatUint(uint64(slo.ID), 10)
	d.SetId(idStr)
	return resourceSloRead(ctx, d, meta)
}

func resourceSloRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	orgID, ok := d.GetOk("org_id")
	if !ok {
		return diag.Errorf("invalid org id provided")
	}

	sloID, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("invalid slo id")
	}

	tflog.Info(ctx, "Reading Slos", map[string]interface{}{
		"id":       d.Id(),
		"owner_id": d.Get("owner_id").(string),
	})

	slo, err := client.GetSlo(ctx, orgID.(string), sloID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	for _, alert := range slo.SloMonitoringChecks {
		switch alert.Name {
		case "is_breached_err_budget":
			alert.Name = "breached_error_budget"
		case "is_unhealthy_slo":
			alert.Name = "unhealthy_slo"
		case "increased_false_positives_threshold":
			alert.Name = "increased_false_positives"
		case "remaining_err_budget_threshold":
			alert.Name = "remaining_error_budget"
		}
	}

	if err = tfutils.EncodeAndSet(slo, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSloUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	var alerts []*api.SloMonitoringCheck
	sloActions := make([]*api.SloAction, 0)
	notify := make([]*api.SloNotify, 0)

	err := Decode(d.Get("rules"), &alerts)
	if err != nil {
		return diag.FromErr(err)
	}

	err = Decode(d.Get("notify"), &notify)
	if err != nil {
		return diag.FromErr(err)
	}

	sloID, _ := strconv.ParseInt(d.Id(), 10, 32)

	for _, alert := range alerts {
		switch alert.Name {
		case "breached_error_budget":
			alert.Name = "is_breached_err_budget"
		case "unhealthy_slo":
			alert.Name = "is_unhealthy_slo"
		case "increased_false_positives":
			alert.Name = "increased_false_positives_threshold"
		case "remaining_error_budget":
			alert.Name = "remaining_err_budget_threshold"
		}
		alert.OwnerType = "team"
		alert.OwnerID = d.Get("owner_id").(string)
		alert.IsChecked = true
		alert.SloID = sloID
	}

	for _, userID := range notify[0].Users {
		user := &api.SloAction{
			Type:      "USER",
			UserID:    userID,
			SloID:     sloID,
			OwnerID:   d.Get("owner_id").(string),
			OwnerType: d.Get("owner_type").(string),
		}
		sloActions = append(sloActions, user)
	}

	if notify[0].Service != "" {
		service := &api.SloAction{
			Type:      "SERVICE",
			UserID:    notify[0].Service,
			SloID:     sloID,
			OwnerID:   d.Get("owner_id").(string),
			OwnerType: d.Get("owner_type").(string),
		}
		sloActions = append(sloActions, service)
	}

	for _, squadID := range notify[0].Squads {
		user := &api.SloAction{
			Type:      "SQUAD",
			UserID:    squadID,
			SloID:     sloID,
			OwnerID:   d.Get("owner_id").(string),
			OwnerType: d.Get("owner_type").(string),
		}
		sloActions = append(sloActions, user)
	}

	tflog.Info(ctx, "Updating Slos", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	orgID := d.Get("org_id").(string)
	id := d.Id()

	tflog.Info(ctx, "Updating Slos", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	_, err = client.UpdateSlo(ctx, orgID, id, &api.Slo{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		TargetSlo:           d.Get("target_slo").(float64),
		ServiceIDs:          tfutils.ListToSlice[string](d.Get("service_ids")),
		Slis:                tfutils.ListToSlice[string](d.Get("slis")),
		TimeIntervalType:    d.Get("time_interval_type").(string),
		DurationInDays:      d.Get("duration_in_days").(int),
		StartTime:           d.Get("start_time").(string),
		EndTime:             d.Get("end_time").(string),
		SloMonitoringChecks: alerts,
		SloActions:          sloActions,
		OwnerType:           d.Get("owner_type").(string),
		OwnerID:             d.Get("owner_id").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSloRead(ctx, d, meta)
}

func resourceSloDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Deleting Slos", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	_, err := client.DeleteSlo(ctx, d.Get("org_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
