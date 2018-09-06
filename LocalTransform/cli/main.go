package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/KablamoOSS/cfn-macros/LocalTransform/transform"
)

func main() {
	var compact bool
	flag.BoolVar(&compact, "c", false, "Produce compact JSON output")
	flag.Parse()

	var err error

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: cfn-transform <config.json> [input-template.json]\n")
		os.Exit(1)
	}

	transformer := &transform.Transformer{}

	err = parseConfig(transformer, flag.Arg(0))
	if err != nil {
		panic(err)
	}

	reader := os.Stdin

	if flag.Arg(1) != "" {
		reader, err = os.Open(flag.Arg(1))
		if err != nil {
			panic(err)
		}
	}

	err = apply(transformer, reader, os.Stdout, compact)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to apply transform: %v", err)
		os.Exit(1)
	}
}

func apply(transformer *transform.Transformer, in io.Reader, out io.Writer, compact bool) error {
	inTmpl := make(map[string]interface{})

	decoder := json.NewDecoder(in)
	err := decoder.Decode(&inTmpl)
	if err != nil {
		return err
	}

	outTmpl, err := transformer.Transform(inTmpl)
	if err != nil {
		panic(err)
	}

	encoder := json.NewEncoder(out)
	if !compact {
		encoder.SetIndent("", "  ")
	}
	err = encoder.Encode(outTmpl)
	if err != nil {
		return err
	}

	return nil
}
