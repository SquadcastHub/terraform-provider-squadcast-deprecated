package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		Description: "Service resource.",

		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServiceImport,
		},

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
			},
			"description": {
				Description:  "Service description.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tfutils.ValidateObjectID,
			},
			"escalation_policy_id": {
				Description:  "Escalation policy id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tfutils.ValidateObjectID,
			},
			"email_prefix": {
				Description: "Email prefix.",
				Type:        schema.TypeString,
				Required:    true,
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
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tfutils.ValidateObjectID,
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

func resourceServiceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	teamID, id, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating service", map[string]interface{}{
		"name": d.Get("name").(string),
	})
	service, err := client.CreateService(ctx, &api.CreateServiceReq{
		Name:               d.Get("name").(string),
		TeamID:             d.Get("team_id").(string),
		Description:        d.Get("description").(string),
		EscalationPolicyID: d.Get("escalation_policy_id").(string),
		EmailPrefix:        d.Get("email_prefix").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(service.ID)

	_, err = client.UpdateServiceDependencies(ctx, service.ID, &api.UpdateServiceDependenciesReq{
		Data: tfutils.ListToSlice[string](d.Get("dependencies")),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading service", map[string]interface{}{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})
	service, err := client.GetServiceById(ctx, teamID.(string), id)
	if err != nil {
		return diag.FromErr(err)
	}

	alertSources, err := client.ListAlertSources(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	service.AlertSources = alertSources.Available().EndpointMap(client.IngestionBaseURL, service)

	if err = tfutils.EncodeAndSet(service, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateService(ctx, d.Id(), &api.UpdateServiceReq{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		EscalationPolicyID: d.Get("escalation_policy_id").(string),
		EmailPrefix:        d.Get("email_prefix").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateServiceDependencies(ctx, d.Id(), &api.UpdateServiceDependenciesReq{
		Data: tfutils.ListToSlice[string](d.Get("dependencies")),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteService(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
