def handler(event, context):
    print("Hello from Lambda!")
    print(event)
    return {
        "statusCode": 200,
        "body": "Lambda executed successfully"
    }
