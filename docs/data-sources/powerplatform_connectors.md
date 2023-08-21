---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerplatform_connectors Data Source - terraform-provider-power-platform"
subcategory: ""
description: |-
  Fetches the list of available connectors in a Power Platform tenant
---

# powerplatform_connectors (Data Source)

Fetches the list of available connectors in a Power Platform tenant



<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `connectors` (Attributes List) List of Connectors (see [below for nested schema](#nestedatt--connectors))
- `id` (String) The ID of this resource.

<a id="nestedatt--connectors"></a>
### Nested Schema for `connectors`

Read-Only:

- `description` (String) Description
- `display_name` (String) Display name
- `id` (String) Id
- `name` (String) Name
- `publisher` (String) Publisher
- `tier` (String) Tier
- `type` (String) Type