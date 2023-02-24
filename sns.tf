resource "aws_sns_topic" "user_subscriber" {
  name = "organic-cache-sns-user-subscriber"
}

resource "aws_sns_topic_subscription" "user_subscription_sns_target" {
  topic_arn              = aws_sns_topic.user_subscriber.arn
  protocol = "lambda"
  endpoint = aws_lambda_function.user_subscribe.arn
  # protocol               = "https"
  # endpoint_auto_confirms = true

  # endpoint               = "https://localhost:4566${aws_api_gateway_resource.org_cache_subscribe.path}"
  # depends_on = [
  #   aws_api_gateway_deployment.subscribe_rest_api_deploy
  # ]
}
