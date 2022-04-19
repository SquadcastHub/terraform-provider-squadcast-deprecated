package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTeams(t *testing.T) {
	rName := acctest.RandomWithPrefix("teams")

	resourceName := "data.squadcast_teams.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamsDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttr(resourceName, "teams.*", "3"),
				),
			},
		},
	})
}

func testAccTeamsDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
data "squadcast_teams" "test" {}
	`)
}
