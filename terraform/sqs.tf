resource "aws_sqs_queue" "price_recalc_queue" {
  name                        = "organic-cache-sqs-price-recalc.fifo"
  max_message_size            = 20480
  message_retention_seconds   = 86400
  fifo_queue                  = true
  content_based_deduplication = true
  # delay_seconds = 1
}

resource "aws_sqs_queue" "quotation_queue" {
  name                      = "organic-cache-sqs-quotation"
  max_message_size          = 20480
  message_retention_seconds = 60
}
