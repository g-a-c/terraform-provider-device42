terraform {
  required_providers {
    device42 = {
      source  = "terraform.example.com/example/device42"
      version = "0.0.1"
    }
  }
}

data "device42_password" "test" {
  id = 4
}

resource "null_resource" "example" {
  triggers = {
    value = data.device42_password.test.password
  }
}

output "output-password" {
  value = data.device42_password.test.password
}

output "output-label" {
  value = data.device42_password.test.label
}

output "output-username" {
  value = data.device42_password.test.username
}
