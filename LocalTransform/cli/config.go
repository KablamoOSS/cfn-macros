package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/KablamoOSS/cfn-macros/LocalTransform/transform"
)

type ConfigFile map[string]ConfigEntry

type ConfigEntry struct {
	Type       string
	Properties json.RawMessage
}

type CmdProperties struct {
	Command string
	Args    []string
}

type DockerProperties struct {
	Image   string
	Runtime string
	Handler string
	File    string
}

type LambdaProperties struct {
	FunctionName string
	// region / profile / etc?
}

func parseConfig(t *transform.Transformer, filename string) error {
	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	conf := make(ConfigFile)

	err = json.Unmarshal(filedata, &conf)
	if err != nil {
		return err
	}

	for key, entry := range conf {
		switch entry.Type {
		case "cmd":
			props := &CmdProperties{}
			err := json.Unmarshal(entry.Properties, props)
			if err != nil {
				return err
			}
			t.RegisterCommand(key, props.Command, props.Args...)
		case "docker":
			props := &DockerProperties{
				Image: "lambci/lamda",
			}
			err := json.Unmarshal(entry.Properties, props)
			if err != nil {
				return err
			}
			t.RegisterDocker(key, props.Runtime, props.Handler, props.File)
		case "lambda":
			props := &LambdaProperties{}
			err := json.Unmarshal(entry.Properties, props)
			if err != nil {
				return err
			}
			t.RegisterLambda(key, props.FunctionName)
		default:
			return fmt.Errorf("parse config: invalid transform type for %s: %s", key, entry.Type)
		}
	}

	return nil
}
