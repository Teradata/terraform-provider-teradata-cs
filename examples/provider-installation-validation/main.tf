terraform {
  required_providers {
    teradata-clearscape = {
      source = "hashicorp.com/edu/teradata-clearscape"
    }
  }
}

provider "teradata-clearscape" {
  token = "***REMOVED***"
}

resource "teradata-clearscape_environment" "edu1" {
  name = "terraformtest"
  region = "us-central"
  password = "terraformtest"
}

output "edu1_environment" {
  value = teradata-clearscape_environment.edu1
   sensitive = true
}
  