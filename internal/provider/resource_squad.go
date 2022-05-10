package provider

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

func resourceSquad() *schema.Resource {
	return &schema.Resource{
		Description: "Squad resource.",

		CreateContext: resourceSquadCreate,
		ReadContext:   resourceSquadRead,
		UpdateContext: resourceSquadUpdate,
		DeleteContext: resourceSquadDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Squad id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Squad name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tfutils.ValidateObjectID,
			},
			"member_ids": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"team": {
				Description: "Team.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Team id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Team name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"members": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "User id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "User name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceSquadCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	squad, err := client.CreateSquad(ctx, d.Get("team_id").(string), &api.CreateSquadReq{
		Name:      d.Get("name").(string),
		MemberIDs: tfutils.SetToSlice[string](d.Get("member_ids")),
	})
	if err != nil {
		fmt.Printf("\nerror here %#v\n", err)
		return diag.FromErr(err)
	}

	d.SetId(squad.ID)

	spew.Dump(d.GetRawState())

	if err = tfutils.EncodeAndSet(squad, d); err != nil {
		return diag.FromErr(err)
	}

	spew.Dump(d.GetRawState())

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	// tflog.Trace(ctx, "created a resource")

	return nil
}

func resourceSquadRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("invalid squad id provided")
	}

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	squad, err := client.GetSquadById(ctx, id.(string), teamID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tfutils.EncodeAndSet(squad, d); err != nil {
		return diag.FromErr(err)
	}

	spew.Dump(d.GetRawState())

	return nil
}

func resourceSquadUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceSquadDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}
