---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerplatform_application Resource - powerplatform"
subcategory: ""
description: |-
  This resource allows you to install a Dynamics 365 application in an environment.
  This is functionally equivalent to the 'Install' button in the Power Platform admin center or pac application install in the Power Platform CLI https://docs.microsoft.com/en-us/powerapps/developer/data-platform/powerapps-cli#pac-application-install.  This resource uses the Install Application Package https://docs.microsoft.com/en-us/rest/api/power-platform/appmanagement/applications/installapplicationpackage endpoint in the Power Platform API.
  NOTE: This resource does not support updating or deleting applications.  The expected behavior is that the application is installed and remains installed until the environment is deleted.
---

# powerplatform_application (Resource)

This resource allows you to install a Dynamics 365 application in an environment.

This is functionally equivalent to the 'Install' button in the Power Platform admin center or [`pac application install` in the Power Platform CLI](https://docs.microsoft.com/en-us/powerapps/developer/data-platform/powerapps-cli#pac-application-install).  This resource uses the [Install Application Package](https://docs.microsoft.com/en-us/rest/api/power-platform/appmanagement/applications/installapplicationpackage) endpoint in the Power Platform API.

NOTE: This resource does not support updating or deleting applications.  The expected behavior is that the application is installed and remains installed until the environment is deleted.

## Example Usage

```terraform
terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

data "powerplatform_environments" "all_environments" {}

data "powerplatform_applications" "application_to_install" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  name           = "Power Platform Pipelines"
  publisher_name = "Microsoft Dynamics 365"
}

data "powerplatform_applications" "all_applications" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
}

locals {
  onboarding_essential_application = toset([for each in data.powerplatform_applications.all_applications.applications : each if each.application_name == "Onboarding essentials"])
}

resource "powerplatform_application" "install_sample_application" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  unique_name    = data.powerplatform_applications.application_to_install.applications[0].unique_name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) Id of the Dynamics 365 environment
- `unique_name` (String) Unique name of the application

### Read-Only

- `id` (String) Unique id (guid)