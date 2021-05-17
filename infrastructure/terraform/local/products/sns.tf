resource "aws_sns_topic" "products_events" {
  name = "ProductsEvents.fifo"
  fifo_topic = true
  tags = {
    Name = "products_sns_topic"
    project = var.project
    env = var.account
  }
}
