# CloudFormation Ingest

Transform that ingests parameters, resources, outputs, etc from other
CloudFormation templates. This allows the creation of reusable components as
standalone CloudFormation templates, which can be composed together to create
more complex configurations.

`cfn-ingest` recognises an additional `Ingest` stanza in CloudFormation
templates, which reference external templates either on the local filesystem or
in an S3 bucket. All the resources, mappings and conditions in the ingested
template will be added to the template, and the output will be a standard
CloudFormation template with the combined set.

When ingesting a template, you may supply resource references (or other values)
as parameters to the ingested template. Any ingested parameters that don't have
values passed are instead added to the ingesting templates parameters.

Outputs of ingested templates can be referenced in the ingesting template via
`{ "Ingest::Ref": "{{TemplateName}}.{{OutputName}}" }` style references. You
may choose whether or not to include the ingested outputs in the final
template, or select specific outputs to include.

For very complex configurations, ingested templates that themselves contain an
`Ingest` stanza will have those ingestions processed, allowing multiple levels
of composition.

## Building

`make cfn-ingest.zip` will build a zip file with a `go1.x` lambda that
can be used as a CloudFormation macro:

Run `make cfn-ingest` to build the CLI, which is usable offline.

Sample template to create the macro:

```
Transform: AWS::Serverless-2016-10-31
Resources:
  Function:
    Type: AWS::Serverless::Function
    Properties:
      Handler: handler
      Runtime: go1.x
      CodeUri: 's3://my-s3-bucket/cfn-ingest.zip'
  Macro:
    Type: AWS::CloudFormation::Macro
    Properties:
      Name: CfnIngest
      Description: Ingests external templates
      FunctionName: !Ref Function
```

With the macro created, adding `Transform: CfnIngest` will apply the
transformation to stacks you create.

## Local usage

When running locally, simply run the `cfn-ingest` command with the base template path as an argument:

    cfn-ingest base-template.json

or:

    cfn-ingest s3://my-bucket/base-template.json

Note: To use shared configuration from `~/.aws`, setting the environment variable `AWS_SDK_LOAD_CONFIG=1` may be required.

## (Badly contrived) Example

We've created a standard CloudFormation template for our preferred EC2 instance
configuration. In this case, an instance which requires a security group, and
has a volume with optional size configuration:

    {
      "Parameters":{
        "VolumeSize": {
          "Type": "Number",
          "DefaultValue": 100
        },
        "SecurityGroup": {
          "Type": "String"
        }
      },
      "Resources": {
        "MyVolumeAttachment": {
          "Type": "AWS::EC2::VolumeAttachment",
          "Properties": {
            "Device": "/dev/sda",
            "InstanceId": { "Ref": "MyInstance" },
            "VolumeIds": [ { "Ref": "MyVolume" } ]
          }
        },
        "MyInstance": {
          "Type": "AWS::EC2::Instance",
          "Properties": {
            "SecurityGroup": { "Ref": "SecurityGroup" }
          }
        },
        "MyVolume": {
          "Type": "AWS::EC2::Volume",
          "Properties": {
            "Size": { "Ref": "VolumeSize" }
          }
        }
      },
      "Outputs":{
        "InstanceId": {
          "Value": { "Ref": "MyInstance" }
        }
      }
    }

We then want to create a more complex template containing such an EC2 instance.
For this example, we'll just create a SecurityGroup for the instance, and an
IAM role with permissions to manage it:

    {
      "Ingest": {
        "InstanceWithVolume": {
          "Source": "instance-with-volume.json",
          "Parameters": {
            "SecurityGroup": { "Ref": "SecurityGroup" }
          },
          "IngestOutputs": true
        }
      },
      "Resources": {
        "MySecurityGroup": {
          "Type": "AWS::EC2::SecurityGroup",
          "Properties": {
          }
        },
        "MyRole": {
          "Type": "AWS::IAM::Role",
          "Properties": {
            "Policy": {
              "PolicyName": "ManageMyInstance",
              "PolicyDocument": {
                "Version" : "2012-10-17",
                "Statement": [
                  {
                    "Effect": "Allow",
                    "Action": "ec2:*",
                    "Resource": { "Fn::Join": [ "arn:aws:ec2:::instance/", { "Ingest::Ref": "InstanceWithVolume.InstanceId" } ] }
                  }
                ]
              }
            }
          }
        }
      }
    }

The ingest stanza here indicates that we will be ingesting one other template,
referenced by `InstanceWithVolume` to be read from the local file
`instance-with-volume.json`. We supply a reference to the SecurityGroup in this
template to the `SecurityGroup` parameter of the ingested template. Since we
don't supply a value to the `VolumeSize` parameter, it will be added to the set
of parameters in the final template.

The `IngestOutputs` parameter indicates that the outputs of
`InstanceWithVolume` should be included in the generated template -- which is
the default behaviour if `IngestOutputs` is omitted. If set to `false` outputs
will not be included in the generated template, but will still be referecable
within the template via `Ingest::Ref`. `IngestOutputs` may also be an array of
output names to select a subset of outputs to be included.

The `MyRole` IAM Role resource contains an `Ingest::Ref` to
`InstanceWithVolume.InstanceId`, which will be replaced with the output value
in the generated template.

Here is the final output, after running `cfn-ingest`:

    {
      "Parameters": {
        "VolumeSize": {
          "DefaultValue": 100,
          "Type": "Number"
        }
      },
      "Resources": {
        "MyInstance": {
          "Properties": {
            "SecurityGroup": {
              "Ref": "SecurityGroup"
            }
          },
          "Type": "AWS::EC2::Instance"
        },
        "MyRole": {
          "Properties": {
            "Policy": {
              "PolicyDocument": {
                "Statement": [
                  {
                    "Action": "ec2:*",
                    "Effect": "Allow",
                    "Resource": {
                      "Fn::Join": [
                        "arn:aws:ec2:::instance/",
                        {
                          "Ref": "MyInstance"
                        }
                      ]
                    }
                  }
                ],
                "Version": "2012-10-17"
              },
              "PolicyName": "ManageMyInstance"
            }
          },
          "Type": "AWS::IAM::Role"
        },
        "MySecurityGroup": {
          "Properties": {},
          "Type": "AWS::EC2::SecurityGroup"
        },
        "MyVolume": {
          "Properties": {
            "Size": {
              "Ref": "VolumeSize"
            }
          },
          "Type": "AWS::EC2::Volume"
        },
        "MyVolumeAttachment": {
          "Properties": {
            "Device": "/dev/sda",
            "InstanceId": {
              "Ref": "MyInstance"
            },
            "VolumeIds": [
              {
                "Ref": "MyVolume"
              }
            ]
          },
          "Type": "AWS::EC2::VolumeAttachment"
        }
      },
      "Outputs": {
        "InstanceId": {
          "Value": {
            "Ref": "MyInstance"
          }
        }
      }
    }
