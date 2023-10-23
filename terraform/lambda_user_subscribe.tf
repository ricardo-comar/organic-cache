data "archive_file" "lambda_user_subscribe_zip" {
  type        = "zip"
  source_file = "../bin/lambda_user_subscribe"
  output_path = "../bin/lambda_user_subscribe.zip"
}

// Function
resource "aws_lambda_function" "user_subscribe" {
  filename         = data.archive_file.lambda_user_subscribe_zip.output_path
  function_name    = "user_subscribe"
  description      = "REST API to subscribe users"
  role             = aws_iam_role.lambda_role_user_subscribe.arn
  handler          = "lambda_user_subscribe"
  source_code_hash = filebase64sha256(data.archive_file.lambda_user_subscribe_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role_user_subscribe]

  environment {
    variables = {
      ACTIVE_USERS_TABLE = aws_dynamodb_table.active_users.name
      RECALC_QUEUE       = aws_sqs_queue.price_recalc_queue.url
    }
  }

}


resource "aws_iam_role" "lambda_role_user_subscribe" {
  name               = "lambda_role_user_subscribe"
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

resource "aws_iam_policy" "iam_policy_for_lambda_user_subscribe" {

  name        = "aws_iam_policy_for_terraform_aws_lambda_role_user_subscribe"
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
        "sqs:SendMessage",
        "sqs:GetQueueAttributes"
     ],
     "Resource": "*",
     "Effect": "Allow"
   },
   {
     "Action": [
        "dynamodb:GetItem",
        "dynamodb:PutItem"
     ],
     "Resource": "*",
     "Effect": "Allow"
   }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role_user_subscribe" {
  role       = aws_iam_role.lambda_role_user_subscribe.name
  policy_arn = aws_iam_policy.iam_policy_for_lambda_user_subscribe.arn
}

resource "aws_lambda_permission" "apigw_lambda_user_subscribe" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.user_subscribe.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.org_cache_api.execution_arn}/*/*/*"
}
