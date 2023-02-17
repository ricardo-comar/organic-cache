
resource "aws_sqs_queue" "refresh_queue" {
  name = "organic-cache-sqs-cache-refresh"
  max_message_size          = 20480
  message_retention_seconds = 86400
  receive_wait_time_seconds = 1
}
