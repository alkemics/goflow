package outputs

type bindings []string

func (b *bindings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multipleBindings []string
	if err := unmarshal(&multipleBindings); err == nil {
		*b = multipleBindings
		return nil
	}

	var simpleBinding string
	if err := unmarshal(&simpleBinding); err != nil {
		return err
	}

	*b = []string{simpleBinding}
	return nil
}
