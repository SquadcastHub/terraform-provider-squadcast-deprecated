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

const serviceMaintenanceID = "service_maintenance"

func resourceServiceMaintenance() *schema.Resource {
	return &schema.Resource{
		Description: "ServiceMaintenance resource.",

		CreateContext: resourceServiceMaintenanceCreate,
		ReadContext:   resourceServiceMaintenanceRead,
		UpdateContext: resourceServiceMaintenanceUpdate,
		DeleteContext: resourceServiceMaintenanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServiceMaintenanceImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ServiceMaintenance id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"service_id": {
				Description:  "Service id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tfutils.ValidateObjectID,
			},
			"window": {
				Description: "window",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Description:  "from",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
						"till": {
							Description:  "till.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
						"repeat_till": {
							Description:  "repeat till.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
						"repeat_frequency": {
							Description:  "repeat frequency.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"day", "week", "2 weeks", "3 weeks", "month"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceServiceMaintenanceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	_, serviceID, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	// d.Set("team_id", teamID)
	d.Set("service_id", serviceID)
	d.SetId(serviceMaintenanceID)

	return []*schema.ResourceData{d}, nil
}

func resourceServiceMaintenanceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	var windows []api.ServiceMaintenanceWindow
	err := Decode(d.Get("window"), &windows)
	if err != nil {
		return diag.FromErr(err)
	}

	updateWindows := make([]api.UpdateServiceMaintenanceWindowsWindow, 0, len(windows))
	for _, w := range windows {
		uw := api.UpdateServiceMaintenanceWindowsWindow{
			From:       w.From,
			Till:       w.Till,
			RepeatTill: w.RepeatTill,
		}
		if w.RepeatFrequency == "" {
			uw.RepeatTill = uw.Till
		}
		if w.RepeatFrequency == "day" {
			uw.Daily = true
		} else if w.RepeatFrequency == "week" {
			uw.Weekly = true
		} else if w.RepeatFrequency == "2 weeks" {
			uw.TwoWeekly = true
		} else if w.RepeatFrequency == "3 weeks" {
			uw.ThreeWeekly = true
		} else if w.RepeatFrequency == "month" {
			uw.Monthly = true
		}
		updateWindows = append(updateWindows, uw)
	}

	_, err = client.UpdateServiceMaintenance(ctx, d.Get("service_id").(string), &api.UpdateServiceMaintenanceWindows{
		OrganizationID: client.OrganizationID,
		ServiceID:      d.Get("service_id").(string),
		Data: api.UpdateServiceMaintenanceWindowsData{
			ServiceMaintenanceWindows: updateWindows,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceMaintenanceID)

	return resourceServiceMaintenanceRead(ctx, d, meta)
}

func resourceServiceMaintenanceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	serviceID, ok := d.GetOk("service_id")
	if !ok {
		return diag.Errorf("invalid service id provided")
	}

	tflog.Info(ctx, "Reading service maintenance", map[string]interface{}{
		"service_id": serviceID,
	})
	serviceMaintenanceWindows, err := client.GetServiceMaintenanceWindows(ctx, serviceID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	windows, err := tfutils.EncodeSlice(serviceMaintenanceWindows)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("window", windows)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServiceMaintenanceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return resourceServiceMaintenanceCreate(ctx, d, meta)
}

func resourceServiceMaintenanceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateServiceMaintenance(ctx, d.Get("service_id").(string), &api.UpdateServiceMaintenanceWindows{
		OrganizationID: client.OrganizationID,
		ServiceID:      d.Get("service_id").(string),
		Data: api.UpdateServiceMaintenanceWindowsData{
			ServiceMaintenanceWindows: []api.UpdateServiceMaintenanceWindowsWindow{},
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}