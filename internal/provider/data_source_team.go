package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "What is a squadcast team?",

		ReadContext: dataSourceTeamRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Team id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description:  "Team name.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
			"description": {
				Description: "Team description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"slug": {
				Description: "Team slug.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default": {
				Description: "Team is default?.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"created_at": {
				Description: "Team created at.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Team updated at.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_by": {
				Description: "Team created by.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Description: "User id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"role_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schema.TypeString,
						},
					},
				},
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Role id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Role name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"slug": {
							Description: "Role slug.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"default": {
							Description: "Role is default?.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"abilities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	id, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	team, err := client.GetTeamById(ctx, id.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tfutils.EncodeAndSet(team, d); err != nil {
		return diag.FromErr(err)
	}

	fmt.Printf("\n\nstate is %s \n\n\n", d.State().String())

	return nil
}
