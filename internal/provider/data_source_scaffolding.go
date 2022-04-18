package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOne() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample data source in the Terraform provider.",

		ReadContext: dataSourceOneRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Escalation policy id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description:  "Escalation policy name.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceOneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// client := meta.(*apiClient)

	// id, ok := d.GetOk("id")
	// if ok {

	// }

	idFromAPI := "my-id"
	d.SetId(idFromAPI)

	return diag.Errorf("not implemented")
}
