
resource "aws_api_gateway_rest_api" "org_cache_api" {
  name = "org_cache_api"
}

resource "aws_api_gateway_rest_api_policy" "rest_api_policy" {
  rest_api_id = aws_api_gateway_rest_api.org_cache_api.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": "*",
            "Action": "execute-api:Invoke",
            "Resource": [
                "${aws_api_gateway_rest_api.org_cache_api.execution_arn}/*"
            ]
        }
    ]
}
EOF

}

############################# SUBSCRIBE API ############################# 

resource "aws_api_gateway_resource" "org_cache_subscribe" {
  path_part   = "subscribe"
  parent_id   = aws_api_gateway_rest_api.org_cache_api.root_resource_id
  rest_api_id = aws_api_gateway_rest_api.org_cache_api.id
}

resource "aws_api_gateway_method" "org_cache_subscribe_method" {
  rest_api_id   = aws_api_gateway_rest_api.org_cache_api.id
  resource_id   = aws_api_gateway_resource.org_cache_subscribe.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "org_cache_subscribe_integration" {
  rest_api_id             = aws_api_gateway_rest_api.org_cache_api.id
  resource_id             = aws_api_gateway_resource.org_cache_subscribe.id
  http_method             = aws_api_gateway_method.org_cache_subscribe_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.user_subscribe.invoke_arn
}

resource "aws_api_gateway_deployment" "subscribe_rest_api_deploy" {
  depends_on = [aws_api_gateway_integration.org_cache_subscribe_integration]
  rest_api_id = aws_api_gateway_rest_api.org_cache_api.id
  stage_name  = "v1"
}

output "url_subscribe" {
  value = "${aws_api_gateway_deployment.subscribe_rest_api_deploy.invoke_url}${aws_api_gateway_resource.org_cache_subscribe.path}"
}

############################# QUOTATION API ############################# 

resource "aws_api_gateway_resource" "org_cache_quotation" {
  path_part   = "quotation"
  parent_id   = aws_api_gateway_rest_api.org_cache_api.root_resource_id
  rest_api_id = aws_api_gateway_rest_api.org_cache_api.id
}

resource "aws_api_gateway_method" "org_cache_quotation_method" {
  rest_api_id   = aws_api_gateway_rest_api.org_cache_api.id
  resource_id   = aws_api_gateway_resource.org_cache_quotation.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "org_cache_quotation_integration" {
  rest_api_id             = aws_api_gateway_rest_api.org_cache_api.id
  resource_id             = aws_api_gateway_resource.org_cache_quotation.id
  http_method             = aws_api_gateway_method.org_cache_quotation_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.quotation_handler.invoke_arn
}

resource "aws_api_gateway_deployment" "quotation_rest_api_deploy" {
  depends_on = [aws_api_gateway_integration.org_cache_quotation_integration]
  rest_api_id = aws_api_gateway_rest_api.org_cache_api.id
  stage_name  = "v1"
}

output "url_quotation" {
  value = "${aws_api_gateway_deployment.quotation_rest_api_deploy.invoke_url}${aws_api_gateway_resource.org_cache_quotation.path}"
}

