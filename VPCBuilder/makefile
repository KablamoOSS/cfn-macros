# Makefile for building python lambda
PROJECT = vpcbuilder
VERSION = $(shell whoami)
DIR = $(shell pwd)
GITHASH = $(shell git rev-parse HEAD | cut -c 1-7)
BUCKET_NAME = yourbucket.lambda.functions
CONTAINER = python:2.7.14-alpine3.6 
WORKING_DIR = /opt/app

buildDeps:
	docker run -v $(DIR):$(WORKING_DIR) -w $(WORKING_DIR) $(CONTAINER) pip install -r requirements.txt -t ./src
.PHONY: buildDeps

buildPackage: buildDeps
	cd ./src && zip -r9 ../$(PROJECT).zip *
.PHONY: buildPackage

uploadToS3: buildPackage
	aws s3 cp ./$(PROJECT).zip s3://$(BUCKET_NAME)/$(PROJECT)/$(PROJECT)-$(GITHASH).zip --acl bucket-owner-full-control --profile yourprofile
	echo 'File version is $(PROJECT)-$(GITHASH).zip'
.PHONY: uploadToS3