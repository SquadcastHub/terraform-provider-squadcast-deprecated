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

func dataSourceEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "What is a squadcast escalation policy?",
		ReadContext: dataSourceEscalationPolicyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Escalation policy id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Escalation policy name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "Escalation policy description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
		},
	}
}

func dataSourceEscalationPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	name, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("invalid escalation policy name provided")
	}

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading escalation policy by name", map[string]interface{}{
		"name": name.(string),
	})

	escalationPolicy, err := client.GetEscalationPolicyByName(ctx, teamID.(string), name.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if escalationPolicy.Name == "" {
		return diag.Errorf("Unable to find escalation policy with the name %s", name)
	}

	if err = tf.EncodeAndSet(escalationPolicy, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
