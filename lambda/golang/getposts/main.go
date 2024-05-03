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
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received event: %+v\n", event)

	// Get the postId from the query string
	postId := event.QueryStringParameters["postId"]
	log.Printf("postId: %s\n", postId)

	// Load the SDK's default configuration, credentials and region
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error loading SDK config: %+v\n", err)
		return util.NewAPIGatewayProxyResponse(http.StatusInternalServerError, `{error: "internal server error"}`), nil
	}

	// Create a DynamoDB client
	table := util.NewDynamodbTable(cfg, os.Getenv("DB_TABLE_NAME"))

	var posts []util.Post
	if postId == "*" {
		posts, err = table.GetAllItems()
		if err != nil {
			log.Printf("Error getting posts: %+v\n", err)
			return util.NewAPIGatewayProxyResponse(http.StatusInternalServerError, `{error: "internal server error"}`), nil
		}
	} else {
		post, err := table.GetItem(postId)
		if err != nil {
			log.Printf("Error getting posts: %+v\n", err)
			return util.NewAPIGatewayProxyResponse(http.StatusInternalServerError, `{error: "internal server error"}`), nil
		}
		posts = append(posts, post)
	}

	body, _ := json.Marshal(posts)
	return util.NewAPIGatewayProxyResponse(http.StatusOK, string(body)), nil
}

func main() {
	lambda.Start(handler)
}
