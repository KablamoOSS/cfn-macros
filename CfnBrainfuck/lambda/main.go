package main

import (
	"encoding/json"
	"github.com/KablamoOSS/cfn-macros/CfnBrainfuck/transform"
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
	RequestID string      `json:"requestId"`
	Status    string      `json:"status"`
	Fragment  interface{} `json:"fragment"`
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

	fragment := make(map[string]interface{})
	output := transform.Transform(req.Fragment)
	err := json.Unmarshal(output, &fragment)
	if err != nil {
		resp.Fragment = string(output)
	} else {
		resp.Fragment = fragment
	}

	return resp, nil
}
