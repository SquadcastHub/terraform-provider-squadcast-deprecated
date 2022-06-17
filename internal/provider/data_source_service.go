package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
	"github.com/hashicorp/terraform-provider-squadcast/internal/tf"
)

func dataSourceService() *schema.Resource {
	return &schema.Resource{
		Description: "What is a squadcast service?",
		ReadContext: dataSourceServiceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Service id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Service name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
				ForceNew:     true,
			},
			"description": {
				Description: "Service description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"escalation_policy_id": {
				Description: "Escalation policy id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email_prefix": {
				Description: "Email prefix.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"api_key": {
				Description: "Email prefix.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "Email.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dependencies": {
				Description: "dependencies.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tf.ValidateObjectID,
				},
			},
			"alert_source_endpoints": {
				Description: "alert sources.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceServiceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	name, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("invalid service name provided")
	}

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading service by name", tf.M{
		"name": name.(string),
	})
	service, err := client.GetServiceByName(ctx, teamID.(string), name.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	alertSources, err := client.ListAlertSources(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	service.AlertSources = alertSources.Available().EndpointMap(client.IngestionBaseURL, service)

	if err = tf.EncodeAndSet(service, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
