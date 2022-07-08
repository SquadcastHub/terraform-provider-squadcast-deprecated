package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceTeamRole(t *testing.T) {
	teamRoleName := acctest.RandomWithPrefix("team_role")

	resourceName := "squadcast_team_role.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckTeamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamRoleConfig(teamRoleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "629a2f542e0a6e82f408f280"),
					resource.TestCheckResourceAttr(resourceName, "name", teamRoleName),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "read-escalation-policies"),
				),
			},
			{
				Config: testAccResourceTeamRoleConfig_update(teamRoleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "629a2f542e0a6e82f408f280"),
					resource.TestCheckResourceAttr(resourceName, "name", teamRoleName),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "read-escalation-policies"),
					resource.TestCheckResourceAttr(resourceName, "abilities.1", "update-runbooks"),
				),
			},
			{
				Config: testAccResourceTeamRoleConfig_noAbilities(teamRoleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "629a2f542e0a6e82f408f280"),
					resource.TestCheckResourceAttr(resourceName, "name", teamRoleName),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "629a2f542e0a6e82f408f280:" + teamRoleName,
			},
		},
	})
}

func testAccCheckTeamRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_team_role" {
			continue
		}

		_, err := client.GetTeamRoleByID(context.Background(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected team role to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceTeamRoleConfig(teamRoleName string) string {
	return fmt.Sprintf(`
resource "squadcast_team_role" "test" {
	name = "%s"
	team_id = "629a2f542e0a6e82f408f280"
	abilities = ["read-escalation-policies"]
}
	`, teamRoleName)
}

func testAccResourceTeamRoleConfig_update(teamRoleName string) string {
	return fmt.Sprintf(`
resource "squadcast_team_role" "test" {
	name = "%s"
	team_id = "629a2f542e0a6e82f408f280"
	abilities = ["read-escalation-policies", "update-runbooks"]
}
	`, teamRoleName)
}

func testAccResourceTeamRoleConfig_noAbilities(teamRoleName string) string {
	return fmt.Sprintf(`
resource "squadcast_team_role" "test" {
	name = "%s"
	team_id = "629a2f542e0a6e82f408f280"
	abilities = []
}
	`, teamRoleName)
}
