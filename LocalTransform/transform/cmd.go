package transform

import (
	"encoding/json"
	"os/exec"
)

type cmdTransformation struct {
	command string
	args    []string
}

func (t *Transformer) RegisterCommand(name, cmd string, args ...string) {
	newTransformation := &cmdTransformation{
		command: cmd,
		args:    args,
	}
	if t.Transforms == nil {
		t.Transforms = make(map[string]transformation)
	}
	t.Transforms[name] = newTransformation
}

func (c *cmdTransformation) apply(tmpl map[string]interface{}) (map[string]interface{}, error) {
	cmd := exec.Command(c.command, c.args...)
	writer, err := cmd.StdinPipe()
	if err != nil {
		return tmpl, err
	}

	go func() {
		encoder := json.NewEncoder(writer)
		encoder.Encode(tmpl)
		writer.Close()
	}()

	outJson, err := cmd.Output()
	if err != nil {
		return tmpl, err
	}

	newMap := make(map[string]interface{})

	err = json.Unmarshal(outJson, &newMap)
	if err != nil {
		return tmpl, err
	}

	return newMap, nil
}
