package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
	"github.com/hashicorp/terraform-provider-squadcast/internal/tf"
)

func resourceSlo() *schema.Resource {
	return &schema.Resource{
		Description: "Slo resource.",

		CreateContext: resourceSloCreate,
		ReadContext:   resourceSloRead,
		UpdateContext: resourceSloUpdate,
		DeleteContext: resourceSloDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Slo id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Slo name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "Slo description.",
				Type:        schema.TypeString,
				Default:     "Slo created from terraform provider",
				Optional:    true,
			},
			"target_slo": {
				Description: "Slo target.",
				Type:        schema.TypeFloat,
				Required:    true,
			},
			"service_ids": {
				Description: "Slo service ids.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"slis": {
				Description: "Slo slis.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"time_interval_type": {
				Description:  "Slo type",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"rolling", "fixed"}, false),
			},
			"duration_in_days": {
				Description: "Slo duration in days.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"start_time": {
				Description:  "Slo start time.",
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"end_time": {
				Description:  "Slo end time.",
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
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
							ValidateFunc: validation.StringInSlice([]string{"breached_error_budget", "unhealthy_slo",
								"increased_false_positives", "remaining_error_budget"}, false),
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
				Description: "Slo notify.",
				Type:        schema.TypeList,
				MaxItems:    1,
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
							Description:  "Service id.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tf.ValidateObjectID,
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
			"owner_type": {
				Description: "Slo owner type",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "team",
			},
			"owner_id": {
				Description:  "Slo team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"org_id": {
				Description:  "Slo org id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
		},
	}
}

var alertsMap = map[string]string{"is_breached_err_budget": "breached_error_budget",
	"breached_error_budget":               "is_breached_err_budget",
	"is_unhealthy_slo":                    "unhealthy_slo",
	"unhealthy_slo":                       "is_unhealthy_slo",
	"increased_false_positives_threshold": "increased_false_positives",
	"increased_false_positives":           "increased_false_positives_threshold",
	"remaining_err_budget_threshold":      "remaining_error_budget",
	"remaining_error_budget":              "remaining_err_budget_threshold",
}

func resourceSloCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	rules := make([]*api.SloMonitoringCheck, 0)
	notify := make([]*api.SloNotify, 0)
	sloActions := make([]*api.SloAction, 0)

	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	err = Decode(d.Get("notify"), &notify)
	if err != nil {
		return diag.FromErr(err)
	}

	ownerID := d.Get("owner_id").(string)
	ownerType := "team"

	sloActions = formatRulesAndNotify(rules, notify, ownerID, 0)

	tflog.Info(ctx, "Creating Slos", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	orgID := d.Get("org_id").(string)

	slo, err := client.CreateSlo(ctx, orgID, &api.Slo{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		TargetSlo:           d.Get("target_slo").(float64),
		ServiceIDs:          tf.ListToSlice[string](d.Get("service_ids")),
		Slis:                tf.ListToSlice[string](d.Get("slis")),
		TimeIntervalType:    d.Get("time_interval_type").(string),
		DurationInDays:      d.Get("duration_in_days").(int),
		StartTime:           d.Get("start_time").(string),
		EndTime:             d.Get("end_time").(string),
		SloMonitoringChecks: rules,
		SloActions:          sloActions,
		OwnerType:           ownerType,
		OwnerID:             ownerID,
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
		alert.Name = alertsMap[alert.Name]
	}

	if err = tf.EncodeAndSet(slo, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSloUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	var rules []*api.SloMonitoringCheck
	sloActions := make([]*api.SloAction, 0)
	notify := make([]*api.SloNotify, 0)

	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	err = Decode(d.Get("notify"), &notify)
	if err != nil {
		return diag.FromErr(err)
	}

	sloID, _ := strconv.ParseInt(d.Id(), 10, 32)
	ownerID := d.Get("owner_id").(string)
	ownerType := "team"

	sloActions = formatRulesAndNotify(rules, notify, ownerID, sloID)

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
		ServiceIDs:          tf.ListToSlice[string](d.Get("service_ids")),
		Slis:                tf.ListToSlice[string](d.Get("slis")),
		TimeIntervalType:    d.Get("time_interval_type").(string),
		DurationInDays:      d.Get("duration_in_days").(int),
		StartTime:           d.Get("start_time").(string),
		EndTime:             d.Get("end_time").(string),
		SloMonitoringChecks: rules,
		SloActions:          sloActions,
		OwnerType:           ownerType,
		OwnerID:             ownerID,
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

// formatRulesAndNotify transform the payload into the format expected by the API and terraform state
func formatRulesAndNotify(rules []*api.SloMonitoringCheck, notify []*api.SloNotify, ownerID string, sloID int64) []*api.SloAction {
	sloActions := make([]*api.SloAction, 0)
	for _, alert := range rules {
		alert.Name = alertsMap[alert.Name]
		alert.IsChecked = true
		alert.SloID = sloID
		alert.OwnerType = "team"
		alert.OwnerID = ownerID
	}

	for _, userID := range notify[0].Users {
		user := &api.SloAction{
			Type:      "USER",
			UserID:    userID,
			SloID:     sloID,
			OwnerID:   ownerID,
			OwnerType: "team",
		}
		sloActions = append(sloActions, user)
	}

	for _, squadID := range notify[0].Squads {
		user := &api.SloAction{
			Type:      "SQUAD",
			UserID:    squadID,
			SloID:     sloID,
			OwnerID:   ownerID,
			OwnerType: "team",
		}
		sloActions = append(sloActions, user)
	}

	if notify[0].Service != "" {
		service := &api.SloAction{
			Type:      "SERVICE",
			UserID:    notify[0].Service,
			SloID:     sloID,
			OwnerID:   ownerID,
			OwnerType: "team",
		}
		sloActions = append(sloActions, service)
	}

	return sloActions
}
