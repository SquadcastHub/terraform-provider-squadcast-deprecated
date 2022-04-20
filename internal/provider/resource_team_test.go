package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTeam(t *testing.T) {
	rName := acctest.RandomWithPrefix("teams")

	resourceName := "squadcast_team.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "My Team"),
				),
			},
		},
	})
}

func testAccResourceTeamConfig(rName string) string {
	return fmt.Sprintf(`
resource "squadcast_team" "test" {
	name = "My Team"
	description = "A great set of people."

	roles {
		name = "Foo"
	}
}
	`)
}
