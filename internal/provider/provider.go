package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-squadcast/internal/types"
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
				"squadcast_data_source": dataSourceOne(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"squadcast_resource": resourceTwo(),
			},
			Schema: map[string]*schema.Schema{
				"region": {
					Type:         schema.TypeString,
					Optional:     true,
					DefaultFunc:  schema.EnvDefaultFunc("SQUADCAST_REGION", "us"),
					ValidateFunc: validation.StringInSlice([]string{"us", "eu", "internal", "staging", "dev"}, false),
				},
				"refresh_token": {
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

type apiClient struct {
	Host   string
	Region string

	RefreshToken string
	AccessToken  string

	UserAgent string
	BaseURL   string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, rd *schema.ResourceData) (c interface{}, diags diag.Diagnostics) {
		client := &apiClient{}
		client.UserAgent = p.UserAgent("terraform-provider-squadcast", version)

		region := rd.Get("region").(string)
		refreshToken := rd.Get("refresh_token").(string)

		client.RefreshToken = refreshToken

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
			client.BaseURL = fmt.Sprintf("http://%s:8081/v3", client.Host)
		} else {
			client.BaseURL = fmt.Sprintf("https://api.%s/v3", client.Host)
		}

		err := client.getAccessToken(ctx)
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

func (client *apiClient) getAccessToken(ctx context.Context) error {
	path := "/oauth/access-token"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.BaseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Refresh-Token", client.RefreshToken)
	req.Header.Set("User-Agent", client.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var response struct {
		Data types.AccessToken `json:"data"`
		*types.Meta
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &response); err != nil {
		fmt.Printf("\n%#v %#v %#v", bytes, err, resp)
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(response.Meta.Meta.Message)
	}

	client.AccessToken = response.Data.AccessToken
	return nil
}
