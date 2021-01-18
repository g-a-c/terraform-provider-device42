variable "secret_id" {
  type        = number
  description = "ID of the secret to retrieve from Device42"
}

data "device42_password" "test" {
  id = var.secret_id
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
