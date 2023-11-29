---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerplatform_billing_policies Data Source - powerplatform"
subcategory: ""
description: |-
  Fetches the list of billing policies in a tenant
---

# powerplatform_billing_policies (Data Source)

Fetches the list of billing policies in a tenant

## Example Usage

```terraform
terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  client_id = var.client_id
  secret    = var.secret
  tenant_id = var.tenant_id
}

data "powerplatform_billing_policies" "all_policies" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `billing_policies` (Attributes List) [Power Platform Billing Policy](https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/billing-policy/get-billing-policy#billingpolicyresponsemodel) (see [below for nested schema](#nestedatt--billing_policies))
- `id` (Number) Placeholder identifier attribute

<a id="nestedatt--billing_policies"></a>
### Nested Schema for `billing_policies`

Required:

- `billing_instrument` (Attributes) The billing instrument of the billing policy (see [below for nested schema](#nestedatt--billing_policies--billing_instrument))
- `location` (String) The location of the billing policy
- `name` (String) The name of the billing policy

Optional:

- `status` (String) The status of the billing policy (Enabled, Disabled)

Read-Only:

- `id` (String) The id of the billing policy

<a id="nestedatt--billing_policies--billing_instrument"></a>
### Nested Schema for `billing_policies.billing_instrument`

Required:

- `resource_group` (String) The resource group of the billing instrument
- `subscription_id` (String) The subscription id of the billing instrument

Read-Only:

- `id` (String) The id of the billing instrument