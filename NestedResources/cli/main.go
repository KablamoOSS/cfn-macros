package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/KablamoOSS/cfn-macros/NestedResources/transform"
)

func main() {
	var compact bool
	flag.BoolVar(&compact, "c", false, "Produce compact JSON output")
	flag.Parse()

	var err error
	reader := os.Stdin

	if flag.Arg(0) != "" {
		reader, err = os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}
	}

	err = apply(reader, os.Stdout, compact)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to apply transform: %v", err)
		os.Exit(1)
	}
}

func apply(in io.Reader, out io.Writer, compact bool) error {
	inTmpl := make(map[string]interface{})

	decoder := json.NewDecoder(in)
	err := decoder.Decode(&inTmpl)
	if err != nil {
		return err
	}

	outTmpl := transform.Transform(inTmpl)

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
