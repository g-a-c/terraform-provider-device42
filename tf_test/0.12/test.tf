data "device42_password" "test" {
  # uses one of two Sumologic syslog tokens to demo the value changing
  # and that the values can be used to change resources, change outputs, etc

  # id = 18756
  id = 18759
}

resource "null_resource" "example" {
  triggers = {
    value = data.device42_password.test.password
  }
}

output "output-value" {
  value = data.device42_password.test.password
}

output "output-label" {
  value = data.device42_password.test.label
}

output "output-username" {
  value = data.device42_password.test.username
}