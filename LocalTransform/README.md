# cfn-transform

A tool for running CloudFormation Macro style transformations outside of a
CloudFormation ChangeSet. Suitable for use in CI pipelines, or for testing
transformation lambda functions.

## Usage

`cfn-transform` will read a CloudFormation template from either stdin, or from
a provided filename. From this template, `cfn-transform` will look for a
`Transform` key with one or more transforms to apply. When all transforms in
the sequence has been applied, the final result will be printed to stdout.

    cat input-template.json | cfn-transform config.json
    # or
    cfn-transform config.json input-template.json

Note: Embedded `Fn::Transform` calls are not currently supported and will be
ignored.

Each of the transforms in the template must be specified in a config file, also
in JSON format. These can be commands that run locally, lambda functions run in
a docker container (see `lambci/lambda`), or actual lambda functions to be
invoked from your AWS account.

    {
      "CommandTransform": {
        "Type": "command",
        "Properties": {
          "Command": "/usr/bin/my-transform"
          "Args": [
            "Lorem",
            "ipsum",
            "etc"
          ]
        }
      },
      "DockerTransform": {
        "Type": "docker",
        "Properties": {
          "Runtime": "go1.x",
          "Handler": "handler",
          "File": "my-lambda.zip"
        }
      },
      "LambdaTransform": {
        "Type": "lambda",
        "Properties": {
          "FunctionName": "my-lambda"
        }
      }
    }


## Commands

Commands can be used to run a transform that supports being run outside of
lambda. They should be able to read a cloudformation template from stdin in
JSON format, and write the transformed template to stdout.

## Docker

Run a standard lambda via a docker compatibility layer (i.e. `lambci/lambda`)
by specifying the runtime, handler and the location of the lambda zip file.

## Lambda

If you already have transform macros, but just want to be able to run them
outside of a CloudFormation change set or mix them with other transform types,
`cfn-transform` is able to invoke them for you.
