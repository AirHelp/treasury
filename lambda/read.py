from __future__ import print_function
import os
import json
import boto3
import decimal
from boto3.dynamodb.conditions import Key, Attr


# Helper class to convert a DynamoDB item to JSON.
class DecimalEncoder(json.JSONEncoder):
    def default(self, o):
        if isinstance(o, decimal.Decimal):
            if o % 1 > 0:
                return float(o)
            else:
                return int(o)
        return super(DecimalEncoder, self).default(o)


def read(event, context):
    try:
        key = event['queryStringParameters']['key']
    except KeyError as e:
        print("ERROR in given data. Event:{}".format(event))
        body = {
            "message": "Missing {} in queryStringParameters.".format(e)
        }
        return {
            "statusCode": 400,
            "body": json.dumps(body)
        }

    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table(os.environ['DYNAMODB_TABLE'])

    dynamodb_response = table.query(
        Limit=1,
        ScanIndexForward=False,
        ConsistentRead=True,
        KeyConditionExpression=Key('key').eq(key)
    )

    if dynamodb_response["ResponseMetadata"]["HTTPStatusCode"] == 200:
        return {
            "statusCode": 200,
            "body": json.dumps(dynamodb_response['Items'][0], cls=DecimalEncoder)
        }
    else:
        print("ERROR in DynamoDB response:{}".format(dynamodb_response))
        return {
            "statusCode": 500,
            "body": json.dumps({"message": "Unable to get data from database."})
        }
