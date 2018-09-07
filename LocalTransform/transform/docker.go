package transform

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type dockerTransformation struct {
	name    string
	image   string // e.g. lambci/lambda
	runtime string
	handler string
	file    string
}

func (t *Transformer) RegisterDocker(name, runtime, handler, file string) {
	// TODO: allow running a docker image (i.e. something lambci/lambda based).
	// Preference might be to supply a lambda .zip file to the base lambci image.
	// TODO: Which (if any) AWS vars to pass through?

	newTransformation := &dockerTransformation{
		name:    name,
		image:   "lambci/lambda", // FIXME: This is a sane default, but shouldn't be hardcoded.
		runtime: runtime,
		handler: handler,
		file:    file,
	}
	if t.Transforms == nil {
		t.Transforms = make(map[string]transformation)
	}
	t.Transforms[name] = newTransformation
}

func (d *dockerTransformation) apply(tmpl map[string]interface{}) (map[string]interface{}, error) {
	tmpdir, err := ioutil.TempDir("", "cfn-transform")
	if err != nil {
		return tmpl, err
	}
	// defer os.RemoveAll(tmpdir)

	err = os.Chmod(tmpdir, 0755)
	if err != nil {
		return tmpl, err
	}

	r, err := zip.OpenReader(d.file)
	if err != nil {
		return tmpl, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			err = os.Mkdir(filepath.Join(tmpdir, f.Name), 0755)
			if err != nil {
				return tmpl, err
			}
		} else {
			srcFd, err := f.Open()
			if err != nil {
				return tmpl, err
			}
			dstFd, err := os.OpenFile(filepath.Join(tmpdir, f.Name), os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				srcFd.Close()
				return tmpl, err
			}

			_, err = io.Copy(dstFd, srcFd)
			srcFd.Close()
			dstFd.Close()
			if err != nil {
				return tmpl, err
			}
		}
	}

	return d.runDockerCommand(tmpdir, tmpl)
}

// TODO: Quick and dirty first pass - consider using the docker api properly
// rather than wrapping the cli.
func (d *dockerTransformation) runDockerCommand(tmpdir string, tmpl map[string]interface{}) (map[string]interface{}, error) {
	cmd := exec.Command(
		"docker",
		"run",
		// "--rm",
		"-v",
		fmt.Sprintf("%s:/var/task", tmpdir),
		"-i",
		"-e",
		"DOCKER_LAMBDA_USE_STDIN=1",
		fmt.Sprintf("%s:%s", d.image, d.runtime),
		d.handler,
	)

	writer, err := cmd.StdinPipe()
	if err != nil {
		return tmpl, err
	}

	go func() {
		input := map[string]interface{}{
			"region":      "docker",
			"accountId":   "docker",
			"transformId": d.name,
			"fragment":    tmpl,
			"requestId":   tmpdir,
			"params":      map[string]interface{}{},
		}
		encoder := json.NewEncoder(writer)
		encoder.Encode(input)
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

	return newMap["fragment"].(map[string]interface{}), nil
}
