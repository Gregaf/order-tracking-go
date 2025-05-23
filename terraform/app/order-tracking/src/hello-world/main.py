import json

def handler(event, ctx):
    message = f"Hello world!"

    return {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "message": message
        })
    }