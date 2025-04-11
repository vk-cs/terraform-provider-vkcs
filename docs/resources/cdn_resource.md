---
subcategory: "CDN"
layout: "vkcs"
page_title: "vkcs: vkcs_cdn_resource"
description: |-
  Manages a CDN resource within VKCS.
---

# vkcs_cdn_resource



## Example Usage
```terraform
resource "vkcs_cdn_resource" "resource" {
  cname        = local.cname
  origin_group = vkcs_cdn_origin_group.origin_group.id
  options = {
    edge_cache_settings = {
      value = "10m"
    }
    forward_host_header = true
    gzip_on             = true
  }
  ssl_certificate = {
    type = "own"
    id   = vkcs_cdn_ssl_certificate.certificate.id
  }
}
```

## Argument Reference
- `cname` **required** *string* &rarr;  Delivery domain that will be used for content delivery through a CDN. Use `secondary_hostnames` to add extra domains. <br>**Note:** Delivery domains should be added to your DNS settings.

- `origin_group` **required** *number* &rarr;  Origin group ID with which the CDN resource is associated.

- `active` optional *boolean* &rarr;  Enables or disables a CDN resource.

- `options` optional &rarr;  Options that configure a CDN resource.
  - `allowed_http_methods` optional &rarr;  HTTP methods allowed for content requests from the CDN.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *string* &rarr;  List of HTTP methods.


  - `brotli_compression` optional &rarr;  Compresses content with Brotli on the CDN side. CDN servers will request only uncompressed content from the origin.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *set of* *string* &rarr;  List of content types to be compressed. It's required to specify text/html here.


  - `browser_cache_settings` optional &rarr;  Cache settings for users browsers.
Cache expiration time is applied to the following response codes: 200, 201, 204, 206, 301, 302, 303, 304, 307, 308. Responses with other codes will not be cached.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *string* &rarr;  Cache expiration time. Use '0s' to disable caching.


  - `cors` optional &rarr;  Enables or disables CORS (Cross-Origin Resource Sharing) header support.
CORS header support allows the CDN to add the Access-Control-Allow-Origin header to a response to a browser.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *string* &rarr;  Value of the Access-Control-Allow-Origin header.
Possible values:
* ["*"] - adds "*" as the header value, content will be uploaded for requests from any domain.
*["domain.com", "second.dom.com"] - adds "$http_origin" as the header value if the origin matches one of the listed domains, content will be uploaded only for requests from the domains specified in the field.
* ["$http_origin"] - adds "$http_origin" as the header value, content will be uploaded for requests from any domain, and the domain from which the request was sent will be added to the header in the response.


  - `country_acl` optional &rarr;  Use this option to control access to the content for specified countries.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `excepted_values` optional *string* &rarr;  List of countries according to ISO-3166-1. The meaning of the argument depends on `policy_type` value.

    - `policy_type` optional *string* &rarr;  The type of CDN resource access policy. Must be one of following: "allow", "deny".


  - `edge_cache_settings` optional &rarr;  Cache settings for CDN servers.
    - `custom_values` optional *map of* *string* &rarr;  A map representing the caching time for a response with a specific response code.
These settings have a higher priority than the value field.
* Use `any` key to specify caching time for all response codes.
* Use `0s` value to disable caching for a specific response code.

    - `default` optional *string* &rarr;  Enables content caching according to the origin cache settings.

    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *string* &rarr;  Caching time. The value is applied to the following response codes: 200, 206, 301, 302.
Responses with codes 4xx, 5xx will not be cached. Use `0s` to disable caching.


  - `fetch_compressed` optional *boolean* &rarr;  If enabled, CDN servers request and cache compressed content from the origin. The origin server should support compression. CDN servers will not decompress your content even if a user browser does not accept compression. Conflicts with `gzip_on` if both enabled simultaneously.

  - `force_return` optional &rarr;  Allows to apply custom HTTP code to the CDN content.
Specify HTTP-code you need and text or URL if you're going to set up redirection.
    - `body` optional *string* &rarr;  URL for redirection or text.

    - `code` optional *number* &rarr;  Status code value.

    - `enabled` optional *boolean* &rarr;  Controls the option state.


  - `forward_host_header` optional *boolean* &rarr;  Forwards the Host header from a end-user request to an origin server. Conflicts with `host_header` if both enabled simultaneously.

  - `gzip_compression` optional &rarr;  Compresses content with GZip on the CDN side. CDN servers will request only uncompressed content from the origin. Conflicts with `fetch_compressed`, `slice` and `gzip_on` if any of them are enabled simultaneously. `application/wasm` value is not supported when the shielded option is disabled, compression in this case is performed on the origin shielding, so it must be active for the MIME type to be compressed.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *set of* *string* &rarr;  List of content types to be compressed. It's required to specify text/html here.


  - `gzip_on` optional *boolean* &rarr;  Enables content compression using gzip on the CDN side. CDN servers will request only uncompressed content from the origin. Conflicts with `fetch_compressed`, `slice` and `gzip_compression` if any of them are enabled simultaneously.

  - `host_header` optional &rarr;  Use this option to specify the Host header that CDN servers use when request content from an origin server. If the option is not set, the header value is equal to the first CNAME. Conflicts with `forward_host_header` if both enabled simultaneously.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *string* &rarr;  Host Header value.


  - `ignore_cookie` optional *boolean* &rarr;  Defines whether the files with the Set-Cookies header are cached as one file or as different ones.

  - `ignore_query_string` optional *boolean* &rarr;  Allows to specify how a file with different query strings is cached: either as one object (option is enabled) or as different objects (option is disabled.). `ignore_query_string`, `query_params_whitelist` and `query_params_blacklist` options cannot be enabled simultaneously.

  - `ip_address_acl` optional &rarr;  The option allows to control access to the CDN Resource content for specific IP addresses.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `excepted_values` optional *string* &rarr;  List of IP addresses with a subnet mask. The meaning of the argument depends on `policy_type` value.

    - `policy_type` optional *string* &rarr;  The type of CDN resource access policy. Must be one of following: "allow", "deny".


  - `query_params_blacklist` optional &rarr;  Use this option to specify query parameters, so files with these query strings will be cached as one object, and files with other parameters will be cached as different objects. `ignore_query_string`, `query_params_whitelist` and `query_params_blacklist` options cannot be enabled simultaneously.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *string* &rarr;  List of query parameters.


  - `query_params_whitelist` optional &rarr;  Use this option to specify query parameters, so files with these query strings will be cached as different objects, and files with other parameters will be cached as one object. `ignore_query_string`, `query_params_whitelist` and `query_params_blacklist` options cannot be enabled simultaneously.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *string* &rarr;  List of query parameters.


  - `referrer_acl` optional &rarr;  Use this option to control access to the CDN resource content for specified domain names.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `excepted_values` optional *string* &rarr;  List of domain names or wildcard domains, without protocol. The meaning of the argument depends on `policy_type` value.

    - `policy_type` optional *string* &rarr;  The type of CDN resource access policy. Must be one of following: "allow", "deny".


  - `secure_key` optional &rarr;  Configures access with tokenized URLs. This makes impossible to access content without a valid (unexpired) token.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `key` optional *string* &rarr;  Secure key generated on your side which will be used for the URL signing.

    - `type` optional *number* &rarr;  Type of the URL signing. Choose one of the values: 0 — to include the end user's IP address to secure token generation, 2 — to exclude the end user's IP address from the secure token generation.


  - `slice` optional *boolean* &rarr;  If enabled, CDN servers request and cache files larger than 10 MB in parts. Origins must support HTTP Range requests.

  - `stale` optional &rarr;  If enabled, CDN serves stale cached content in case of origin unavailability.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *string* &rarr;  The list of errors to which the option is applied.


  - `static_headers` optional &rarr;  Custom HTTP Headers that a CDN server adds to a response.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *map of* *string* &rarr;  A map of static headers in the format "header_name": "header_value".


  - `static_request_headers` optional &rarr;  Custom HTTP Headers for a CDN server to add to a request.
    - `enabled` optional *boolean* &rarr;  Controls the option state.

    - `value` optional *map of* *string* &rarr;  A map of static headers in the format "header_name": "header_value".



- `origin_protocol` optional *string* &rarr;  Protocol used by CDN servers to request content from an origin source. If protocol is not specified, HTTP is used to connect to an origin server.

- `region` optional *string* &rarr;  The region in which to obtain the CDN client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.

- `secondary_hostnames` optional *string* &rarr;  Additional delivery domains (CNAMEs) that will be used to deliver content via the CDN.

- `shielding` optional &rarr;  Use this attribute to configure origin shielding.
  - `enabled` optional *boolean* &rarr;  Defines whether origin shielding feature is enabled for the resource.

  - `pop_id` optional *number* &rarr;  ID of the origin shielding point of presence.


- `ssl_certificate` optional &rarr;  SSL certificate settings for content delivery over HTTPS protocol.
  - `id` optional *number* &rarr;  ID of the SSL certificate linked to the CDN resource. Must be configured when `type` is "own".

  - `type` optional *string* &rarr;  Type of the SSL certificate. Must be one of following: "not_used", "own", "lets_encrypt".

  - `status` read-only *string* &rarr;  Status of the SSL certificate.



## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *number* &rarr;  ID of the CDN resource.

- `preset_applied` *boolean* &rarr;  Protocol used by CDN servers to request content from an origin source.

- `status` *string* &rarr;  CDN resource status.

- `vp_enabled` *boolean* &rarr;  Defines whether the CDN resource is integrated with the Streaming Platform.



## Using a Let's Encrypt certificate

To issue a free [Let's Encrypt](https://letsencrypt.org/) certificate, specify "lets_encrypt" as 
the value for `ssl_certificate.type` argument. The certificate will be issued after the CDN 
resource is established, once the origin servers are available and DNS changes involving the 
CNAME records for personal domains have propagated.

~> **Note:** The option is only available for an active CDN resource, to achieve this, set the 
value of `active` argument to "true".

## Configuring ACLs

To enhance security, you can specify Access Control Lists (ACLs) options. All of the follow the 
same principles: when `policy_type` is "allow", it means that CDN server will allow access for all 
possible values of ACL subject except for those specified in `excepted_values` argument, and when 
`policy_type` is "deny", CDN will deny access with the same logic for excepted values.

### Example Configuration

For example, to protect content from unauthorized access from certain countries, you could use 
`country_acl` option:

```hcl
resource "vkcs_cdn_resource" "resource" {
  ...
  options = {
    country_acl = {
      policy_type     = "allow"
      excepted_values = ["GB", "DE"]
    }
  }
  ...
}

## Import

A CDN resource can be imported using the `id`, e.g.
```shell
terraform import vkcs_cdn_resource.resource <resource_id>
```
