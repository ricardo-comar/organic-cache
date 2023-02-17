resource "aws_dynamodb_table" "active_users" {
  name           = "active_users"
  billing_mode   = "PROVISIONED"
  read_capacity  = "2"
  write_capacity = "1"
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  ttl {
    enabled        = true
    attribute_name = "ttl"
  }

  #  point_in_time_recovery { enabled = true } 
  # server_side_encryption { enabled = true } 
}
