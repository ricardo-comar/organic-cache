data "archive_file" "lambda_quotation_handler_zip" {
  type        = "zip"
  source_file = "bin/lambda_quotation_handler"
  output_path = "bin/quotation_handler.zip"
}

// Function
resource "aws_lambda_function" "quotation_handler" {
  filename         = data.archive_file.lambda_quotation_handler_zip.output_path
  function_name    = "organic-cache-quotation-handler"
  description      = "REST API to reply quotations"
  role             = aws_iam_role.lambda_role_quotation_handler.arn
  handler          = "lambda_quotation_handler"
  source_code_hash = filebase64sha256(data.archive_file.lambda_quotation_handler_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role_quotation_handler, aws_sns_topic.user_subscriber]

  environment {
    variables = {
      USER_SUBSCRIBE_TOPIC = aws_sns_topic.user_subscriber.arn
      QUOTATION_QUEUE = aws_sqs_queue.quotation_queue.url
      QUOTATION_TABLE = aws_dynamodb_table.quotations.name
    }
  }

}


resource "aws_iam_role" "lambda_role_quotation_handler" {
  name               = "lambda_role_quotation_handler"
  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "lambda.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF
}

resource "aws_iam_policy" "iam_policy_for_lambda_quotation_handler" {

  name        = "aws_iam_policy_for_terraform_aws_lambda_role_quotation_handler"
  path        = "/"
  description = "AWS IAM Policy for managing aws lambda role"
  policy      = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": [
       "logs:CreateLogGroup",
       "logs:CreateLogStream",
       "logs:PutLogEvents"
     ],
     "Resource": "arn:aws:logs:*:*:*",
     "Effect": "Allow"
   },
   {
     "Action": [
        "ec2:CreateNetworkInterface",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DeleteNetworkInterface"
     ],
     "Resource": "*",
     "Effect": "Allow"
   },
   {
     "Action": [
        "sns:Subscribe",
        "sns:Receive"
     ],
     "Resource": "*",
     "Effect": "Allow"
   },
   {
     "Action": [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes"
     ],
     "Resource": "*",
     "Effect": "Allow"
   },
   {
     "Action": [
        "dynamodb:PutItem"
     ],
     "Resource": "*",
     "Effect": "Allow"
   }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role_quotation_handler" {
  role       = aws_iam_role.lambda_role_quotation_handler.name
  policy_arn = aws_iam_policy.iam_policy_for_lambda_quotation_handler.arn
}

resource "aws_lambda_permission" "apigw_lambda_quotation_handler" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.quotation_handler.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.org_cache_api.execution_arn}/*/*/*"
}
