package main

import (
	"context"
	"io"
	"log"
	"os"
	"ttswebsite/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

// handler is your Lambda function handler
func handler(ctx context.Context, event events.SNSEvent) {
	log.Printf("Received event: %+v\n", event)

	// Get the post ID from the SNS message
	postId := event.Records[0].SNS.Message
	log.Printf("Post ID: %s\n", postId)

	// Load the SDK's default configuration, credentials and region
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error loading SDK config: %+v\n", err)
		return
	}

	// Create a DynamoDB client
	table := util.NewDynamodbTable(cfg, os.Getenv("DB_TABLE_NAME"))

	// Get the post
	post, err := table.GetItem(postId)
	if err != nil {
		log.Printf("Error getting posts: %+v\n", err)
		return
	}

	// Synthesize the post's text
	pollyInstance := util.NewPollyInstance(cfg)
	output, err := pollyInstance.SynthesizeSpeech(post.Voice, "<speak>"+post.Text+"</speak>")
	if err != nil {
		log.Printf("Error synthesizing speech: %+v\n", err)
		return
	}
	// Save the audio to a file
	filePath := "/tmp/" + postId + ".mp3"
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %+v\n", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, output.AudioStream)
	if err != nil {
		log.Printf("Error copying audio stream: %+v\n", err)
		return
	}

	log.Println("Audio synthesized")

	// Upload the audio to S3
	objectKey := postId + ".mp3"
	bucket := util.NewS3Bucket(cfg, os.Getenv("BUCKET_NAME"))
	err = bucket.UploadFile(objectKey, filePath)
	if err != nil {
		log.Printf("Error uploading file: %+v\n", err)
		return
	}
	defer os.Remove(filePath)
	log.Println("Audio uploaded to S3")

	// Update the post with the audio URL and new status
	url, err := bucket.GetObjectUrl(objectKey)
	if err != nil {
		log.Printf("Error getting object URL: %+v\n", err)
		return
	}
	log.Printf("Audio URL: %s\n", url)
	post.Url = url
	post.Status = "SYNTHESIZED"
	attributeMap, err := table.UpdateItem(post)
	if err != nil {
		log.Printf("Error updating item: %+v\n", err)
		return
	}
	log.Printf("Post updated: %+v\n", attributeMap)
}

func main() {
	lambda.Start(handler)
}
