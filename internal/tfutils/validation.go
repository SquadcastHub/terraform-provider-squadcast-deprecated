package tfutils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

var ValidateObjectID = validation.StringLenBetween(24, 24)
