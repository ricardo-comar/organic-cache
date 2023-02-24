resource "aws_sqs_queue" "refresh_queue" {
  name = "organic-cache-sqs-cache-refresh"
  max_message_size          = 20480
  message_retention_seconds = 86400
  delay_seconds = 1
}

resource "aws_sqs_queue" "quotation_queue" {
  name = "organic-cache-sqs-quotation"
  max_message_size          = 20480
  message_retention_seconds = 86400
}
