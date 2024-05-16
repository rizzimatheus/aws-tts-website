import boto3
import os
from contextlib import closing
from boto3.dynamodb.conditions import Key, Attr

def lambda_handler(event, context):

    post_id = event["Records"][0]["Sns"]["Message"]

    print("Text to Speech function. Post ID in DynamoDB: " + post_id)

    # Retrieving information about the post from DynamoDB table
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table(os.environ['DB_TABLE_NAME'])
    post_item = table.get_item(
        Key = {
            'id': post_id,
        }
    )

    text = '<speak>' + post_item["Item"]["text"] + '</speak>'
    voice = post_item["Item"]["voice"]

    print("Post Text: " + text)
    print("Post Voice: " + voice)

    # Invoke Polly API, which will transform text into audio
    polly = boto3.client('polly')
    response = polly.synthesize_speech(
        Engine = "standard",
        OutputFormat='mp3',
        Text = text,
        TextType = 'ssml',
        VoiceId = voice,
    )
    
    # Save the audio stream returned by Amazon Polly on Lambda's temp directory.
    if "AudioStream" in response:
        with closing(response["AudioStream"]) as stream:
            output = os.path.join("/tmp/", post_id)
            with open(output, "ab") as file:
                file.write(stream.read())
    
    print("Audio synthesized")

    s3 = boto3.client('s3')
    s3.upload_file(
        '/tmp/' + post_id,
        os.environ['BUCKET_NAME'],
        post_id + ".mp3"
    )
    s3.put_object_acl(
        ACL='public-read',
        Bucket=os.environ['BUCKET_NAME'],
        Key= post_id + ".mp3"
    )
    
    print("Audio stored")

    location = s3.get_bucket_location(Bucket=os.environ['BUCKET_NAME'])
    region = location['LocationConstraint']

    if region is None:
        url_beginning = "https://s3.amazonaws.com/"
    else:
        url_beginning = "https://s3." + str(region) + ".amazonaws.com/" \

    url = url_beginning \
        + str(os.environ['BUCKET_NAME']) \
        + "/" \
        + str(post_id) \
        + ".mp3"

    # Updating the item in DynamoDB
    response = table.update_item(
        Key={'id':post_id},
        UpdateExpression = "SET #statusAtt = :statusValue, #urlAtt = :urlValue",
        ExpressionAttributeValues = {
            ':statusValue': 'Updated',
            ':urlValue': url
        },
        ExpressionAttributeNames = {
            '#statusAtt': 'status',
            '#urlAtt': 'url'
        },
    )

    return response