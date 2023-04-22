# Data Source: bitbucket_ip_ranges
Use this data source to get a list of IP Address information for Bitbucket Cloud, this can be used as part of allow-lists.
See here for more information: https://support.atlassian.com/bitbucket-cloud/docs/what-are-the-bitbucket-cloud-ip-addresses-i-should-use-to-configure-my-corporate-firewall/

## Example Usage
```hcl
data "bitbucket_ip_ranges" "example" {}
```

## Argument Reference
No arguments are passed in.

## Attribute Reference
The following attributes are exported:
* `ranges` - A list of IP Address information, of which each entry in the list contains:
  * `network` - The IP Address.
  * `cidr` - The CIDR.
  * `mask` - The Mask.
  * `mask_len` - The Mask Length.
  * `regions` - A list of regions this IP Address resides in. Follows AWS Region Code Notation.
  * `directions` - A list defining if the IP address is ingress, egress or both.
