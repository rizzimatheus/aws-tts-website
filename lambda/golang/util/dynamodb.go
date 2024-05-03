package util

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamodbTable represents a DynamoDB table.
type DynamodbTable struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func NewDynamodbTable(config aws.Config, tableName string) DynamodbTable {
	return DynamodbTable{
		DynamoDbClient: dynamodb.NewFromConfig(config),
		TableName:      tableName,
	}
}

// AddItem adds a post to the DynamoDB table.
func (t DynamodbTable) AddItem(post Post) error {
	item, err := attributevalue.MarshalMap(post)
	if err != nil {
		panic(err)
	}
	_, err = t.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

// UpdateItem updates the status and url of a post that already exists in the
// DynamoDB table. This function uses the `expression` package to build the update
// expression.
func (t DynamodbTable) UpdateItem(post Post) (map[string]map[string]interface{}, error) {
	var err error
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]map[string]interface{}
	update := expression.Set(expression.Name("status"), expression.Value(post.Status))
	update.Set(expression.Name("url"), expression.Value(post.Url))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		response, err = t.DynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName:                 aws.String(t.TableName),
			Key:                       post.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			log.Printf("Couldn't update post %v. Here's why: %v\n", post.Id, err)
		} else {
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				log.Printf("Couldn't unmarshall update response. Here's why: %v\n", err)
			}
		}
	}
	return attributeMap, err
}

// GetItem retrieves a post from the DynamoDB table.
func (t DynamodbTable) GetItem(id string) (Post, error) {
	post := Post{Id: id}
	response, err := t.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: post.GetKey(), TableName: aws.String(t.TableName),
	})
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", id, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &post)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}
	return post, err
}

// GetAllItems retrieves all posts from the DynamoDB table.
func (t DynamodbTable) GetAllItems() ([]Post, error) {
	var posts []Post
	var err error
	var response *dynamodb.ScanOutput
	// Create a paginator for the DynamoDB table.
	scanPaginator := dynamodb.NewScanPaginator(t.DynamoDbClient, &dynamodb.ScanInput{
		TableName: aws.String(t.TableName),
	})
	// Scan the DynamoDB table and unmarshal the response into a slice of Post structs.
	for scanPaginator.HasMorePages() {
		response, err = scanPaginator.NextPage(context.TODO())
		if err != nil {
			log.Printf("Couldn't scan for posts. Here's why: %v\n", err)
			break
		} else {
			var postPage []Post
			err = attributevalue.UnmarshalListOfMaps(response.Items, &postPage)
			if err != nil {
				log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
				break
			} else {
				posts = append(posts, postPage...)
			}
		}
	}
	return posts, err
}

// Query gets all posts in the DynamoDB table that is owned by a specified user.
// The function uses the `expression` package to build the key condition expression
// that is used in the query.
func (t DynamodbTable) Query(userId string) ([]Post, error) {
	var err error
	var response *dynamodb.QueryOutput
	var posts []Post

	keyEx := expression.Key("owner").Equal(expression.Value(userId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)
	} else {
		// Create a paginator for the DynamoDB table.
		queryPaginator := dynamodb.NewQueryPaginator(t.DynamoDbClient, &dynamodb.QueryInput{
			TableName:                 aws.String(t.TableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
		})
		// Query the DynamoDB table and unmarshal the response into a slice of Post structs.
		for queryPaginator.HasMorePages() {
			response, err = queryPaginator.NextPage(context.TODO())
			if err != nil {
				log.Printf("Couldn't query for posts owned by %v. Here's why: %v\n", userId, err)
				break
			} else {
				var postPage []Post
				err = attributevalue.UnmarshalListOfMaps(response.Items, &postPage)
				if err != nil {
					log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
					break
				} else {
					posts = append(posts, postPage...)
				}
			}
		}
	}
	return posts, err
}

// TableExists determines whether a DynamoDB table exists.
func (t DynamodbTable) TableExists() (bool, error) {
	exists := true
	_, err := t.DynamoDbClient.DescribeTable(
		context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(t.TableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exist.\n", t.TableName)
			err = nil
		} else {
			log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", t.TableName, err)
		}
		exists = false
	}
	return exists, err
}
