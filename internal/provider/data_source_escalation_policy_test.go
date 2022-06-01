package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEscalationPolicy(t *testing.T) {

	//TODO: input escalation policy name with random name prefix
	// acctest.RandomWithPrefix("escalation-policy-")
	escalationPolicyName := "Support Team Escalation"

	resourceName := "data.squadcast_escalation_policy.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: Support this back once we we have resource to create escalation policy
		// When that happens, update the test case to create escalation policy first and later fetch that with datasource
		// CheckDestroy:      testAccCheckEscalationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyDataSourceConfig(escalationPolicyName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "625644143b923689c1015240"),
					resource.TestCheckResourceAttr(resourceName, "name", escalationPolicyName),
					resource.TestCheckResourceAttr(resourceName, "description", "some description here."),
				),
			},
		},
	})
}

func testAccEscalationPolicyDataSourceConfig(escalationPolicyName string) string {
	return fmt.Sprintf(`
data "squadcast_escalation_policy" "test" {
	name = "%s"
	team_id = "625644143b923689c1015240"
}
	`, escalationPolicyName)
}
