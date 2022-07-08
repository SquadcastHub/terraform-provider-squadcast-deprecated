package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func dataSourceEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "EscalationPolicy data source.",

		ReadContext: dataSourceEscalationPolicyRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "EscalationPolicy id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "EscalationPolicy name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "EscalationPolicy description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"repeat": {
				Description: "repeat this policy",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"times": {
							Description: "repeat times",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"delay_minutes": {
							Description: "repeat after minutes",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"rules": {
				Description: "rules.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delay_minutes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"targets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"notification_channels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"round_robin": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Description: "enable rotation within",
										Type:        schema.TypeBool,
										Computed:    true,
									},
									"rotation": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Description: "enable rotation within",
													Type:        schema.TypeBool,
													Computed:    true,
												},
												"delay_minutes": {
													Description: "repeat after minutes",
													Type:        schema.TypeInt,
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"repeat": {
							Description: "repeat this rule",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"times": {
										Description: "repeat times",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"delay_minutes": {
										Description: "repeat after minutes",
										Type:        schema.TypeInt,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceEscalationPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading escalation_policy", tf.M{
		"name": d.Get("name").(string),
	})
	escalationPolicy, err := client.GetEscalationPolicyByName(ctx, d.Get("team_id").(string), d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(escalationPolicy, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
