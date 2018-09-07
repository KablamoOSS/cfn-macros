package transform

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ingestorJson = `{
	"Parameters": {
		"FooParameter": {
			"Type": "String"
		}
	},
	"Resources": {
		"FooResource": {
			"Type": "CI::Example",
			"Properties": {}
		}
	},
	"Outputs": {
		"FooOutput": {
			"Value": { "Ref": "FooResource" }
		}
	}
}`

var ingesteeJson = `{
	"Parameters": {
		"BarParameter": {
			"Type": "String"
		}
	},
	"Resources": {
		"BarResource": {
			"Type": "CI::Example",
			"Properties": {
			}
		}
	},
	"Outputs": {
		"BarOutput": {
			"Value": { "Ref": "BarResource" }
		}
	}
}`

// Simple ingestion - result should include the full set of parameters/resources/outputs
func TestBasicIngestion(t *testing.T) {
	ingestorTmpl := &CfnTemplate{}
	ingesteeTmpl := &CfnTemplate{}

	err := json.Unmarshal([]byte(ingestorJson), ingestorTmpl)
	assert.Nil(t, err)

	err = json.Unmarshal([]byte(ingesteeJson), ingesteeTmpl)
	assert.Nil(t, err)

	err = ingestorTmpl.ingest("Ingestee", ingesteeTmpl, map[string]interface{}{})
	assert.Nil(t, err)

	assert.Len(t, ingestorTmpl.Parameters, 2)
	assert.Len(t, ingestorTmpl.Resources, 2)
	assert.Len(t, ingestorTmpl.Outputs, 2)
}

var conflictJson = `{
	"Resources": {
		"FooResource": {
			"Type": "CI::Example",
			"Properties": {}
		}
	}
}`

// Ingesting a template with conflicting parameter/resource/etc names should return an error
func TestConflictIngestion(t *testing.T) {
	ingestorTmpl := &CfnTemplate{}
	ingesteeTmpl := &CfnTemplate{}

	err := json.Unmarshal([]byte(ingestorJson), ingestorTmpl)
	assert.Nil(t, err)

	err = json.Unmarshal([]byte(conflictJson), ingesteeTmpl)
	assert.Nil(t, err)

	err = ingestorTmpl.ingest("Ingestee", ingesteeTmpl, map[string]interface{}{})
	assert.NotNil(t, err)
}

var refIngestorJson = `{
	"Parameters": {
		"FooParameter": {
			"Type": "String"
		}
	},
	"Resources": {
		"FooResource": {
			"Type": "CI::Example",
			"Properties": {
				"Value": { "Ingest::Ref": "Ingestee.BarOutput" }
			}
		}
	},
	"Outputs": {
		"FooOutput": {
			"Value": { "Ref": "FooResource" }
		}
	}
}`

var refIngesteeJson = `{
	"Parameters": {
		"BarParameter": {
			"Type": "String"
		}
	},
	"Resources": {
		"BarResource": {
			"Type": "CI::Example",
			"Properties": {
				"Value": { "Ref": "BarParameter" }
			}
		}
	},
	"Outputs": {
		"BarOutput": {
			"Value": "ABC"
		}
	}
}`

// Simple ingestion - result should include the full set of parameters+resources
// default behaviour doesn't include outputs, so only one expected
func TestRefIngestion(t *testing.T) {
	ingestorTmpl := &CfnTemplate{}
	ingesteeTmpl := &CfnTemplate{}

	err := json.Unmarshal([]byte(refIngestorJson), ingestorTmpl)
	assert.Nil(t, err)

	err = json.Unmarshal([]byte(refIngesteeJson), ingesteeTmpl)
	assert.Nil(t, err)

	err = ingestorTmpl.ingest("Ingestee", ingesteeTmpl, map[string]interface{}{"BarParameter": "XYZ"})
	assert.Nil(t, err)

	assert.Len(t, ingestorTmpl.Parameters, 1)
	assert.Len(t, ingestorTmpl.Resources, 2)
	assert.Len(t, ingestorTmpl.Outputs, 2)

	barResource, _ := ingestorTmpl.Resources["BarResource"].(map[string]interface{})
	barProperties := barResource["Properties"].(map[string]interface{})
	barValue, barValueIsString := barProperties["Value"].(string)
	assert.True(t, barValueIsString, "BarResource.Properties: %#v", barResource)
	assert.Equal(t, "XYZ", barValue)

	fooResource, _ := ingestorTmpl.Resources["FooResource"].(map[string]interface{})
	fooProperties := fooResource["Properties"].(map[string]interface{})
	fooValue, fooValueIsString := fooProperties["Value"].(string)
	assert.True(t, fooValueIsString, "FooResource.Properties: %#v", fooResource)
	assert.Equal(t, "ABC", fooValue)
}
