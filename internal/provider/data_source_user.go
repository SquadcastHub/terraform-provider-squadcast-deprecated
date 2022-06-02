package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Squadcast user",
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "User id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "User Email",
				Type:        schema.TypeString,
				Required:    true,
			},
			"first_name": {
				Description: "User First Name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_name": {
				Description: "User Last Name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"org_role_id": {
				Description: "User Org level role id, for account_owner/ user/ stakeholder",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"org_role_name": {
				Description: "User Org level role name: account_owner/ user/ stakeholder",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	email := d.Get("email")

	tflog.Info(ctx, "Reading user by email", map[string]interface{}{
		"email": email.(string),
	})
	user, err := client.GetUserByEmail(ctx, email.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tfutils.EncodeAndSet(user, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
