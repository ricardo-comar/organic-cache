resource "aws_sns_topic" "user_prices_topic" {
  name = "organic-cache-sns-quotations-topic"
}

resource "aws_sns_topic_subscription" "user_prices_topic_sns_target" {
  topic_arn = aws_sns_topic.user_prices_topic.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.quotation_provider.arn
}