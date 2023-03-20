data "archive_file" "lambda_cache_refresh_zip" {
  type        = "zip"
  source_file = "bin/lambda_cache_refresh"
  output_path = "bin/cache_refresh.zip"
}

// Function
resource "aws_lambda_function" "cache_refresh" {
  filename         = data.archive_file.lambda_cache_refresh_zip.output_path
  function_name    = "cache_refresh"
  description      = "Price Cache Refresh Lambda"
  role             = aws_iam_role.lambda_role_cache_refresh.arn
  handler          = "lambda_cache_refresh"
  source_code_hash = filebase64sha256(data.archive_file.lambda_cache_refresh_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role_cache_refresh]

  environment {
    variables = {
      ACTIVE_USERS_TABLE = aws_dynamodb_table.active_users.name
      RECALC_QUEUE       = aws_sqs_queue.price_recalc_queue.url
    }
  }
}



resource "aws_iam_role" "lambda_role_cache_refresh" {
  name               = "lambda_role_cache_refresh"
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

resource "aws_iam_policy" "iam_policy_for_lambda_cache_refresh" {

  name        = "aws_iam_policy_for_terraform_aws_lambda_role_cache_refresh"
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

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role_cache_refresh" {
  role       = aws_iam_role.lambda_role_cache_refresh.name
  policy_arn = aws_iam_policy.iam_policy_for_lambda_cache_refresh.arn
}

resource "aws_cloudwatch_event_rule" "cache_refresh_lambda_event_rule" {
  name                = "cache-refresh-lambda-event-rule"
  description         = "retry scheduled every 1 min"
  schedule_expression = "rate(1 minutes)"
}

resource "aws_cloudwatch_event_target" "cache_refresh_lambda_target" {
  arn       = aws_lambda_function.cache_refresh.arn
  target_id = aws_lambda_function.cache_refresh.id
  rule      = aws_cloudwatch_event_rule.cache_refresh_lambda_event_rule.name
}

# resource "aws_lambda_permission" "allow_cloudwatch_to_call_cache_refresh_lambda" {
#   statement_id  = "AllowExecutionFromCloudWatch"
#   action        = "lambda:InvokeFunction"
#   function_name = aws_lambda_function.cache_refresh.function_name
#   principal     = "events.amazonaws.com"
#   source_arn    = aws_cloudwatch_event_rule.cache_refresh_lambda_event_rule.arn
# }
