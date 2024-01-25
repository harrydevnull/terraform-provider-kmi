---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kmi_template Resource - terraform-provider-kmi"
subcategory: ""
description: |-
  
---

# kmi_template (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ca_collection` (String) CA collection name to be created on KMI
- `ca_definition` (String) CA definition name to be created on KMI
- `client_collection` (String) Client collection name to be created on KMI
- `options` (Attributes) (see [below for nested schema](#nestedatt--options))
- `template_name` (String) Certificate Signing Request template name to be created on KMI

### Read-Only

- `last_updated` (String) The last time the group was updated.

<a id="nestedatt--options"></a>
### Nested Schema for `options`

Optional:

- `allow_ca` (String) Whether the signed secret can have the CA option set in the BasicConstraints extension.
- `common_name` (String) Common name of the certificate. Can be "*" to allow all values, or a string with '*' as a glob character
- `dns_san` (String) Comma delimited list of acceptable domain names for the Subject Alternative Name extension. Names use '*' as a glob character Default no values allowed
- `hash_type` (String) Comma delimited list of acceptable hash_types for the signed certificate. Can be '*' to allow all key types. This constraint is ignored for key_types that don't use hashing as part of the signature (ed25519)
- `ip_san` (String) Comma delimited list of acceptable IPs for the Subject Alternative Name extension. Can be IP addresses or CIDRs. Default no values allowed
- `key_type` (String) Comma delimited list of acceptable key types for the signed certificate. Can be '*' to allow all key types.Default rsa:2048,rsa:4096,ec:secp256r1
- `leaf_exceeds_ca_ttl` (String) Boolean flag as to whether or not the signed secret's notAfter date can exceed that of the CA certificate
- `max_ttl` (String) The maximum period of time that the signed secret can be valid for Default is 90 days
- `min_ttl` (String) The minimum period of time that the signed secret can be valid for Default is 7 day
- `uri_san` (String) Comma delimited list of acceptable URIs for the Subject Alternative Name extension. Names use '*' as a glob character. Default no values allowed