package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"squadcast_squad":   dataSourceSquad(),
				"squadcast_service": dataSourceService(),
				"squadcast_team":    dataSourceTeam(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"squadcast_squad":               resourceSquad(),
				"squadcast_suppression_rules":   resourceSuppressionRules(),
				"squadcast_deduplication_rules": resourceDeduplicationRules(),
				"squadcast_routing_rules":       resourceRoutingRules(),
				"squadcast_tagging_rules":       resourceTaggingRules(),
				"squadcast_service":             resourceService(),
				"squadcast_service_maintenance": resourceServiceMaintenance(),
				"squadcast_team":                resourceTeam(),
				"squadcast_team_role":           resourceTeamRole(),
				"squadcast_team_members":        resourceTeamMembers(),
			},
			Schema: map[string]*schema.Schema{
				"organization_id": {
					Description:  "org id",
					Type:         schema.TypeString,
					Required:     true,
					DefaultFunc:  schema.EnvDefaultFunc("SQUADCAST_ORGANIZATION_ID", ""),
					ValidateFunc: tfutils.ValidateObjectID,
				},
				"region": {
					Description:  "region",
					Type:         schema.TypeString,
					Optional:     true,
					DefaultFunc:  schema.EnvDefaultFunc("SQUADCAST_REGION", "us"),
					ValidateFunc: validation.StringInSlice([]string{"us", "eu", "internal", "staging", "dev"}, false),
				},
				"refresh_token": {
					Description: "refresh token",
					Type:        schema.TypeString,
					Sensitive:   true,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SQUADCAST_REFRESH_TOKEN", nil),
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, rd *schema.ResourceData) (c interface{}, diags diag.Diagnostics) {
		client := &api.Client{}
		client.UserAgent = p.UserAgent("terraform-provider-squadcast", version)

		region := rd.Get("region").(string)
		refreshToken := rd.Get("refresh_token").(string)
		organizationID := rd.Get("organization_id").(string)

		client.RefreshToken = refreshToken
		client.OrganizationID = organizationID

		switch region {
		case "us":
			client.Host = "squadcast.com"
		case "eu":
			client.Host = "eu.squadcast.com"
		case "internal":
			client.Host = "squadcast.xyz"
		case "staging":
			client.Host = "squadcast.tech"
		case "dev":
			client.Host = "localhost"
		}

		if region == "dev" {
			client.BaseURLV3 = fmt.Sprintf("http://%s:8081/v3", client.Host)
			client.BaseURLV2 = fmt.Sprintf("http://%s:8080/v2", client.Host)
			client.AuthBaseURL = fmt.Sprintf("http://%s:8081/v3", client.Host)
		} else {
			client.BaseURLV3 = fmt.Sprintf("https://api.%s/v3", client.Host)
			client.BaseURLV2 = fmt.Sprintf("https://platform-backend.%s/v2", client.Host)
			client.AuthBaseURL = fmt.Sprintf("https://api.%s/v3", client.Host)
		}

		err := client.GetAccessToken(ctx)
		if err != nil {
			return nil, append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred while fetching the access token.",
				Detail:   err.Error(),
			})
		}

		return client, nil
	}
}
