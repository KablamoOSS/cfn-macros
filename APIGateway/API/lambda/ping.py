import json
import logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)


def main(event, context):
    response = {
        "statusCode": 200,
        "body": {
            "message": "pong",
            "response": {},
        }
    }

    response["body"] = json.dumps(response["body"])
    return response

