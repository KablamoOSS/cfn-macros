package transform

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type lambdaTransformation struct {
	name         string
	functionName string
}

func (t *Transformer) RegisterLambda(name, functionName string) {
	newTransformation := &lambdaTransformation{
		name:         name,
		functionName: functionName,
	}
	if t.Transforms == nil {
		t.Transforms = make(map[string]transformation)
	}
	t.Transforms[name] = newTransformation
}

func (t *lambdaTransformation) apply(tmpl map[string]interface{}) (map[string]interface{}, error) {
	sess := session.New()
	client := lambda.New(sess)

	inputTmpl := map[string]interface{}{
		"region":      "lambda",
		"accountId":   "lambda",
		"transformId": t.name,
		"fragment":    tmpl,
		"requestId":   t.name,
		"params":      map[string]interface{}{},
	}

	input, err := json.Marshal(inputTmpl)
	if err != nil {
		return tmpl, err
	}

	invokeOutput, err := client.Invoke(
		&lambda.InvokeInput{
			FunctionName: &t.functionName,
			Payload:      input,
		},
	)
	if err != nil {
		return tmpl, err
	}

	if invokeOutput.FunctionError != nil {
		return tmpl, fmt.Errorf("invoke lambda %s: %s functionError", t.functionName, *invokeOutput.FunctionError)
	}

	if *invokeOutput.StatusCode != 200 {
		return tmpl, fmt.Errorf("invoke lambda %s: Status Code %d", t.functionName, *invokeOutput.StatusCode)
	}

	newTmpl := map[string]interface{}{}
	err = json.Unmarshal(invokeOutput.Payload, &newTmpl)
	if err != nil {
		return tmpl, err
	}

	return newTmpl["fragment"].(map[string]interface{}), nil
}
