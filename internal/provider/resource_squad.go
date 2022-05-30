package provider

import (
	"context"
	"log"

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
		},
	}
}

func resourceSquadCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	squad, err := client.CreateSquad(ctx, &api.CreateSquadReq{
		Name:      d.Get("name").(string),
		MemberIDs: tfutils.SetToSlice[string](d.Get("member_ids")),
		TeamID:    d.Get("team_id").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(squad.ID)

	return resourceSquadRead(ctx, d, meta)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	// tflog.Trace(ctx, "created a resource")
}

func resourceSquadRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	squad, err := client.GetSquadById(ctx, id, teamID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tfutils.EncodeAndSet(squad, d); err != nil {
		return diag.FromErr(err)
	}

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

	client := meta.(*api.Client)

	log.Printf("[DEBUG] Deleting Squad: %s", d.Id())
	_, err := client.DeleteSquad(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
