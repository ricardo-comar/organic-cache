resource "aws_dynamodb_table" "active_users" {
  name           = "active_users"
  billing_mode   = "PROVISIONED"
  read_capacity  = "2"
  write_capacity = "1"
  hash_key       = "user_id"

  attribute {
    name = "user_id"
    type = "S"
  }

  ttl {
    enabled        = true
    attribute_name = "ttl"
  }
}


resource "aws_dynamodb_table" "user_discounts" {
  name           = "user_discounts"
  billing_mode   = "PROVISIONED"
  read_capacity  = "2"
  write_capacity = "1"
  hash_key       = "user_id"

  attribute {
    name = "user_id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "products" {
  name           = "products"
  billing_mode   = "PROVISIONED"
  read_capacity  = "2"
  write_capacity = "1"
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "user_prices" {
  name           = "user_prices"
  billing_mode   = "PROVISIONED"
  read_capacity  = "2"
  write_capacity = "1"
  hash_key       = "user_id"
 
  attribute {
    name = "user_id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "quotations" {
  name           = "quotations"
  billing_mode   = "PROVISIONED"
  read_capacity  = "2"
  write_capacity = "1"
  hash_key       = "request_id"
 
  attribute {
    name = "request_id"
    type = "S"
  }

  ttl {
    enabled        = true
    attribute_name = "ttl"
  }

}
