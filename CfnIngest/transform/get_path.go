package transform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// lambda handler can set this false to indicate that local files are unavailable
var HaveLocalFilesystem = true

func GetPath(path string) (*CfnTemplate, error) {
	if strings.HasPrefix(path, "s3://") {
		parts := strings.SplitN(path[5:], "/", 2)
		return getS3Object(parts[0], parts[1])
	}

	if !HaveLocalFilesystem {
		return nil, fmt.Errorf("no local filesystem available")
	}

	return getLocalFile(path)
}

func getLocalFile(path string) (*CfnTemplate, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmpl := &CfnTemplate{}
	err = json.Unmarshal(body, tmpl)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func getS3Object(bucket, key string) (*CfnTemplate, error) {
	tmpl := &CfnTemplate{}
	sess := session.Must(session.NewSession())
	client := s3.New(sess)
	out, err := client.GetObject(
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		},
	)
	if err != nil {
		return nil, err
	}

	defer out.Body.Close()
	body, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, tmpl)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
