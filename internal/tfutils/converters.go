package tfutils

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SetToSlice[T any](set interface{}) []T {
	islice := set.(*schema.Set).List()
	slicelen := len(islice)
	slice := make([]T, slicelen, slicelen)

	if slicelen == 0 {
		return slice
	}

	for i, v := range islice {
		slice[i] = v.(T)
	}

	return slice
}

func ListToSlice[T any](list interface{}) []T {
	islice := list.([]interface{})
	slicelen := len(islice)
	slice := make([]T, slicelen, slicelen)

	if slicelen == 0 {
		return slice
	}

	for i, v := range islice {
		slice[i] = v.(T)
	}

	return slice
}
