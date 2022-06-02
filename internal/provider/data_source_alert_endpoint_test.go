package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAlertEndpoint(t *testing.T) {
	resourceName := "data.squadcast_alert_endpoint.test_uptime_robot"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertEndpointDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "webhook_url",
						"https://api.squadcast.tech/v1/incidents/uptime-robot/d5f207aabfe813e618d1496c1834ddce62cd201e"),
				),
			},
		},
	})
}

func testAccAlertEndpointDataSourceConfig() string {
	return fmt.Sprint(`
data "squadcast_alert_endpoint" "test_uptime_robot" {
	service_key = "d5f207aabfe813e618d1496c1834ddce62cd201e"
	alert_source_name = "uptime-robot"
}
	`)
}
