package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
)

func TestAccResourceTeamMembers(t *testing.T) {
	resourceName := "squadcast_team_members.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckTeamMembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamMembersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "629a2f542e0a6e82f408f280"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "members.0.user_id", "5ef5de4259c32c7ca25b0bfa"),
					resource.TestCheckResourceAttr(resourceName, "members.0.roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.0.roles.0", "Manage Team"),
					resource.TestCheckResourceAttr(resourceName, "members.0.roles.1", "Admin"),
				),
			},
			{
				Config: testAccResourceTeamMembersConfig_addMember(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "629a2f542e0a6e82f408f280"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.0.user_id", "5ef5de4259c32c7ca25b0bfa"),
					resource.TestCheckResourceAttr(resourceName, "members.0.roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.0.roles.0", "Manage Team"),
					resource.TestCheckResourceAttr(resourceName, "members.0.roles.1", "Admin"),
					resource.TestCheckResourceAttr(resourceName, "members.1.user_id", "5eb26b36ec9f070550204c85"),
					resource.TestCheckResourceAttr(resourceName, "members.1.roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "members.1.roles.0", "Observer"),
				),
			},
			{
				Config: testAccResourceTeamMembersConfig_removeMember(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "629a2f542e0a6e82f408f280"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "members.0.user_id", "5eb26b36ec9f070550204c85"),
					resource.TestCheckResourceAttr(resourceName, "members.0.roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.0.roles.0", "Manage Team"),
					resource.TestCheckResourceAttr(resourceName, "members.0.roles.1", "Observer"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "629a2f542e0a6e82f408f280",
			},
		},
	})
}

func testAccCheckTeamMembersDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_team_members" {
			continue
		}

		teamMembers, err := client.GetTeamById(context.Background(), rs.Primary.Attributes["team_id"])
		if err != nil {
			return err
		}

		count := len(teamMembers.Members)
		if count != 1 {
			return fmt.Errorf("expected only team admin to be present, %d found", count)
		}

		count = len(teamMembers.Members[0].RoleIDs)
		if count != 1 {
			return fmt.Errorf("expected team admin to not have other roles, %d found", count)
		}
	}

	return nil
}

func testAccResourceTeamMembersConfig() string {
	return fmt.Sprintf(`
resource "squadcast_team_members" "test" {
	team_id = "629a2f542e0a6e82f408f280"

	members {
		user_id = "5ef5de4259c32c7ca25b0bfa"
		roles = ["Manage Team", "Admin"]
	}
}
	`)
}

func testAccResourceTeamMembersConfig_addMember() string {
	return fmt.Sprintf(`
resource "squadcast_team_members" "test" {
	team_id = "629a2f542e0a6e82f408f280"

	members {
		user_id = "5ef5de4259c32c7ca25b0bfa"
		roles = ["Manage Team", "Admin"]
	}

	members {
		user_id = "5eb26b36ec9f070550204c85"
		roles = ["Observer"]
	}
}
	`)
}

func testAccResourceTeamMembersConfig_removeMember() string {
	return fmt.Sprintf(`
resource "squadcast_team_members" "test" {
	team_id = "629a2f542e0a6e82f408f280"

	members {
		user_id = "5eb26b36ec9f070550204c85"
		roles = ["Manage Team", "Observer"]
	}
}
	`)
}
