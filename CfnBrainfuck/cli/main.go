package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/KablamoOSS/cfn-macros/CfnBrainfuck/transform"
)

func main() {
	flag.Parse()

	var err error
	reader := os.Stdin

	if flag.Arg(0) != "" {
		reader, err = os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}
	}

	err = apply(reader, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to apply transform: %v", err)
		os.Exit(1)
	}
}

func apply(in io.Reader, out io.Writer) error {
	inTmpl := make(map[string]interface{})

	decoder := json.NewDecoder(in)
	err := decoder.Decode(&inTmpl)
	if err != nil {
		return err
	}

	output := transform.Transform(inTmpl)
	fmt.Print(string(output))

	return nil
}
