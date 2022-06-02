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
				// TODO: Update email and test data while setting up tests
				Config: testAccUserDataSourceConfig("gurunandan@squadcast.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "org_role_name", "account_owner"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig(email string) string {
	return fmt.Sprintf(`
	data "squadcast_user" "test" {
		email = "%s"
	}`, email)
}
