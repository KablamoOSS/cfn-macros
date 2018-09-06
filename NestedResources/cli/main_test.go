package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var sampleJson = `{
	"Foo": "Bar"
}`

func TestCLITransform(t *testing.T) {
	// Not really trying to test any actual transformation here, just that it
	// reads/writes json correctly.
	in := bytes.NewBufferString(sampleJson)
	out := &bytes.Buffer{}

	apply(in, out, true)

	assert.Equal(t, `{"Foo":"Bar"}`+"\n", out.String())
}
