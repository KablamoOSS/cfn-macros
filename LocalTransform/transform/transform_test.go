package transform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockTransformer struct {
	calledNTimes int
}

func (t *mockTransformer) apply(tmpl map[string]interface{}) (map[string]interface{}, error) {
	t.calledNTimes++
	return tmpl, nil
}

func TestRemovesTransformsKey(t *testing.T) {
	mockTransformation := &mockTransformer{}
	tf := &Transformer{}

	tf.Transforms = make(map[string]transformation)
	tf.Transforms["MockTransform"] = mockTransformation

	input := map[string]interface{}{
		"Transform": "MockTransform",
		"Data": map[string]interface{}{
			"Foo": "Bar",
		},
	}

	output, err := tf.Transform(input)
	// It should not produce an error
	assert.Nil(t, err)

	// It should have called the registered transformation
	assert.Equal(t, 1, mockTransformation.calledNTimes)

	_, ok := output["Transform"]
	// It should not include the "Transform" key in the output
	assert.False(t, ok)
}
