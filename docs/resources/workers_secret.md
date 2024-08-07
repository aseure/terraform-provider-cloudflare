---
page_title: "cloudflare_workers_secret Resource - Cloudflare"
subcategory: ""
description: |-
  Provides a Cloudflare Worker secret resource.
---

# cloudflare_workers_secret (Resource)

Provides a Cloudflare Worker secret resource.

## Example Usage

```terraform
resource "cloudflare_workers_secret" "my_secret" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  name        = "MY_EXAMPLE_SECRET_TEXT"
  script_name = "script_1"
  secret_text = "my_secret_value"
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_id` (String) The account identifier to target for the resource.
- `name` (String) The name of the Worker secret. **Modifying this attribute will force creation of a new resource.**
- `script_name` (String) The name of the Worker script to associate the secret with. **Modifying this attribute will force creation of a new resource.**
- `secret_text` (String, Sensitive) The text of the Worker secret. **Modifying this attribute will force creation of a new resource.**

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
$ terraform import cloudflare_workers_secret.example <account_id>/<script_name>/<secret_name>
```
