package main

import (
	"encoding/json"
	"flag"
	"github.com/KablamoOSS/cfn-macros/CfnIngest/transform"
	"os"
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
