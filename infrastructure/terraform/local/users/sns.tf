resource "aws_sns_topic" "users_events" {
  name = "UsersEvents.fifo"
  fifo_topic = true

  tags = {
    Name = "users_sns_topic"
    project = var.project
    env = var.account
  }
}
