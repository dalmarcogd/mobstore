provider "mysql" {
  endpoint = "localhost:3306"
  username = "root"
  password = "mysql"
}

resource "mysql_database" "mysql_database_users" {
  provider = mysql
  name = var.project
  lifecycle {
    prevent_destroy = true
  }
}

resource "mysql_user" "mysql_user_users" {
  provider = mysql
  user = var.project
  plaintext_password = "my-password"
  host = "%"
  lifecycle {
    prevent_destroy = true
  }
}

resource "mysql_grant" "mysql_grant_users" {
  provider = mysql
  user = mysql_user.mysql_user_users.user
  host = mysql_user.mysql_user_users.host
  database = mysql_database.mysql_database_users.name
  privileges = [
    "SELECT",
    "INSERT",
    "UPDATE",
    "DELETE",
    "CREATE",
    "DROP"]
  lifecycle {
    prevent_destroy = true
  }
}