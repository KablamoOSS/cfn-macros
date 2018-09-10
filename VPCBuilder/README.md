# VPC Builder

Builds out a "fully" featured VPC summarising the complexity associated with a VPC such as Internet & Customer Gateways, Subnets, Routetables and NATGateways.

It also adds in VPC Flowlogs with an IAM role and supports full dynamic allocation of IPv6 with the VPC and to each subnet.

The IPv6 handles Egress Internet Gateway and default route against ::/0

## Deploying the Transform

```
ARTIFACT_BUCKET=<your_s3_bucket>
aws cloudformation package --template-file Transform/transform.yaml --s3-bucket $ARTIFACT_BUCKET --s3-prefix macros/VPCBuilder  --output-template /tmp/packaged.yaml
aws cloudformation deploy --capabilities CAPABILITY_IAM --template-file /tmp/packaged.yaml --stack-name VPCBuilder-macro
```

## Basic Usage

Utilise the yaml structure below in the [template](VPC/example.yaml). It will support the removal of Subnets, RouteTables, NATGateways and NetworkACLs.

```
aws cloudformation deploy --capabilities CAPABILITY_IAM \
        --template-file VPC/example.yaml \
        --stack-name VPC \
        --parameter-overrides VGW=vgw-030ba2b7b2c0ce5d5
```

## To Do

- Explaination about dependencies (VGW)
- Add Outputs with Exports for critical resources
- VPC Endpoints for all AWS Services
- Add a little better handling of custom pieces (e.g. different route gateways)
- Adding proper IPv6 regex and handling with NetworkACLs

