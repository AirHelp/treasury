from __future__ import print_function
import os
import time
import json
import decimal
import boto3
from botocore.exceptions import ClientError


def write(event, context):
    try:
        input = json.loads(event['body'])
        key = input['key']
        value = input['value']
        kms_arn = input['kms_arn']
        author = input['author']
    except KeyError:
        print("ERROR in given data. Event:{}".format(event))
        body = {
            "message": "Missing required fields. Check documentation."
        }
        return {
            "statusCode": 400,
            "body": json.dumps(body)
        }
    except ValueError as e:
        print("ERROR in given data. Event:{}".format(event))
        body = {
            "message": "{}".format(e)
        }
        return {
            "statusCode": 400,
            "body": json.dumps(body)
        }

    data = {
        'key': key,
        'version': decimal.Decimal(time.time()),
        'kms_arn': kms_arn,
        'value': value,
        'author': author
    }

    print("Try to write data:{} into dynamo".format(data))

    try:
        dynamodb = boto3.resource('dynamodb')
        dynamo_table = dynamodb.Table(os.environ['DYNAMODB_TABLE'])
        dynamo_response = dynamo_table.put_item(Item=data)
    except ClientError as e:
        print("DynamoDB ERROR:{}".format(e.response))
        body = {
            "message": e.response['Error']['Message']
        }
        return {
            "statusCode": 500,
            "body": json.dumps(body)
        }

    print("Success! DynamoDB respose:{}".format(dynamo_response))
    body = {
        "message": "Success! Data written to: {}".format(key)
    }
    return {
        "statusCode": 200,
        "body": json.dumps(body)
    }
