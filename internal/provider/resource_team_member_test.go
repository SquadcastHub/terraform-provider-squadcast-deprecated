package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceTeamMember(t *testing.T) {
	resourceName := "squadcast_team_member.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckTeamMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamMemberConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "629a2f542e0a6e82f408f280"),
					resource.TestCheckResourceAttr(resourceName, "user_id", "5f8891527f735f0a6646f3b7"),
					resource.TestCheckResourceAttr(resourceName, "role_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "role_ids.0", "629a2f542e0a6e82f408f281"),
					resource.TestCheckResourceAttr(resourceName, "role_ids.1", "629a2f542e0a6e82f408f282"),
				),
			},
			{
				Config: testAccResourceTeamMemberConfig_observer(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "629a2f542e0a6e82f408f280"),
					resource.TestCheckResourceAttr(resourceName, "user_id", "5eb26b36ec9f070550204c85"),
					resource.TestCheckResourceAttr(resourceName, "role_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_ids.0", "629a2f542e0a6e82f408f284"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "629a2f542e0a6e82f408f280:diane@squadcast.com",
			},
		},
	})
}

func testAccCheckTeamMemberDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_team_member" {
			continue
		}

		_, err := client.GetTeamMemberByID(context.Background(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected member to be deleted, but was found")
		}

	}

	return nil
}

func testAccResourceTeamMemberConfig() string {
	return fmt.Sprintf(`
resource "squadcast_team_member" "test" {
	team_id = "629a2f542e0a6e82f408f280"
	user_id = "5f8891527f735f0a6646f3b7"
	role_ids = ["629a2f542e0a6e82f408f281", "629a2f542e0a6e82f408f282"]
}
	`)
}

func testAccResourceTeamMemberConfig_observer() string {
	return fmt.Sprintf(`
resource "squadcast_team_member" "test" {
	team_id = "629a2f542e0a6e82f408f280"
	user_id = "5eb26b36ec9f070550204c85"
	role_ids = ["629a2f542e0a6e82f408f284"]
}
	`)
}
