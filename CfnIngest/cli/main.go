package main

import (
	"encoding/json"
	"flag"
	"os"

	transform "github.com/KablamoOSS/cfn-macros/CfnIngest/transform"
)

func main() {
	flag.Parse()
	template, err := transform.GetPath(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	err = template.Transform()
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(template)
}
