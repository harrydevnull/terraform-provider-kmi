---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kmi Provider"
subcategory: ""
description: |-
  
---

# kmi Provider



## Example Usage

```terraform
provider "kmi" {
  host      = "kdc host"
  api_key   = "<path to api_cert.key>"
  api_crt   = "<path to api_cert.crt>"
  akamai_ca = "<path to akamai_ca_list.pem>"

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `akamai_ca` (String)
- `akamai_ca_path` (String)
- `api_crt` (String)
- `api_crt_path` (String)
- `api_key` (String, Sensitive)
- `api_key_path` (String, Sensitive)
- `host` (String)
- `proxy_host` (String)
