resource "aws_iam_role" "ws_quotation_api_gateway_role" {
  name = "WsQuotationAPIGatewayRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "apigateway.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_apigatewayv2_api" "ws_quotation_api_gateway" {
  name                       = "ws-quotation-api-gateway"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_integration" "ws_quotation_api_integration" {
  api_id                    = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  integration_type          = "AWS_PROXY"
  integration_uri           = aws_lambda_function.quotation_handler.invoke_arn
  credentials_arn           = aws_iam_role.ws_quotation_api_gateway_role.arn
  content_handling_strategy = "CONVERT_TO_TEXT"
  passthrough_behavior      = "WHEN_NO_MATCH"
}

resource "aws_apigatewayv2_integration_response" "ws_quotation_api_integration_response" {
  api_id                   = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  integration_id           = aws_apigatewayv2_integration.ws_quotation_api_integration.id
  integration_response_key = "/200/"
}

resource "aws_apigatewayv2_route" "ws_quotation_api_default_route" {
  api_id    = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.ws_quotation_api_integration.id}"
}

resource "aws_apigatewayv2_route_response" "ws_quotation_api_default_route_response" {
  api_id             = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  route_id           = aws_apigatewayv2_route.ws_quotation_api_default_route.id
  route_response_key = "$default"
}

resource "aws_apigatewayv2_route" "ws_quotation_api_connect_route" {
  api_id    = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  route_key = "$connect"
  target    = "integrations/${aws_apigatewayv2_integration.ws_quotation_api_integration.id}"
}

resource "aws_apigatewayv2_route_response" "ws_quotation_api_connect_route_response" {
  api_id             = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  route_id           = aws_apigatewayv2_route.ws_quotation_api_connect_route.id
  route_response_key = "$default"
}

resource "aws_apigatewayv2_route" "ws_quotation_api_disconnect_route" {
  api_id    = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  route_key = "$disconnect"
  target    = "integrations/${aws_apigatewayv2_integration.ws_quotation_api_integration.id}"
}

resource "aws_apigatewayv2_route_response" "ws_quotation_api_disconnect_route_response" {
  api_id             = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  route_id           = aws_apigatewayv2_route.ws_quotation_api_disconnect_route.id
  route_response_key = "$default"
}

resource "aws_apigatewayv2_route" "ws_quotation_api_message_route" {
  api_id    = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  route_key = "MESSAGE"
  target    = "integrations/${aws_apigatewayv2_integration.ws_quotation_api_integration.id}"
}

resource "aws_apigatewayv2_route_response" "ws_quotation_api_message_route_response" {
  api_id             = aws_apigatewayv2_api.ws_quotation_api_gateway.id
  route_id           = aws_apigatewayv2_route.ws_quotation_api_message_route.id
  route_response_key = "$default"
}


# resource "aws_api_gateway_deployment" "quotation_rest_api_deploy" {
#   depends_on = [aws_api_gateway_integration.org_cache_quotation_integration]
#   rest_api_id = aws_api_gateway_rest_api.org_cache_api.id
#   stage_name  = "v1"
# }

# output "url_quotation" {
#   # value = "${aws_api_gateway_deployment.quotation_rest_api_deploy.invoke_url}${aws_api_gateway_resource.org_cache_quotation.path}"
#   value = "http://localhost:4566${aws_apigatewayv2_route.ws_quotation_api_default_route.api}"
# }

