package tfutils

import (
	"github.com/mitchellh/mapstructure"
)

const EncoderStructTag = "tf"

type StateEncoder interface {
	Encode() (map[string]interface{}, error)
}

func EncodeGeneric(input interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &m,
		TagName: EncoderStructTag,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(input)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func EncodeSliceGeneric[T interface{}](input []T) ([]interface{}, error) {
	slice := make([]interface{}, len(input))
	for i, v := range input {
		m, err := EncodeGeneric(v)
		if err != nil {
			return nil, err
		}
		slice[i] = m
	}

	return slice, nil
}

func EncodeSlice[T StateEncoder](input []T) ([]interface{}, error) {
	slice := make([]interface{}, len(input))
	for i, v := range input {
		m, err := v.Encode()
		if err != nil {
			return nil, err
		}
		slice[i] = m
	}

	return slice, nil
}
