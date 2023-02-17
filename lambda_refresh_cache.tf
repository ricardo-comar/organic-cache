data "archive_file" "lambda_refresh_cache_zip" {
  type        = "zip"
  source_file = "bin/lambda_refresh_cache"
  output_path = "bin/refresh_cache.zip"
}

// Function
resource "aws_lambda_function" "refresh_cache" {
  filename         = data.archive_file.lambda_refresh_cache_zip.output_path
  function_name    = "organic-cache-refresh_cache"
  description      = "Price Cache Refresh Lambda"
  role             = aws_iam_role.lambda_role_refresh_cache.arn
  handler          = "lambda_refresh_cache"
  source_code_hash = filebase64sha256(data.archive_file.lambda_refresh_cache_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role_refresh_cache]

  environment {
    variables = {
      ACTIVE_USERS_TABLE = aws_dynamodb_table.active_users.name
      REFRESH_QUEUE = aws_sqs_queue.refresh_queue.url
    }
  }
}



resource "aws_iam_role" "lambda_role_refresh_cache" {
  name               = "lambda_role_refresh_cache"
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

resource "aws_iam_policy" "iam_policy_for_lambda_refresh_cache" {

  name        = "aws_iam_policy_for_terraform_aws_lambda_role_refresh_cache"
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
        "lambda:InvokeFunction"
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
   }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role_refresh_cache" {
  role       = aws_iam_role.lambda_role_refresh_cache.name
  policy_arn = aws_iam_policy.iam_policy_for_lambda_refresh_cache.arn
}

resource "aws_cloudwatch_event_rule" "refresh_cache_lambda_event_rule" {
  name                = "refresh-cache-lambda-event-rule"
  description         = "retry scheduled every 1 min"
  schedule_expression = "rate(1 minutes)"
}

resource "aws_cloudwatch_event_target" "refresh_cache_lambda_target" {
  arn       = aws_lambda_function.refresh_cache.arn
  target_id = aws_lambda_function.refresh_cache.id
  rule      = aws_cloudwatch_event_rule.refresh_cache_lambda_event_rule.name
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_refresh_cache_lambda" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.refresh_cache.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.refresh_cache_lambda_event_rule.arn
}
