package main

import (
	"github.com/KablamoOSS/cfn-macros/NestedResources/transform"
	"github.com/aws/aws-lambda-go/lambda"
)

type MacroRequest struct {
	Region         string                 `json:"region"`
	AccountID      string                 `json:"accountId"`
	Fragment       map[string]interface{} `json:"fragment"`
	TransformID    string                 `json:"transformId"`
	Params         map[string]interface{} `json:"params"`
	RequestID      string                 `json:"requestId"`
	TemplateParams map[string]interface{} `json:"templateParameterValues"`
}

type MacroResponse struct {
	RequestID string                 `json:"requestId"`
	Status    string                 `json:"status"`
	Fragment  map[string]interface{} `json:"fragment"`
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
		Fragment:  make(map[string]interface{}),
	}

	resp.Fragment = transform.Transform(resp.Fragment)

	return resp, nil
}
