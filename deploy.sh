FUNCTION_TARGET=PredAlert
FUNCTION_NAME=go-http-function

gcloud beta functions deploy ${FUNCTION_NAME} \
--gen2 \
--runtime go116 \
--trigger-http \
--entry-point ${FUNCTION_TARGET} \
--env-vars-file=env.yaml \
--source . \
--allow-unauthenticated