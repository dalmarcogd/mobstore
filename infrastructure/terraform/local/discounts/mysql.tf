provider "mysql" {
  endpoint = "localhost:3306"
  username = "root"
  password = "mysql"
}

resource "mysql_database" "mysql_database_discounts" {
  provider = mysql
  name = var.project
  lifecycle {
    prevent_destroy = true
  }
}

resource "mysql_user" "mysql_user_discounts" {
  provider = mysql
  user = var.project
  host = "%"
  plaintext_password = "my-password"
  lifecycle {
    prevent_destroy = true
  }
}

resource "mysql_grant" "mysql_grant_discounts" {
  provider = mysql

  user = mysql_user.mysql_user_discounts.user
  host = mysql_user.mysql_user_discounts.host
  database = mysql_database.mysql_database_discounts.name
  privileges = [
    "SELECT",
    "INSERT",
    "UPDATE",
    "DELETE",
    "CREATE",
    "DROP",
    "INDEX"]
  lifecycle {
    prevent_destroy = true
  }
}