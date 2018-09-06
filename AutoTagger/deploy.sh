#!/bin/bash

if [ -z "${ARTIFACT_BUCKET}" ]; then
    echo "This deployment script needs an S3 bucket to store CloudFormation artifacts."
    echo "You can also set this by doing: export ARTIFACT_BUCKET=my-bucket-name"
    echo
    read -p "S3 bucket to store artifacts: " ARTIFACT_BUCKET
fi

aws cloudformation package \
    --template-file transform.yaml \
    --s3-bucket ${ARTIFACT_BUCKET} \
    --output-template-file packaged.yaml

aws cloudformation deploy \
    --stack-name autotagger-transform \
    --template-file packaged.yaml \
    --capabilities CAPABILITY_IAM

aws cloudformation deploy \
    --stack-name autotagger-example \
    --template-file example.yaml \
    --capabilities CAPABILITY_IAM
