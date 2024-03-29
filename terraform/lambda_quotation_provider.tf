data "archive_file" "lambda_quotation_provider_zip" {
  type        = "zip"
  source_file = "../bin/lambda_quotation_provider"
  output_path = "../bin/lambda_quotation_provider.zip"
}

// Function
resource "aws_lambda_function" "quotation_provider" {
  filename         = data.archive_file.lambda_quotation_provider_zip.output_path
  function_name    = "quotation_provider"
  description      = "DB Update Lambda"
  role             = aws_iam_role.lambda_role_quotation_provider.arn
  handler          = "lambda_quotation_provider"
  source_code_hash = filebase64sha256(data.archive_file.lambda_quotation_provider_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]

  environment {
    variables = {
      USER_PRICES_TABLE = aws_dynamodb_table.user_prices.name
      QUOTATIONS_TABLE  = aws_dynamodb_table.quotations.name
      RECALC_QUEUE      = aws_sqs_queue.price_recalc_queue.url
      QUOTATIONS_QUEUE  = aws_sqs_queue.quotation_queue.url
    }
  }

}

resource "aws_iam_role" "lambda_role_quotation_provider" {
  name               = "lambda_role_quotation_provider"
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

resource "aws_iam_policy" "iam_policy_for_lambda_quotation_provider" {

  name        = "aws_iam_policy_for_terraform_aws_lambda_role_quotation_provider"
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
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes"
     ],
     "Resource": "*",
     "Effect": "Allow"
   },
   {
     "Action": [
        "execute-api:*"
      ],
     "Resource": "*",
     "Effect": "Allow"
   }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role_provider" {
  role       = aws_iam_role.lambda_role_quotation_provider.name
  policy_arn = aws_iam_policy.iam_policy_for_lambda_quotation_provider.arn
}
