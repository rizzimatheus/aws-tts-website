package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"ttswebsite/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
)

// NewPostBodyResponse is the response body for the new post endpoint
type NewPostBodyResponse struct {
	PostId string `json:"postId"`
}

// handler is the lambda handler function
func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received event: %+v\n", event)

	// Unmarshal the body of the request
	var post util.Post
	err := json.Unmarshal([]byte(event.Body), &post)
	if err != nil {
		log.Printf("Error unmarshalling post: %+v\n", err)
		return util.NewAPIGatewayProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	// Validate the post text length
	if len(post.Text) > 100 {
		log.Printf("Error: text too long (%d characters)\n", len(post.Text))
		return util.NewAPIGatewayProxyResponse(http.StatusBadRequest, `{error: "text too long"}`), nil
	}

	post.Id = uuid.New().String()
	post.Status = "PROCESSING"
	log.Printf("Post to add to the table: %+v\n", post)

	// Load the SDK's default configuration, credentials and region
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error loading SDK config: %+v\n", err)
		return util.NewAPIGatewayProxyResponse(http.StatusInternalServerError, `{error: "internal server error"}`), nil
	}

	// Create a DynamoDB client
	table := util.NewDynamodbTable(cfg, os.Getenv("DB_TABLE_NAME"))

	// Add the post to the DynamoDB table
	err = table.AddItem(post)
	if err != nil {
		log.Printf("Error adding item to table: %+v\n", err)
		return util.NewAPIGatewayProxyResponse(http.StatusInternalServerError, `{error: "internal server error"}`), nil
	}
	log.Printf("Post added to the table: %+v\n", post.Id)

	// Create a SNS client
	sns := util.SnsInstance{
		SnsClient: sns.NewFromConfig(cfg),
	}
	// Publish a message to the SNS topic
	err = sns.Publish(os.Getenv("SNS_TOPIC"), post.Id, "", "", "", "")
	if err != nil {
		log.Printf("Error publishing SNS message: %+v\n", err)
		return util.NewAPIGatewayProxyResponse(http.StatusInternalServerError, `{error: "internal server error"}`), nil
	}
	log.Printf("SNS message published: %+v\n", post.Id)

	body, _ := json.Marshal(NewPostBodyResponse{PostId: post.Id})
	return util.NewAPIGatewayProxyResponse(http.StatusCreated, string(body)), nil
}

func main() {
	lambda.Start(handler)
}
