import boto3
import os
from boto3.dynamodb.conditions import Key, Attr

def lambda_handler(event, context):
    post_id = event['queryStringParameters']['postId']
    
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table(os.environ['DB_TABLE_NAME'])
    
    if post_id=="*":
        items = table.scan()
    else:
        items = table.query(KeyConditionExpression=Key('id').eq(post_id))
    
    return items["Items"]