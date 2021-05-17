module "discounts_products_crud_sqs_queue" {
  source = "../modules//sqs"
  name = "Discounts-ProductsCrud.fifo"
  fifo_queue = true
  content_based_deduplication = true

  tags = {
    Name = "discounts_products_sqs_queue"
    project = var.project
    env = var.account
  }
}

data "aws_sns_topic" "products_events_sns" {
  name = "ProductsEvents.fifo"
  provider = aws
}

resource "aws_sns_topic_subscription" "discounts_products_crud_subscription_queue" {
  topic_arn = data.aws_sns_topic.products_events_sns.arn
  protocol = "sqs"
  endpoint = module.discounts_products_crud_sqs_queue.this_sqs_queue_arn
  raw_message_delivery = true
  filter_policy = jsonencode({
    "type": [
      "products"],
    "operation": [
      "create",
      "update",
      "delete"],
  })
}

module "discounts_users_crud_sqs_queue" {
  source = "../modules//sqs"
  name = "Discounts-UsersCrud.fifo"
  fifo_queue = true
  content_based_deduplication = true

  tags = {
    Name = "discounts_users_sqs_queue"
    project = var.project
    env = var.account
  }
}

data "aws_sns_topic" "users_events_sns" {
  name = "UsersEvents.fifo"
  provider = aws
}

resource "aws_sns_topic_subscription" "discounts_users_crud_subscription_queue" {
  topic_arn = data.aws_sns_topic.users_events_sns.arn
  protocol = "sqs"
  endpoint = module.discounts_users_crud_sqs_queue.this_sqs_queue_arn
  raw_message_delivery = true
  filter_policy = jsonencode({
    "type": [
      "users"],
    "operation": [
      "create",
      "update",
      "delete"],
  })
}
