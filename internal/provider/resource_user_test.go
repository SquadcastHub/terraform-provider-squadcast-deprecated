package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
)

func TestAccResourceUser(t *testing.T) {
	userName := fmt.Sprintf("%s%s", "user", acctest.RandStringFromCharSet(10, "abcdefghijlkmnopqrstuvwxyz"))

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_user(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", userName),
					resource.TestCheckResourceAttr(resourceName, "last_name", "lastname"),
					resource.TestCheckResourceAttr(resourceName, "email", userName+"@example.com"),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
			{
				Config: testAccResourceUserConfig_abilities(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", userName),
					resource.TestCheckResourceAttr(resourceName, "last_name", "lastname"),
					resource.TestCheckResourceAttr(resourceName, "email", userName+"@example.com"),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "manage-billing"),
				),
			},
			{
				Config: testAccResourceUserConfig_stakeholder(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", userName),
					resource.TestCheckResourceAttr(resourceName, "last_name", "lastname"),
					resource.TestCheckResourceAttr(resourceName, "email", userName+"@example.com"),
					resource.TestCheckResourceAttr(resourceName, "role", "stakeholder"),
				),
			},
			{
				Config: testAccResourceUserConfig_user(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", userName),
					resource.TestCheckResourceAttr(resourceName, "last_name", "lastname"),
					resource.TestCheckResourceAttr(resourceName, "email", userName+"@example.com"),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     userName + "@example.com",
			},
		},
	})
}

func testAccCheckUserDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_user" {
			continue
		}

		_, err := client.GetUserById(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected user to be destroyed, %s found", rs.Primary.ID)
		}

		// FIXME: check for 404 errors, any other error is not acceptable.
		// if !err.IsNotFoundError() {
		// 	return err
		// }
	}

	return nil
}

func testAccResourceUserConfig_user(userName string) string {
	return fmt.Sprintf(`
resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "lastname"
	email = "%s@example.com"
	role = "user"
}
	`, userName, userName)
}

func testAccResourceUserConfig_abilities(userName string) string {
	return fmt.Sprintf(`
resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "lastname"
	email = "%s@example.com"
	role = "user"

	abilities = ["manage-billing"]
}
	`, userName, userName)
}

func testAccResourceUserConfig_stakeholder(userName string) string {
	return fmt.Sprintf(`
resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "lastname"
	email = "%s@example.com"
	role = "stakeholder"

	abilities = ["manage-billing"]
}
		`, userName, userName)
}
