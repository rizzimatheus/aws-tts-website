package main

import (
	"context"
	"log"
	"os"
	"ttswebsite/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// NewPostBodyResponse is the response body for the new post endpoint
type NewPostBodyResponse struct {
	PostId string `json:"postId"`
}

// handler is the lambda handler function
func handler(ctx context.Context, event events.DynamoDBEvent) {
	log.Printf("Received event: %+v\n", event)

	// Load the SDK's default configuration, credentials and region
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error loading SDK config: %+v\n", err)
	}

	// Create a SNS client
	sns := util.SnsInstance{
		SnsClient: sns.NewFromConfig(cfg),
	}

	for _, record := range event.Records {
		postId := record.Change.Keys["id"].String()
		tableId, err := util.DynamoDBTableNameFromArn(record.EventSourceArn)
		if err != nil {
			log.Printf("Error getting table name from ARN: %+v\n", err)
			continue
		}

		// Publish a message to the SNS topic
		msg := util.NewSnsMessage(tableId, postId)
		err = sns.Publish(os.Getenv("SNS_TOPIC"), msg, "", "", "", "")
		if err != nil {
			log.Printf("Error publishing SNS message: %+v\n", err)
		}
		log.Printf("SNS message published: %+v\n", msg)
	}

}

func main() {
	lambda.Start(handler)
}
