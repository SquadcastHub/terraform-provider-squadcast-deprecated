package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTeam(t *testing.T) {
	rName := acctest.RandomWithPrefix("team")

	resourceName := "data.squadcast_team.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "Default Team"),
					resource.TestCheckResourceAttr(resourceName, "description", "Default team"),
					resource.TestCheckResourceAttr(resourceName, "slug", "default-team"),
					resource.TestCheckResourceAttr(resourceName, "default", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
		},
	})
}

func testAccTeamDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
data "squadcast_team" "test" {
	id = "613611c1eb22db455cfa789f"
}
	`)
}
