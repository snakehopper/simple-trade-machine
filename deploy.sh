#!/bin/bash

set -e

go generate ./...
go test ./...

FUNCTION_TARGET=AlertHandler
FUNCTION_NAME=simple-trade-machine

gcloud beta functions deploy ${FUNCTION_NAME} \
--gen2 \
--runtime go116 \
--trigger-http \
--entry-point ${FUNCTION_TARGET} \
--env-vars-file=env.yaml \
--source . \
--allow-unauthenticated