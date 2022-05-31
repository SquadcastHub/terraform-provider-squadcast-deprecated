package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-squadcast/internal/api"
)

func TestAccResourceTaggingRules(t *testing.T) {
	resourceName := "squadcast_tagging_rules.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckTaggingRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTaggingRulesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.is_basic", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "payload[\"event_id\"] == 40"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.key", "MyTag"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.value", "foo"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.color", "#ababab"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
				),
			},
			{
				Config: testAccResourceTaggingRulesConfig_updateRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.is_basic", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "payload[\"event_id\"] == 40"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.basic_expression.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.key", "MyTag"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.value", "foo"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.color", "#ababab"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.is_basic", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.expression", ""),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expression.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expression.0.lhs", "payload[\"foo\"]"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expression.0.op", "is"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expression.0.rhs", "bar"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.0.key", "MyTag"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.0.value", "foo"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.0.color", "#ababab"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.1.key", "MyTag2"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.1.value", "bar"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.1.color", "#f0f0f0"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "613611c1eb22db455cfa789f:61361611c2fc70c3101ca7dd",
			},
		},
	})
}

func testAccCheckTaggingRulesDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_tagging_rules" {
			continue
		}

		taggingRules, err := client.GetTaggingRules(context.Background(), rs.Primary.Attributes["service_id"], rs.Primary.Attributes["team_id"])
		if err != nil {
			return err
		}
		count := len(taggingRules.Rules)
		if count > 0 {
			return fmt.Errorf("expected all tagging rules to be destroyed, %d found", count)
		}
	}

	return nil
}

func testAccResourceTaggingRulesConfig() string {
	return fmt.Sprintf(`
resource "squadcast_tagging_rules" "test" {
	team_id = "613611c1eb22db455cfa789f"
	service_id = "61361611c2fc70c3101ca7dd"

	rules {
		is_basic = false
		expression = "payload[\"event_id\"] == 40"

		tags {
			key = "MyTag"
			value = "foo"
			color = "#ababab"
		}
	}
}
	`)
}

func testAccResourceTaggingRulesConfig_updateRules() string {
	return fmt.Sprintf(`
resource "squadcast_tagging_rules" "test" {
	team_id = "613611c1eb22db455cfa789f"
	service_id = "61361611c2fc70c3101ca7dd"

	rules {
		is_basic = false
		expression = "payload[\"event_id\"] == 40"

		tags {
			key = "MyTag"
			value = "foo"
			color = "#ababab"
		}
	}

	rules {
		is_basic = true

		basic_expression {
			lhs = "payload[\"foo\"]"
			op = "is"
			rhs = "bar"
		}

		tags {
			key = "MyTag"
			value = "foo"
			color = "#ababab"
		}

		tags {
			key = "MyTag2"
			value = "bar"
			color = "#f0f0f0"
		}
	}
}
	`)
}
