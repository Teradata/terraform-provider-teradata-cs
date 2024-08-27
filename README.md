<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for Teradata ClearScape

The Teradata ClearScape Terraform Provider allows managing resources within ClearScape platform.

## Usage Example

```hcl
# 1. Specify the version of the AzureRM Provider to use
terraform {
  required_providers {
    teradata-clearscape = {
      source = "hashicorp.com/edu/teradata-clearscape"
    }
  }
}

# 2. Configure the Teradata ClearScape Provider

provider "teradata-clearscape" {
  token = "CLEARSCAPE_TOKEN"
}

# 3. Create a resource group
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "teradata-clearscape_environment" "example" {
  name = "example-resource"
  region = "example"
  password = "sensitive"
}


```

* [Additional examples can be found in the `./examples` folder within this repository](https://github.com/teradata/terraform-provider-teradata-clearscape/tree/main/examples).