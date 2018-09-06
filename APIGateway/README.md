# API Gateway Macro

#### Aims:
 - Demonstrate the Use of both AWS hosted transforms and Custom transforms
 - Demonstrate the use of multiple transforms in a single template

## Transform
This transform adds standard Outputs to any cloudformation template as well as specific ones for a SAM API deployment.

Observer that the API specification file in `API/template.yaml` does not create any outputs however, when the stack successfully deploys the URL for the deployed stage is listed in the stack outputs.  This is done by the `APIG-2018-08-18` transform.

### Deploying the Transform

```
ARTIFACT_BUCKET=<your_s3_bucket>
aws cloudformation package --template-file Macro/macro.yaml --s3-bucket $ARTIFACT_BUCKET --s3-prefix macros/API  --output-template /tmp/packaged.yaml
aws cloudformation deploy --capabilities CAPABILITY_IAM --template-file /tmp/packaged.yaml --stack-name API-macro
```

## Deploy API

```
ARTIFACT_BUCKET=<your_existing_bucket>
aws cloudformation package --template-file API/template.yaml --s3-bucket $ARTIFACT_BUCKET --s3-prefix API/lambdas  --output-template /tmp/packaged.yaml
aws cloudformation deploy --capabilities CAPABILITY_IAM --template-file /tmp/packaged.yaml --stack-name my-API-dev
```

## Testing
Check the stack outputs of `my-API-dev` to find `ENDPOINT_URL`.  It has been added by the `APIG-2018-08-18` transform.

Test your API.

```bash
curl ${ENDPOINT_URL}/ping
{"message": "pong", "response": {}}
```

