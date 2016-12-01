package main

import (
	"encoding/json"
	"log"

	"github.com/eawsy/aws-lambda-go/service/lambda/runtime"
)

func handleLambda(rawEvent json.RawMessage, context *runtime.Context) (interface{}, error) {
	log.Printf("Starting %s…", context.FunctionName)

	event, err := ParseEc2StateChangeNotification(rawEvent)
	if err != nil {
		log.Printf("Failed to parse event: %v", err)
		return nil, err
	}

	handler := Handler{}
	if err = handler.HandleEc2StateChange(event.Detail); err != nil {
		return nil, err
	}

	log.Printf("Done %s.", context.FunctionName)
	return nil, nil
}

func init() {
	runtime.HandleFunc(handleLambda)
}

func main() {}
