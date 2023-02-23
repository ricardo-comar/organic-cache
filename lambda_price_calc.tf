data "archive_file" "lambda_price_calc_zip" {
  type        = "zip"
  source_file = "bin/lambda_price_calc"
  output_path = "bin/price_calc.zip"
}

// Function
resource "aws_lambda_function" "price_calc" {
  filename         = data.archive_file.lambda_price_calc_zip.output_path
  function_name    = "identity-provider-price-calc"
  description      = "DB Update Lambda"
  role             = aws_iam_role.lambda_role_price_calc.arn
  handler          = "lambda_price_calc"
  source_code_hash = filebase64sha256(data.archive_file.lambda_price_calc_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]

  environment {
    variables = {
      PRODUCTS_TABLE = aws_dynamodb_table.products.name
      USER_DISCOUNTS_TABLE = aws_dynamodb_table.user_discounts.name
      USER_PRICES_TABLE = aws_dynamodb_table.user_prices.name
    }
  }

}

resource "aws_iam_role" "lambda_role_price_calc" {
  name               = "lambda_role_price_calc"
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

resource "aws_iam_policy" "iam_policy_for_lambda" {

  name        = "aws_iam_policy_for_terraform_aws_lambda_role_price_calc"
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
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes"
     ],
     "Resource": "*",
     "Effect": "Allow"
   },
   {
     "Action": [
        "dynamodb:Query",
        "dynamodb:Scan",
        "dynamodb:GetItem",
        "dynamodb:UpdateItem",
        "dynamodb:PutItem",
        "dynamodb:UpdateItem"
     ],
     "Resource": "*",
     "Effect": "Allow"
   }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role" {
  role       = aws_iam_role.lambda_role_price_calc.name
  policy_arn = aws_iam_policy.iam_policy_for_lambda.arn
}

# Event source from SQS
resource "aws_lambda_event_source_mapping" "price_calc_event_source_mapping" {
  event_source_arn = aws_sqs_queue.refresh_queue.arn
  enabled          = true
  function_name    = aws_lambda_function.price_calc.arn
  batch_size       = 10
}
