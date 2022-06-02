package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
)

func dataSourceAlertEndpoint() *schema.Resource {
	return &schema.Resource{
		Description: "Get Webhook URL for given service and alertsource",
		ReadContext: dataSourceAlertEndpointRead,
		Schema: map[string]*schema.Schema{
			"alert_source_name": {
				Description:  "Short Name of the alertsource",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_key": {
				Description: "Service id for which the incidents should be assigned",
				Type:        schema.TypeString,
				Required:    true,
			},
			"webhook_url": {
				Description: "url endpoint where the alerts needs to be sent to",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"version": {
				Description:  "version of the api",
				Type:         schema.TypeString,
				Default:      "v1",
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 3),
			},
		},
	}
}

func dataSourceAlertEndpointRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	shortName, ok := d.GetOk("alert_source_name")
	if !ok {
		return diag.Errorf("invalid alert source name name provided")
	}

	serviceKey, ok := d.GetOk("service_key")
	if !ok {
		return diag.Errorf("invalid service id provided")
	}

	version, ok := d.GetOk("version")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	webhookurl := fmt.Sprintf("%s/%s/incidents/%s/%s", client.BaseIngestionURL, version, shortName, serviceKey)
	err := d.Set("webhook_url", webhookurl)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(webhookurl)
	return nil
}
