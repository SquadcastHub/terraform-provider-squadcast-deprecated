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

func TestAccResourceSlo(t *testing.T) {
	sloName := acctest.RandomWithPrefix("terraform-acc-test-slo-")

	resourceName := "squadcast_slo.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSloDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSloConfig(sloName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "owner_id", "611262fcd5b4ea846b534a8a"),
					resource.TestCheckResourceAttr(resourceName, "name", sloName),
					resource.TestCheckResourceAttr(resourceName, "duration_in_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "notify.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.users.0", "5e1c2309342445001180f9c2"),
				),
			},
			{
				Config: testAccResourceSloConfig_update(sloName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "owner_id", "611262fcd5b4ea846b534a8a"),
					resource.TestCheckResourceAttr(resourceName, "name", sloName),
					resource.TestCheckResourceAttr(resourceName, "duration_in_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "notify.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.users.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.users.0", "5e1c2309342445001180f9c2"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.users.1", "617793e650d38001057faaaf"),
				),
			},
		},
	})
}

func testAccCheckSloDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*api.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_slo" {
			continue
		}

		slo, _ := client.GetSlo(context.Background(), rs.Primary.Attributes["org_id"], rs.Primary.Attributes["id"])
		if slo != nil {
			return fmt.Errorf("expected slo to be destroyed, %s found", slo.Name)
		}
	}
	return nil
}

func testAccResourceSloConfig(sloName string) string {
	return fmt.Sprintf(`

resource "squadcast_slo" "test" {
	name = "%s"
	description = "Tracks some slo for some service"
	target_slo = 99.9
	service_ids = ["615d3e23aff6885f46d291be"]
	slis = ["latency"]
	time_interval_type = "rolling"
	duration_in_days = 30
	org_id = "604592dabc35ea0008bb0584"

	rules {
		name = "breached_error_budget"
	}
	
	rules {
		name = "remaining_error_budget"
		threshold = 10
	}

	notify {
		users = ["5e1c2309342445001180f9c2"]
	}
	
	owner_type="team"
	owner_id = "611262fcd5b4ea846b534a8a"
}
	`, sloName)
}

func testAccResourceSloConfig_update(sloName string) string {
	return fmt.Sprintf(`

resource "squadcast_slo" "test" {
	name = "%s"
	description = "Tracks some slo for some service"
	target_slo = 99.9
	service_ids = ["615d3e23aff6885f46d291be"]
	slis = ["latency"]
	time_interval_type = "rolling"
	duration_in_days = 7
	org_id = "604592dabc35ea0008bb0584"

	rules {
		name = "breached_error_budget"
	}
	
	rules {
		name = "remaining_error_budget"
		threshold = 11
	}

	rules {
		name = "unhealthy_slo"
		threshold = 1
	}
	
	notify {
		users = ["5e1c2309342445001180f9c2", "617793e650d38001057faaaf"]
	}

	owner_type="team"
	owner_id = "611262fcd5b4ea846b534a8a"
}
	`, sloName)
}
