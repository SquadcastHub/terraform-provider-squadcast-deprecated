package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

const teamMembersID = "team_members"

func resourceTeamMembers() *schema.Resource {
	return &schema.Resource{
		Description: "TeamMembers resource.",

		CreateContext: resourceTeamMembersCreate,
		ReadContext:   resourceTeamMembersRead,
		UpdateContext: resourceTeamMembersUpdate,
		DeleteContext: resourceTeamMembersDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTeamMembersImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tfutils.ValidateObjectID,
				ForceNew:     true,
			},
			"member": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Description:  "user id?.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tfutils.ValidateObjectID,
						},
						"roles": {
							Description: "role names.",
							Type:        schema.TypeList,
							Required:    true,
							MinItems:    1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceTeamMembersImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	id := d.Id()

	_, err := client.GetTeamById(ctx, id)
	if err != nil {
		return nil, err
	}

	d.Set("team_id", id)
	d.SetId(teamMembersID)

	return []*schema.ResourceData{d}, nil
}

func decodeTeamMembers(team *api.Team, memberslist interface{}) ([]*api.TeamMember, map[string]interface{}, error) {
	roleIDsMap := map[string]interface{}{}
	for _, r := range team.Roles {
		roleIDsMap[r.Name] = r.ID
	}

	tfmembers := memberslist.([]interface{})
	members := make([]*api.TeamMember, 0, len(tfmembers))
	for _, mem := range tfmembers {
		roleNames := tfutils.ListToSlice[string](mem.(map[string]interface{})["roles"].([]interface{}))
		roleIDs := make([]string, 0, len(roleNames))
		for _, roleName := range roleNames {
			roleID, ok := roleIDsMap[roleName]
			if !ok {
				return nil, nil, fmt.Errorf("cannot find role %s in team %s", roleName, team.Name)
			}
			roleIDs = append(roleIDs, roleID.(string))
		}

		member := &api.TeamMember{
			UserID:  mem.(map[string]interface{})["user_id"].(string),
			RoleIDs: roleIDs,
		}

		members = append(members, member)
	}

	return members, roleIDsMap, nil
}

func resourceTeamMembersCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	team, err := client.GetTeamById(ctx, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	members, _, err := decodeTeamMembers(team, d.Get("member"))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateTeamMembers(ctx, d.Get("team_id").(string), &api.UpdateTeamMembersReq{Members: members})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(teamMembersID)

	return resourceTeamMembersRead(ctx, d, meta)
}

func resourceTeamMembersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	team, err := client.GetTeamById(ctx, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	roleNamesMap := map[string]interface{}{}
	for _, r := range team.Roles {
		roleNamesMap[r.ID] = r.Name
	}

	members := make([]map[string]interface{}, 0, len(team.Members))
	for _, mem := range team.Members {
		roleNames := make([]interface{}, 0, len(mem.RoleIDs))
		for _, rid := range mem.RoleIDs {
			roleNames = append(roleNames, roleNamesMap[rid])
		}

		member := map[string]interface{}{
			"user_id": mem.UserID,
			"roles":   roleNames,
		}

		members = append(members, member)
	}

	d.Set("member", members)

	return nil
}

func resourceTeamMembersUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return resourceTeamMembersCreate(ctx, d, meta)
}

func resourceTeamMembersDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	team, err := client.GetTeamById(ctx, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	members, roleIDsMap, err := decodeTeamMembers(team, d.Get("member"))
	if err != nil {
		return diag.FromErr(err)
	}

	manageTeamRoleID := roleIDsMap["Manage Team"].(string)

	var manager *api.TeamMember
	for _, mem := range members {
		for _, rid := range mem.RoleIDs {
			if rid == manageTeamRoleID {
				manager = &api.TeamMember{
					UserID:  mem.UserID,
					RoleIDs: []string{manageTeamRoleID},
				}
			}
		}

		if manager != nil {
			break
		}
	}

	_, err = client.UpdateTeamMembers(ctx, d.Get("team_id").(string), &api.UpdateTeamMembersReq{Members: []*api.TeamMember{manager}})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
