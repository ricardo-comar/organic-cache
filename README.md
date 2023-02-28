# Organic Cache
Proof of Concept for a organic cache in AWS, offering a price quotation API that responds immediately when the price table is already available, and if not calculates on demand and saves the user ID to periodically recalculate and allows to be available on time.
If the user stops to consume the API after a pre determined time, that price table stops to be recalculated periodically, saving processing time.

## Solution Diagram
![Infrastructure Diagram](doc/diagram.png)
[Source](https://app.cloudcraft.co/view/c11241e7-f79b-42b3-b008-85ca557f501c?key=5618624e-2104-4aec-8a13-c1d94a20a96c)


## Scenarions
### User Subscription
1. On user first access, on calling ***Subscribe User PUT***, it's unique ID is registred on DynamoDB **Active Users** to be recalculated regurlaly, and a TTL (time-to-live) attribute.
1. If the result record stills un-updated until the TTL expires, it's automaticaly removed by DynamoDB and user price table will no longer be recalculated regurlaly.
1. If the same user is subscribed again, the TTL is updated with the new expiration time.

### Price Table
1. A scheduled CloudWatch Event Rule triggers the Lambda **Refresh Cache**, scans all items in DynamoDB **Active Users**, and for each of them generates a message to SQS **Refresh Queue**.
1. For each message on SQS **Refresh Queue**, Lambda **Price Calc** scan for all DynamoDB **Products** and query for user's DynamoDB **Discounts´**, calculating user's price table (with original and final prices) and saving into DybanoDB **Prices by User**. 

### Quotation Available
1. To request a quotation, the user can call the ***Quotation Websocket***, sending his ID and a list of products IDs and quantities. You can test it localy, as described below. 
1. The Lambda **Quotation Handler** saves the quotation data into DybamoDB **Quotations** and sends a message to SNS **Quotations Topic** with the quotation data, and returns a message *quotation under analisys*. 
1. The Lambda **Quotation Provider** is notified by SNS and, queries for the corresponding prices table from the receiving user's ID in DynamoDB **Prices by User**. 
1. Using the user's price table and products final prices and expected quantities, generates the quotation and notifies the client by it's *ConnectionID*, using API Gateway response channel.

### Quotation Unavailable
1. As described before, if Lambda **Quotation Provider** cannot find the user price table, it sends a message to SQS **Quotation Queue** with a different **group id**, to be processed before the refresh messages :sunglasses: and finish processing.
1. When Lambda **Price Calc** receives that message with a **RequestId** attribute, after price calculation it sends a message to SNS **Quotations Topic**, to Lambda **Quotation Provider** run again and be able to reply the quotation response.


## Local setup
### Installed programs:
- Docker
- Docker Compose
- Go
- Terraform
- Terraform Local
- AWS CLI

### IMPORTANT - Localstack Pro License

Because of the recent necessity to modify the solution to communicate by websocket, the API Gateway must be created using resource [apigatewayv2](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_api), but in localstack it's only available in Pro Version.

To be able to run localy, you can create an API Key in (https://app.localstack.cloud/) valid for 15 days for free.


### Usage
I recommend multiple terminals (like Ubuntu Terminator) to keep track of all running events

#### Terminal 1 - docker-compose
```
$ cd localstack
$ docker-compose up
```

#### Terminal 2 - terraform

First, check if the service is running by calling the health endpoint:
```
curl http://localhost:4566/_localstack/health | jq
```

Now you can deploy everything :metal:
```
$ make deploy
```
Keep this terminal on sight, copy the output ***url_quotation*** value to be used a few moments later... 

And load the prices and user discounts into DynamoDB:
```
$ make load
```

#### Terminal 3 - user subscribe
```
$ make subscribe ID=AAA
```
Now eyes on Terminal 1.... A few moments later Lambda **Refresh Cache** will be triggered and you be able to notice the user AAA price table be calculated.

#### Browser - Quotation

Open the file [test_page.html](test_page.html), open DevTools (F12) and switch to Console view.

Click on Setup and paste the ***url_quotation*** value into input box, and Click on Close. You should see a message "*opened*" on console view.

Sometimes after 60s you will see also a "pong!" message there... Don't mind, it's the javascript making the websocket active :smile:

On *Client ID* input, you can use any string (like Foo) and on console you will receive a first message "quotation under analisys", and right after another response with 3 products.

If you use AAA or BBB as *Client ID*, some products will be answered with a discount, present in [user_discounts.json](localstack/dynamodb_user_discounts.json). You can modify to test different results, loading the data on Terminal 2 again.


#### TIP
You can monitor a lambda individually by connecting into the corresponding docker process:
```
docker logs --follow $(docker ps -f name=price-calc -q)
```

## Conclusion
This project is for educational purpuse only, used to learn more about AWS and GoLang features (like parallel programming and circuit breaker strategy).