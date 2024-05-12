import boto3
import os
import uuid
import json

def lambda_handler(event, context):
    record_id = str(uuid.uuid4())

    body = json.loads(event["body"])
    voice = body["voice"]
    text = body["text"]

    print('Generating new DynamoDB record, with ID: ' + record_id)
    print('Input Text: ' + text)
    print('Selected voice: ' + voice)
    
    #Creating new record in DynamoDB table
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table(os.environ['DB_TABLE_NAME'])
    table.put_item(
        Item={
            'id' : record_id,
            'text' : text,
            'voice' : voice,
            'status' : 'Processing'
        }
    )

    print("DynamoDB item created")
    
    #Sending notification about new post to SNS
    client = boto3.client('sns')
    client.publish(
        TopicArn = os.environ['SNS_TOPIC'],
        Message = record_id
    )
    
    return {"postId": record_id}