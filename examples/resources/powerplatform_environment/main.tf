terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "github.com/microsoft/terraform-provider-power-platform"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  tenant_id = var.tenant_id
}

resource "powerplatform_environment" "development" {
  display_name     = "example_environment"
  location         = "europe"
  language_code    = "1033"
  currency_code    = "USD"
  environment_type = "Sandbox"
  domain           = "mydomain"
  security_group_id = "00000000-0000-0000-0000-000000000000"
}

