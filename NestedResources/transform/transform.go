package transform

type kvpair struct {
	name string
	body map[string]interface{}
}

func Transform(fragment map[string]interface{}) map[string]interface{} {
	if _, ok := fragment["Resources"]; !ok {
		return fragment
	}

	pending := []kvpair{}

	if resources, ok := fragment["Resources"].(map[string]interface{}); ok {
		for key, value := range resources {
			pending = append(pending, kvpair{name: key, body: value.(map[string]interface{})})
		}

		finalResources := map[string]interface{}{}
		for len(pending) > 0 {
			item := pending[len(pending)-1]
			if len(pending) == 1 {
				pending = []kvpair{}
			} else {
				pending = pending[:len(pending)-1]
			}
			this, adds := transformResource(item.body)
			finalResources[item.name] = this
			pending = append(pending, adds...)
		}
		fragment["Resources"] = finalResources
	}

	return fragment
}

func transformResource(resource map[string]interface{}) (map[string]interface{}, []kvpair) {
	adds := make([]kvpair, 0)

	propsIface, ok := resource["Properties"]
	if !ok {
		return resource, []kvpair{}
	}

	props, ok := propsIface.(map[string]interface{})
	if !ok {
		return resource, []kvpair{}
	}

	for key, value := range props {
		if resource, ok := value.(map[string]interface{}); ok && looksLikeResource(resource) {
			name := resource["Name"].(string)
			delete(resource, "Name")
			adds = append(adds, kvpair{name: name, body: resource})
			props[key] = map[string]interface{}{"Ref": name}
		} else if properties, ok := value.([]interface{}); ok {
			for i, property := range properties {
				if resource, ok := property.(map[string]interface{}); ok && looksLikeResource(resource) {
					name := resource["Name"].(string)
					delete(resource, "Name")
					adds = append(adds, kvpair{name: name, body: resource})
					properties[i] = map[string]interface{}{"Ref": name}
				}
			}
		}
	}
	resource["Properties"] = props

	return resource, adds
}

func looksLikeResource(resource map[string]interface{}) bool {
	_, typeIsString := resource["Type"].(string)
	_, propsIsObject := resource["Properties"].(map[string]interface{})
	return typeIsString && propsIsObject
}
