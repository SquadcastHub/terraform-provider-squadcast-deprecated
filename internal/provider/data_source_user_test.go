package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	resourceName := "data.squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig("gurunandan@squadcast.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig(email string) string {
	fmt.Printf(`
	data "squadcast_user" "test" {
		email = "%s"
	}`, email)
	return fmt.Sprintf(`
	data "squadcast_user" "test" {
		email = "%s"
	}`, email)
}
