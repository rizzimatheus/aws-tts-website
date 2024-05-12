package util

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

// SnsInstance is a wrapper around the SNS client.
type SnsInstance struct {
	SnsClient *sns.Client
}

// Publish publishes a message to an Amazon SNS topic. The message is then sent to all
// subscribers. When the topic is a FIFO topic, the message must also contain a group ID
// and, when ID-based deduplication is used, a deduplication ID. An optional key-value
// filter attribute can be specified so that the message can be filtered according to
// a filter policy.
func (s SnsInstance) Publish(topicArn string, message string, groupId string, dedupId string, filterKey string, filterValue string) error {
	publishInput := sns.PublishInput{TopicArn: aws.String(topicArn), Message: aws.String(message)}
	if groupId != "" {
		publishInput.MessageGroupId = aws.String(groupId)
	}
	if dedupId != "" {
		publishInput.MessageDeduplicationId = aws.String(dedupId)
	}
	if filterKey != "" && filterValue != "" {
		publishInput.MessageAttributes = map[string]types.MessageAttributeValue{
			filterKey: {DataType: aws.String("String"), StringValue: aws.String(filterValue)},
		}
	}
	_, err := s.SnsClient.Publish(context.TODO(), &publishInput)
	if err != nil {
		log.Printf("Couldn't publish message to topic %v. Here's why: %v", topicArn, err)
	}
	return err
}

// SnsMessage is a struct that represents a message sent to an Amazon SNS topic.
type SnsMessage struct {
	DynamodbTableId  string `json:"dynamodb_table_id"`
	DynamodbRecordId string `json:"dynamodb_record_id"`
}

// NewSnsMessage returns a string of a new SnsMessage struct with the specified table and record IDs.
func NewSnsMessage(tableId string, recordId string) string {
	j, _ := json.Marshal(SnsMessage{DynamodbTableId: tableId, DynamodbRecordId: recordId})
	return string(j)
}
