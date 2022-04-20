package provider

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Team resource.",

		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Team id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Team name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Description:  "Team description.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
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
			// "members": {
			// 	Type:     schema.TypeSet,
			// 	Required: true,
			// 	MinItems: 1,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"user_id": {
			// 				Description: "User id.",
			// 				Type:        schema.TypeString,
			// 				Required:    true,
			// 			},
			// 			"roles": {
			// 				Type:     schema.TypeSet,
			// 				Required: true,
			// 				Elem:     schema.TypeString,
			// 				MinItems: 1,
			// 			},
			// 		},
			// 	},
			// },
			"roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Role id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description:  "Role name.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
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
						// "abilities": {
						// 	Type:     schema.TypeList,
						// 	Computed: true,
						// 	Elem:     schema.TypeString,
						// },
					},
				},
			},
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// client := meta.(*api.Client)

	roles := d.Get("roles").(*schema.Set).List()
	spew.Dump(roles)

	// team, err := client.CreateTeam(ctx, api.TeamCreateReq{
	// 	Name:        d.Get("name").(string),
	// 	Description: d.Get("description").(string),
	// 	// MembersIDs:  memberIds,
	// })
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// d.SetId(team.ID)

	// spew.Dump(d.GetRawState())

	// if err = tfutils.EncodeAndSet(team, d); err != nil {
	// 	return diag.FromErr(err)
	// }

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	// tflog.Trace(ctx, "created a resource")

	return nil
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}
