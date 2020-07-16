package inputs

import (
	"fmt"

	"github.com/alkemics/goflow"
)

type binding goflow.Field

func (b *binding) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var simpleBinding string
	if err := unmarshal(&simpleBinding); err == nil {
		*b = binding(goflow.Field{
			Name: simpleBinding,
			Type: "@any",
		})
		return nil
	}

	var mapBinding map[string]string
	if err := unmarshal(&mapBinding); err != nil {
		return err
	}

	if len(mapBinding) != 1 {
		var key string
		for k := range mapBinding {
			key = k
			break
		}
		return BindingError{
			Input: key,
			Err:   fmt.Errorf("cannot declare %d (!= 1) inputs", len(mapBinding)),
		}
	}

	for i, t := range mapBinding {
		switch t {
		case "single":
			*b = binding(goflow.Field{
				Name: i,
				Type: "@single",
			})
		default:
			return BindingError{
				Input: i,
				Err:   fmt.Errorf("invalid modifier '%s'", t),
			}
		}
	}

	return nil
}
