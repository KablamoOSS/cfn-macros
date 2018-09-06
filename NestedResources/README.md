# NestedResources Transformation

Reduces the reference complexity of CloudFormation templates by allowing
single-use child resources to be embedded into parameter lists. Resources can
be embedded arbitrarily deep, and the macro will flatten the tree out to a
standard CloudFormation resource list with references to logical ids.

## CloudFormation Macro

NestedResources is primarily intended to used as a CloudFormation macro.

`make nested-resources.zip` will build a zip file with a `go1.x` lambda that
can be used as a CloudFormation macro:

```
Transform: AWS::Serverless-2016-10-31
Resources:
  Function:
    Type: AWS::Serverless::Function
    Properties:
      Handler: handler
      Runtime: go1.x
      CodeUri: 's3://my-s3-bucket/nested-resources.zip'
  Macro:
    Type: AWS::CloudFormation::Macro
    Properties:
      Name: NestedResources
      Description: Reorganises nested resources from resource properties
      FunctionName: !Ref Function
```

With the macro created, adding `Transform: NestedResources` will apply the
transformation to stacks you create.

## Offline Usage

There may be cases where a CloudFormation macro may in inconvenient or
inappropriate. For these cases, `make nested-resources` will produce a binary
which can be used offline.

Only json templates are supported for offline usage at this time.

```
nested-resources my-template.json > my-transformed-template.json
```

Note: This will NOT manage the `Transform` key in the template. If
`NestedResources` is present in the `Transform` list, an offline transformation
will not remove it.

## Example:

Turns this:

```
{
  "Resources": {
    "MyGateway": {
      "Type": "AWS::ApiGateway::Deployment",
      "Properties": {
        "Description": "My API gateway deployment",
        "RestApiId": {
          "Name": "MyRestApi",
          "Type": "AWS::ApiGateway::RestApi",
          "Parameters": {
            ...
          }
        },
        "StageDescription": { ... }
      }
    }
  }
}
```

into this:

```
{
  "Resources": {
    "MyGateway": {
      "Type": "AWS::ApiGateway::Deployment",
      "Properties": {
        "Description": "My API gateway deployment",
        "RestApiId": { "Ref": "MyRestApi" },
        "StageDescription": { ... }
      }
    },
    "MyRestApi": {
      "Type": "AWS::ApiGateway::RestApi",
      "Parameters": {
        ...
      }
    }
  }
}
```
