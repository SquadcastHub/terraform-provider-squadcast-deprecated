package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSlo(t *testing.T) {
	sloName := acctest.RandomWithPrefix("slo")

	resourceName := "squadcast_slo.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		// CheckDestroy:      testAccCheckSloDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSloConfig(sloName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "owner_id", "611262fcd5b4ea846b534a8a"),
					resource.TestCheckResourceAttr(resourceName, "name", sloName),
					resource.TestCheckResourceAttr(resourceName, "duration_in_days", "30"),
					// TODO: Add more attributes for monitoring checks and actions
				),
			},
			{
				Config: testAccResourceSloConfig_update(sloName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "owner_id", "611262fcd5b4ea846b534a8a"),
					resource.TestCheckResourceAttr(resourceName, "name", sloName),
					resource.TestCheckResourceAttr(resourceName, "duration_in_days", "7"),
				),
			},
			// {
			// 	ResourceName:        resourceName,
			// 	ImportState:         true,
			// 	ImportStateVerify:   true,
			// 	ImportStateIdPrefix: "613611c1eb22db455cfa789f:",
			// },
		},
	})
}

// func testAccCheckSloDestroy(s *terraform.State) error {
// 	client := testProvider.Meta().(*api.Client)

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "squadcast_slo" {
// 			continue
// 		}

// 		_, err := client.GetServiceById(context.Background(), rs.Primary.Attributes["owner_id"], rs.Primary.ID)
// 		if err == nil {
// 			return fmt.Errorf("expected service to be destroyed, %s found", rs.Primary.ID)
// 		}

// 		// FIXME: check for 404 errors, any other error is not acceptable.
// 		// if !err.IsNotFoundError() {
// 		// 	return err
// 		// }
// 	}

// 	return nil
// }

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
