package main

import (
	"github.com/KablamoOSS/cfn-macros/CfnIngest/transform"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	// Don't allow attempting to load local files when running from a lambda
	transform.HaveLocalFilesystem = false
}

type MacroRequest struct {
	Region         string                 `json:"region"`
	AccountID      string                 `json:"accountId"`
	Fragment       *transform.CfnTemplate `json:"fragment"`
	TransformID    string                 `json:"transformId"`
	Params         map[string]interface{} `json:"params"`
	RequestID      string                 `json:"requestId"`
	TemplateParams map[string]interface{} `json:"templateParameterValues"`
}

type MacroResponse struct {
	RequestID string                 `json:"requestId"`
	Status    string                 `json:"status"`
	Fragment  *transform.CfnTemplate `json:"fragment"`
}

type Resource struct {
	Type       string
	Parameters map[string]interface{}
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(req MacroRequest) (MacroResponse, error) {
	resp := MacroResponse{
		Status:    "success",
		RequestID: req.RequestID,
		Fragment:  nil,
	}

	err := req.Fragment.Transform()
	if err != nil {
		return resp, err
	}

	resp.Fragment = req.Fragment

	return resp, nil
}
