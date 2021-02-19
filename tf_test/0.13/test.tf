terraform {
  required_providers {
    device42 = {
      source  = "g-a-c/device42"
      version = "~>1.0.0"
    }
  }
}

variable "secret_id" {
  type        = number
  description = "what is the number of the secret to retrieve from Device42?"
}

data "device42_password" "test" {
  id = var.secret_id
}

resource "null_resource" "example" {
  triggers = {
    value = data.device42_password.test.password
  }
}

# resource "random_password" "test_pw" {
#   length  = 8
#   special = true
# }

output "output-password" {
  value = data.device42_password.test.password
}

output "output-label" {
  value = data.device42_password.test.label
}

output "output-username" {
  value = data.device42_password.test.username
}
