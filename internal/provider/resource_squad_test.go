package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSquad(t *testing.T) {
	rName := acctest.RandomWithPrefix("squad")

	resourceName := "squadcast_squad.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSquadConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "My Squad"),
				),
			},
		},
	})
}

func testAccResourceSquadConfig(rName string) string {
	return fmt.Sprintf(`
resource "squadcast_squad" "test" {
	name = "My Squad"
	team_id = "613611c1eb22db455cfa789f"
	member_ids = toset(["5f8891527f735f0a6646f3b6"])
}
	`)
}
