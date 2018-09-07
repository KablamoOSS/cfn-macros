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

	finalResources := map[string]interface{}{}
	finalParameters := map[string]interface{}{}

	if parameters, ok := fragment["Parameters"].(map[string]interface{}); ok {
		finalParameters = parameters
	}

	if resources, ok := fragment["Resources"].(map[string]interface{}); ok {
		for key, value := range resources {
			pending = append(pending, kvpair{name: key, body: value.(map[string]interface{})})
		}

		for len(pending) > 0 {
			item := pending[len(pending)-1]
			if len(pending) == 1 {
				pending = []kvpair{}
			} else {
				pending = pending[:len(pending)-1]
			}
			this, newRs, newPs := transformResource(item.body)
			finalResources[item.name] = this
			for _, newParam := range newPs {
				finalParameters[newParam.name] = newParam.body
			}
			pending = append(pending, newRs...)
		}
		fragment["Resources"] = finalResources
		fragment["Parameters"] = finalParameters
	}

	return fragment
}

func transformResource(resource map[string]interface{}) (map[string]interface{}, []kvpair, []kvpair) {
	newResources := make([]kvpair, 0)
	newParameters := make([]kvpair, 0)

	propsIface, ok := resource["Properties"]
	if !ok {
		return resource, []kvpair{}, []kvpair{}
	}

	props, ok := propsIface.(map[string]interface{})
	if !ok {
		return resource, []kvpair{}, []kvpair{}
	}

	for key, value := range props {
		if node, ok := value.(map[string]interface{}); ok {
			kind, _ := node["Kind"].(string)
			switch kind {
			case "Resource":
				name := node["Name"].(string)
				delete(node, "Kind")
				delete(node, "Name")
				newResources = append(newResources, kvpair{name: name, body: node})
				props[key] = map[string]interface{}{"Ref": name}
			case "Parameter":
				name := node["Name"].(string)
				delete(node, "Kind")
				delete(node, "Name")
				newParameters = append(newParameters, kvpair{name: name, body: node})
				props[key] = map[string]interface{}{"Ref": name}
			}
		} else if properties, ok := value.([]interface{}); ok {
			for i, property := range properties {
				if node, ok := property.(map[string]interface{}); ok {
					kind, _ := node["Kind"].(string)
					switch kind {
					case "Resource":
						name := node["Name"].(string)
						delete(node, "Name")
						delete(node, "Kind")
						newResources = append(newResources, kvpair{name: name, body: node})
						properties[i] = map[string]interface{}{"Ref": name}
					case "Parameter":
						name := node["Name"].(string)
						delete(node, "Name")
						delete(node, "Kind")
						newParameters = append(newParameters, kvpair{name: name, body: node})
						properties[i] = map[string]interface{}{"Ref": name}
					}
				}
			}
		}
	}
	resource["Properties"] = props

	return resource, newResources, newParameters
}
