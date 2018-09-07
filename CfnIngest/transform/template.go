package transform

import (
	"fmt"
)

// TODO: This makes parsing JSON templates a little simpler, but reduces
// compatibility with other template transforms which might include additional
// top level sections. Replacing this with a map[string]interface{} would solve
// this, at the cost of making it a little more awkward to work with.
type CfnTemplate struct {
	Ingest     map[string]interface{} `json:"Ingest,omitempty"`
	Transforms interface{}            `json:"Transform,omitempty"`
	Parameters map[string]interface{} `json:"Parameters,omitempty"`
	Mappings   map[string]interface{} `json:"Mappings,omitempty"`
	Conditions map[string]interface{} `json:"Conditions,omitempty"`
	Resources  map[string]interface{} `json:"Resources,omitempty"`
	Outputs    map[string]interface{} `json:"Outputs,omitempty"`
}

func getObject(node map[string]interface{}, key string) (map[string]interface{}, bool) {
	if value, hasValue := node[key]; hasValue {
		if valueAsMap, valueIsMap := value.(map[string]interface{}); valueIsMap {
			return valueAsMap, true
		}
	}
	return nil, false
}

func (tmpl *CfnTemplate) Transform() error {
	for ingestName, ingestIface := range tmpl.Ingest {
		ingestMap, ok := ingestIface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Ingest %s definition is not a map", ingestName)
		}

		sourceIface, ok := ingestMap["Source"]
		if !ok {
			return fmt.Errorf("Ingest %s definitions `Source` key", ingestName)
		}

		sourcePath, ok := sourceIface.(string)
		if !ok {
			return fmt.Errorf("Ingest %s `Source` is not a string", ingestName)
		}

		paramsMap := make(map[string]interface{})
		paramsIface, ok := ingestMap["Parameters"]
		if ok {
			paramsMap, ok = paramsIface.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Ingest %s `Parameters` is not a map", ingestName)
			}
		}

		child, err := GetPath(sourcePath)
		if err != nil {
			return err
		}
		child.Transform()

		err = tmpl.ingest(ingestName, child, paramsMap)
		if err != nil {
			return err
		}

		delete(tmpl.Ingest, ingestName)
	}

	return nil
}

func (tmpl *CfnTemplate) ingest(templateName string, child *CfnTemplate, parameters map[string]interface{}) error {
	for paramKey, paramValue := range parameters {
		child.SupplyParameter(paramKey, paramValue)
	}

	for outputKey, outputIface := range child.Outputs {
		outputMap, ok := outputIface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Ingested template %s output %s is not a map", templateName, outputKey)
		}

		outputValue, ok := outputMap["Value"]
		if !ok {
			return fmt.Errorf("Ingested template %s output %s is missing a `Value` key", templateName, outputKey)
		}

		refKey := fmt.Sprintf("%s.%s", templateName, outputKey)
		inject(tmpl.Resources, "Ingest::Ref", refKey, outputValue)
		inject(tmpl.Mappings, "Ingest::Ref", refKey, outputValue)
		inject(tmpl.Conditions, "Ingest::Ref", refKey, outputValue)
		inject(tmpl.Outputs, "Ingest::Ref", refKey, outputValue)
	}

	for paramKey, paramIface := range child.Parameters {
		if _, ok := tmpl.Parameters[paramKey]; ok {
			return fmt.Errorf("Ingested template %s has a conflicting parameter named %s", templateName, paramKey)
		}
		if tmpl.Parameters == nil {
			tmpl.Parameters = make(map[string]interface{})
		}
		tmpl.Parameters[paramKey] = paramIface
	}

	for mappingKey, mappingIface := range child.Mappings {
		if _, ok := tmpl.Mappings[mappingKey]; ok {
			return fmt.Errorf("Ingested template %s has a conflicting mapping named %s", templateName, mappingKey)
		}
		if tmpl.Mappings == nil {
			tmpl.Mappings = make(map[string]interface{})
		}
		tmpl.Mappings[mappingKey] = mappingIface
	}

	for condKey, condIface := range child.Conditions {
		if _, ok := tmpl.Conditions[condKey]; ok {
			return fmt.Errorf("Ingested template %s has a conflicting condition named %s", templateName, condKey)
		}
		if tmpl.Conditions == nil {
			tmpl.Conditions = make(map[string]interface{})
		}
		tmpl.Conditions[condKey] = condIface
	}

	for resourceKey, resourceIface := range child.Resources {
		if _, ok := tmpl.Resources[resourceKey]; ok {
			return fmt.Errorf("Ingested template %s has a conflicting resource named %s", templateName, resourceKey)
		}
		if tmpl.Resources == nil {
			tmpl.Resources = make(map[string]interface{})
		}
		tmpl.Resources[resourceKey] = resourceIface
	}

	// TODO: Add a flag to skip adding these (or only add selected outputs)
	for outputKey, outputIface := range child.Outputs {
		if _, ok := tmpl.Outputs[outputKey]; ok {
			return fmt.Errorf("Ingested template %s has a conflicting output named %s", templateName, outputKey)
		}
		if tmpl.Outputs == nil {
			tmpl.Outputs = make(map[string]interface{})
		}
		tmpl.Outputs[outputKey] = outputIface
	}

	// TODO: Merge Transforms

	return nil
}

func (tmpl *CfnTemplate) SupplyParameter(paramName string, paramIface interface{}) error {
	if _, ok := tmpl.Parameters[paramName]; !ok {
		return fmt.Errorf("template does not accept parameter %s", paramName)
	}

	inject(tmpl.Resources, "Ref", paramName, paramIface)
	inject(tmpl.Mappings, "Ref", paramName, paramIface)
	inject(tmpl.Conditions, "Ref", paramName, paramIface)
	inject(tmpl.Outputs, "Ref", paramName, paramIface)
	delete(tmpl.Parameters, paramName)

	return nil
}

func inject(nodeIface interface{}, token, paramKey string, newValue interface{}) {
	switch node := nodeIface.(type) {
	case map[string]interface{}:
		for nodeKey, nodeValue := range node {
			vMap, nodeValueIsMap := nodeValue.(map[string]interface{})
			if nodeValueIsMap {
				tokenEntry, tokenEntryIsString := vMap[token].(string)
				if tokenEntryIsString && len(vMap) == 1 && tokenEntry == paramKey {
					node[nodeKey] = newValue
					return
				}
			}
			inject(nodeValue, token, paramKey, newValue)
		}
	case []interface{}:
		for i, v := range node {
			m, isMap := v.(map[string]interface{})
			if isMap {
				s, isString := m[token].(string)
				if isString && len(m) == 1 && s == paramKey {
					node[i] = newValue
					return
				}
			}
			inject(v, token, paramKey, newValue)

		}
	}
}
