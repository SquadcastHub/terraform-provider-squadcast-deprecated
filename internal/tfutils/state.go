package tfutils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SetState(d *schema.ResourceData, m map[string]interface{}) error {
	id, ok := m["id"].(string)
	if !ok {
		return fmt.Errorf("invalid id")
	}

	d.SetId(id)

	for k, v := range m {
		if k == "id" {
			continue
		}

		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("cannot set `%s: %v` : %w", k, v, err)
		}
	}

	return nil
}

func EncodeAndSet(input StateEncoder, d *schema.ResourceData) error {
	m, err := input.Encode()
	if err != nil {
		return err
	}
	return SetState(d, m)
}
