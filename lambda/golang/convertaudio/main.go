package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"ttswebsite/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

func updatePost(table util.DynamodbTable, post util.Post, url string) {
	status := "Synthesized"
	if url == "" {
		status = "Failed"
	}
	post.Url = url
	post.Status = status
	attributeMap, err := table.UpdateItem(post)
	if err != nil {
		log.Printf("Error updating item: %+v\n", err)
		return
	}
	log.Printf("Post updated: %+v\n", attributeMap)
}

// handler is your Lambda function handler
func handler(ctx context.Context, event events.SNSEvent) {
	log.Printf("Received event: %+v\n", event)
	url := ""
	
	// Get the dynamodb table and post IDs from the SNS message
	var message util.SnsMessage
	err := json.Unmarshal([]byte(event.Records[0].SNS.Message), &message)
	if err != nil {
		log.Printf("Error unmarshalling message: %+v\n", err)
		return
	}
	log.Printf("Message: %+v\n", message)
	
	// Load the SDK's default configuration, credentials and region
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error loading SDK config: %+v\n", err)
		return
	}
	
	// Create a DynamoDB client
	table := util.NewDynamodbTable(cfg, message.DynamodbTableId)

	// Get the post
	post, err := table.GetItem(message.DynamodbRecordId)
	if err != nil {
		log.Printf("Error getting posts: %+v\n", err)
		updatePost(table, post, url)
		return
	}

	// Synthesize the post's text
	pollyInstance := util.NewPollyInstance(cfg)
	output, err := pollyInstance.SynthesizeSpeech(post.Voice, "<speak>"+post.Text+"</speak>")
	if err != nil {
		log.Printf("Error synthesizing speech: %+v\n", err)
		updatePost(table, post, url)
		return
	}
	// Save the audio to a file
	filePath := "/tmp/" + message.DynamodbRecordId + ".mp3"
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %+v\n", err)
		updatePost(table, post, url)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, output.AudioStream)
	if err != nil {
		log.Printf("Error copying audio stream: %+v\n", err)
		updatePost(table, post, url)
		return
	}

	log.Println("Audio synthesized")

	// Upload the audio to S3
	objectKey := message.DynamodbRecordId + ".mp3"
	bucket := util.NewS3Bucket(cfg, os.Getenv("BUCKET_NAME"))
	err = bucket.UploadFile(objectKey, filePath)
	if err != nil {
		log.Printf("Error uploading file: %+v\n", err)
		updatePost(table, post, url)
		return
	}
	defer os.Remove(filePath)
	log.Println("Audio uploaded to S3")

	// Update the post with the audio URL and new status
	url, err = bucket.GetObjectUrl(objectKey)
	if err != nil {
		log.Printf("Error getting object URL: %+v\n", err)
		updatePost(table, post, url)
		return
	}
	log.Printf("Audio URL: %s\n", url)
	updatePost(table, post, url)
}

func main() {
	lambda.Start(handler)
}
