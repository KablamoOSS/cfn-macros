package transform

import (
	"fmt"
)

type Transformer struct {
	Transforms map[string]transformation
}

type transformation interface {
	apply(map[string]interface{}) (map[string]interface{}, error)
}

func (t *Transformer) Transform(tmpl map[string]interface{}) (map[string]interface{}, error) {
	transformsKey, ok := tmpl["Transform"]
	if !ok {
		// No Transform key, nothing to do
		return tmpl, nil
	}

	delete(tmpl, "Transform")

	switch s := transformsKey.(type) {
	case string:
		newTmpl, err := t.apply(s, tmpl)
		if err != nil {
			return tmpl, err
		}
		tmpl = newTmpl
	case []interface{}:
		for _, sn := range s {
			newTmpl, err := t.apply(sn.(string), tmpl)
			if err != nil {
				return tmpl, err
			}
			tmpl = newTmpl
		}
	}

	return tmpl, nil
}

func (t *Transformer) apply(name string, tmpl map[string]interface{}) (map[string]interface{}, error) {
	transformer, ok := t.Transforms[name]
	if !ok {
		return tmpl, fmt.Errorf("unrecognised transform: %s", name)
	}
	return transformer.apply(tmpl)
}
