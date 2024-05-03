package util

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Post struct {
	Id     string `dynamodbav:"id" json:"id"`
	Text   string `dynamodbav:"text" json:"text"`
	Voice  string `dynamodbav:"voice" json:"voice"`
	Status string `dynamodbav:"status" json:"status"`
	Url    string `dynamodbav:"url" json:"url"`
}

// GetKey returns the primary key of the post in a format that can be
// sent to DynamoDB.
func (p Post) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(p.Id)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"id": id}
}
