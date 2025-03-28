package cdn_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccCDNResourceResource_basic(t *testing.T) {
	nameSuffix := acctest.GenerateNameSuffix()
	baseName := "tfacc-resource-base-" + nameSuffix
	cname := "tfacc-basic-" + nameSuffix + ".vk.com"

	ogBaseConfig := acctest.AccTestRenderConfig(testAccCDNResourceResourceOriginGroupBase, map[string]string{"Name": baseName})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNResourceResourceBasic, map[string]string{"Cname": cname, "TestAccCDNResourceResourceOriginGroupBase": ogBaseConfig}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_resource.basic", "cname", cname),
					resource.TestCheckResourceAttrPair("vkcs_cdn_resource.basic", "origin_group", "vkcs_cdn_origin_group.base", "id"),
				),
			},
			acctest.ImportStep("vkcs_cdn_resource.basic"),
		},
	})
}

func TestAccCDNResourceResource_full(t *testing.T) {
	nameSuffix := acctest.GenerateNameSuffix()
	baseName := "tfacc-resource-base-" + nameSuffix
	cname := "tfacc-full-" + nameSuffix + ".vk.com"

	ogBaseConfig := acctest.AccTestRenderConfig(testAccCDNResourceResourceOriginGroupBase, map[string]string{"Name": baseName})
	sslCertificateBaseConfig := acctest.AccTestRenderConfig(testAccCDNResourceResourceSslCertificateBase, map[string]string{"Name": baseName, "Certificate": sslCert, "PrivateKey": sslPrivateKey})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNResourceResourceFull, map[string]string{"Cname": cname, "TestAccCDNResourceResourceOriginGroupBase": ogBaseConfig, "TestAccCDNResourceResourceShieldingPopBase": testAccCDNResourceResourceShieldingPopBase, "TestAccCDNResourceResourceSslCertificateBase": sslCertificateBaseConfig}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "cname", cname),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "active", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.allowed_http_methods.value.#", "3"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.allowed_http_methods.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.brotli_compression.value.#", "3"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.brotli_compression.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.browser_cache_settings.value", "3600s"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.browser_cache_settings.enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.cors.value.#", "2"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.cors.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.edge_cache_settings.value", "10m"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.edge_cache_settings.custom_values.%", "2"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.edge_cache_settings.custom_values.200", "60s"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.edge_cache_settings.custom_values.404", "30m"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.fetch_compressed", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.force_return.code", "301"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.force_return.body", "https://vk.com/redirect"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.forward_host_header", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.gzip_on", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.ignore_cookie", "false"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.full", "options.query_params_blacklist.value", []string{"some", "query"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.country_acl.policy_type", "allow"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.full", "options.country_acl.excepted_values", []string{"GB", "DE"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.referrer_acl.policy_type", "deny"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.full", "options.referrer_acl.excepted_values", []string{"example.com", "*.example.net"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.ip_address_acl.policy_type", "allow"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.full", "options.ip_address_acl.excepted_values", []string{"192.168.1.100/32"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.slice", "false"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.full", "options.stale.value", []string{"http_403", "http_404"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.stale.enabled", "true"),
					acctest.TestCheckResourceMapAttr("vkcs_cdn_resource.full", "options.static_request_headers.value", map[string]string{"Header-One": "Value 1", "Header-Two": "Value 2"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.static_request_headers.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.secure_key.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.secure_key.key", "mysupersecretkey"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.secure_key.type", "0"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "options.static_response_headers.enabled", "true"),
					acctest.TestCheckResourceAttrDeepEqual("vkcs_cdn_resource.full", "options.static_response_headers.value", []map[string]any{{"name": "First-Header", "value": []string{"Header1"}, "always": true}, {"name": "Second-Header", "value": []string{"Header2"}, "always": false}}),
					resource.TestCheckResourceAttrPair("vkcs_cdn_resource.full", "origin_group", "vkcs_cdn_origin_group.base", "id"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "origin_protocol", "MATCH"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.full", "secondary_hostnames", []string{"cdn1.vk.com", "cdn2.vk.com"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "ssl_certificate.type", "own"),
					resource.TestCheckResourceAttrPair("vkcs_cdn_resource.full", "ssl_certificate.id", "vkcs_cdn_ssl_certificate.base", "id"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.full", "shielding.enabled", "true"),
					resource.TestCheckResourceAttrPair("vkcs_cdn_resource.full", "shielding.pop_id", "data.vkcs_cdn_shielding_pop.base", "id"),
				),
			},
			acctest.ImportStep("vkcs_cdn_resource.full"),
		},
	})
}

func TestAccCDNResourceResource_update(t *testing.T) {
	nameSuffix := acctest.GenerateNameSuffix()
	ogOldName := "tfacc-resource-update-old" + nameSuffix
	cname := "tfacc-update-" + nameSuffix + ".vk.com"
	ogNewName := "tfacc-resource-update-new-" + nameSuffix

	ogOldConfig := acctest.AccTestRenderConfig(testAccCDNResourceResourceOriginGroupBase, map[string]string{"Name": ogOldName})
	ogNewConfig := acctest.AccTestRenderConfig(testAccCDNResourceResourceOriginGroupBase, map[string]string{"Name": ogNewName})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNResourceResourceUpdateOld, map[string]string{"Cname": cname, "TestAccCDNResourceResourceOriginGroupBase": ogOldConfig, "TestAccCDNResourceResourceShieldingPopBase": testAccCDNResourceResourceShieldingPopBase}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "cname", cname),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "active", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.allowed_http_methods.value.#", "4"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.allowed_http_methods.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.brotli_compression.value.#", "3"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.brotli_compression.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.browser_cache_settings.value", "3600s"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.browser_cache_settings.enabled", "false"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.cors.value", []string{"app1.vk.com", "app2.vk.com"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.edge_cache_settings.value", "10m"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.edge_cache_settings.custom_values.%", "2"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.edge_cache_settings.custom_values.200", "60s"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.edge_cache_settings.custom_values.404", "30m"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.fetch_compressed", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.force_return.code", "301"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.force_return.body", "https://vk.com/redirect/old"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.forward_host_header", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.gzip_on", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.ignore_cookie", "false"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.query_params_blacklist.value", []string{"some", "query"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.country_acl.policy_type", "allow"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.country_acl.excepted_values", []string{"GB", "DE"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.referrer_acl.policy_type", "deny"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.referrer_acl.excepted_values", []string{"example1.com", "*.example2.net"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.ip_address_acl.policy_type", "allow"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.ip_address_acl.excepted_values", []string{"192.168.1.100/32"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.slice", "false"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.stale.value", []string{"http_403", "http_404"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.stale.enabled", "true"),
					acctest.TestCheckResourceMapAttr("vkcs_cdn_resource.update", "options.static_request_headers.value", map[string]string{"Header-One": "Old Value 1", "Header-Two": "Old Value 2"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.static_request_headers.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.secure_key.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.secure_key.key", "mysupersecretkey"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.secure_key.type", "0"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.static_response_headers.enabled", "true"),
					acctest.TestCheckResourceAttrDeepEqual("vkcs_cdn_resource.update", "options.static_response_headers.value", []map[string]any{{"name": "First-Header", "value": []string{"Header1"}, "always": true}, {"name": "Second-Header", "value": []string{"Header2"}, "always": false}}),
					resource.TestCheckResourceAttrPair("vkcs_cdn_resource.update", "origin_group", "vkcs_cdn_origin_group.base", "id"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "origin_protocol", "HTTP"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "shielding.enabled", "true"),
					resource.TestCheckResourceAttrPair("vkcs_cdn_resource.update", "shielding.pop_id", "data.vkcs_cdn_shielding_pop.base", "id"),
				),
			},
			acctest.ImportStep("vkcs_cdn_resource.update"),
			{
				Config: acctest.AccTestRenderConfig(testAccCDNResourceResourceUpdateNew, map[string]string{"Cname": cname, "TestAccCDNResourceResourceOriginGroupBase": ogNewConfig}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("vkcs_cdn_resource.update", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "cname", cname),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "active", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.allowed_http_methods.value.#", "3"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.allowed_http_methods.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.brotli_compression.value.#", "3"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.brotli_compression.enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.browser_cache_settings.value", "5m"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.browser_cache_settings.enabled", "true"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.cors.value", []string{"app3.vk.com", "app1.vk.com"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.cors.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.edge_cache_settings.default", "10m"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.edge_cache_settings.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.fetch_compressed", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.force_return.code", "301"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.force_return.body", "https://vk.com/redirect/old"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.force_return.enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.forward_host_header", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.gzip_on", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.host_header.value", "host.com"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.host_header.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.ignore_cookie", "true"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.query_params_blacklist.value", []string{"some", "query"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.query_params_blacklist.enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.country_acl.policy_type", "deny"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.country_acl.excepted_values", []string{"BI", "JE"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.country_acl.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.referrer_acl.policy_type", "deny"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.referrer_acl.excepted_values", []string{"example1.com", "*.example2.net"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.referrer_acl.enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.ip_address_acl.policy_type", "allow"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.ip_address_acl.excepted_values", []string{"192.168.1.100/32", "192.168.1.200/32"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.slice", "false"),
					acctest.TestCheckResourceListAttr("vkcs_cdn_resource.update", "options.stale.value", []string{"http_403", "http_404"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.stale.enabled", "false"),
					acctest.TestCheckResourceMapAttr("vkcs_cdn_resource.update", "options.static_request_headers.value", map[string]string{"Header-Two": "Value 2", "Header-Three": "Value 3"}),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.static_request_headers.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.secure_key.enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.secure_key.key", "mysimplekey"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.secure_key.type", "2"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "options.static_response_headers.enabled", "true"),
					acctest.TestCheckResourceAttrDeepEqual("vkcs_cdn_resource.update", "options.static_response_headers.value", []map[string]any{{"name": "Second-Header", "value": []string{"Header1", "Header2"}, "always": true}}),
					resource.TestCheckResourceAttrPair("vkcs_cdn_resource.update", "origin_group", "vkcs_cdn_origin_group.base", "id"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "origin_protocol", "HTTPS"),
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "shielding.enabled", "false"),
				),
			},
			acctest.ImportStep("vkcs_cdn_resource.update"),
		},
	})
}

func TestAccCDNResourceResource_editSSLCertificate(t *testing.T) {
	nameSuffix := acctest.GenerateNameSuffix()
	baseName := "tfacc-resource-base-" + nameSuffix
	cname := "tfacc-update-ssl-" + nameSuffix + ".vk.com"
	oldCertName := "tfacc-resource-update-ssl-old-" + nameSuffix
	newCertName := "tfacc-resource-update-ssl-new-" + nameSuffix

	ogBaseConfig := acctest.AccTestRenderConfig(testAccCDNResourceResourceOriginGroupBase, map[string]string{"Name": baseName})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNResourceResourceUpdateSSLCertificateNotUsed, map[string]string{"Cname": cname, "TestAccCDNResourceResourceOriginGroupBase": ogBaseConfig}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "ssl_certificate.type", "not_used"),
				),
			},
			acctest.ImportStep("vkcs_cdn_resource.update"),
			{
				Config: acctest.AccTestRenderConfig(testAccCDNResourceResourceUpdateSSLCertificateOwnOld, map[string]string{"Cname": cname, "OldCertName": oldCertName, "Certificate": sslCert, "PrivateKey": sslPrivateKey, "TestAccCDNResourceResourceOriginGroupBase": ogBaseConfig}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("vkcs_cdn_resource.update", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "ssl_certificate.type", "own"),
					resource.TestCheckResourceAttrPair("vkcs_cdn_resource.update", "ssl_certificate.id", "vkcs_cdn_ssl_certificate.old", "id"),
				),
			},
			acctest.ImportStep("vkcs_cdn_resource.update"),
			{
				Config: acctest.AccTestRenderConfig(testAccCDNResourceResourceUpdateSSLCertificateOwnNew, map[string]string{"Cname": cname, "OldCertName": oldCertName, "NewCertName": newCertName, "Certificate": sslCert, "PrivateKey": sslPrivateKey, "TestAccCDNResourceResourceOriginGroupBase": ogBaseConfig}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("vkcs_cdn_resource.update", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "ssl_certificate.type", "own"),
					resource.TestCheckResourceAttrPair("vkcs_cdn_resource.update", "ssl_certificate.id", "vkcs_cdn_ssl_certificate.new", "id"),
				),
			},
			acctest.ImportStep("vkcs_cdn_resource.update"),
			{
				Config: acctest.AccTestRenderConfig(testAccCDNResourceResourceUpdateSSLCertificateLetsEncrypt, map[string]string{"Cname": cname, "NewCertName": newCertName, "Certificate": sslCert, "PrivateKey": sslPrivateKey, "TestAccCDNResourceResourceOriginGroupBase": ogBaseConfig}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("vkcs_cdn_resource.update", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_resource.update", "ssl_certificate.type", "lets_encrypt"),
				),
			},
			acctest.ImportStep("vkcs_cdn_resource.update"),
		},
	})
}

const testAccCDNResourceResourceOriginGroupBase = `
resource "vkcs_cdn_origin_group" "base" {
  name  = "{{ .Name }}"
  origins = [
	{
      source = "origin.vk.com"
    }
  ]
}
`

const testAccCDNResourceResourceShieldingPopBase = `
data "vkcs_cdn_shielding_pop" "base" {
  city = "Moscow-Megafon"
}
`

const testAccCDNResourceResourceSslCertificateBase = `
resource "vkcs_cdn_ssl_certificate" "base" {
  name        = "{{ .Name }}"
  certificate = <<EOT
{{ .Certificate }}
EOT

  private_key = <<EOT
{{ .PrivateKey }}
EOT
}
`

const testAccCDNResourceResourceBasic = `
{{ .TestAccCDNResourceResourceOriginGroupBase }}

resource "vkcs_cdn_resource" "basic" {
  cname        = "{{ .Cname }}"
  origin_group = vkcs_cdn_origin_group.base.id
}
`

const testAccCDNResourceResourceFull = `
{{ .TestAccCDNResourceResourceOriginGroupBase }}
{{ .TestAccCDNResourceResourceShieldingPopBase }}
{{ .TestAccCDNResourceResourceSslCertificateBase }}

resource "vkcs_cdn_resource" "full" {
  cname      = "{{ .Cname }}"
  active     = true
  options = {
    allowed_http_methods = {
      value   = ["GET", "HEAD", "OPTIONS"]
      enabled = true
    }
    brotli_compression = {
      value = ["text/html", "text/css", "application/javascript"]
    }
    browser_cache_settings = {
      value   = "3600s"
      enabled = false
    }
    cors = {
      value = ["app1.vk.com", "app2.vk.com"]
    }
    edge_cache_settings = {
      value = "10m"
      custom_values = {
        "200" : "60s",
        "404" : "30m"
      }
    }
    fetch_compressed = false
    force_return = {
      code = 301
      body = "https://vk.com/redirect"
    }
    forward_host_header = true
    gzip_on             = true
    ignore_cookie       = false
    query_params_blacklist = {
      value = ["some", "query"]
    }
    country_acl = {
      policy_type     = "allow"
      excepted_values = ["GB", "DE"]
    }
    referrer_acl = {
      policy_type     = "deny"
      excepted_values = ["example.com", "*.example.net"]
    }
    ip_address_acl = {
      policy_type     = "allow"
      excepted_values = ["192.168.1.100/32"]
    }
    slice = false
    stale = {
      value = ["http_403", "http_404"]
    }
    static_request_headers = {
      value = {
        "Header-One" : "Value 1",
        "Header-Two" : "Value 2"
      }
    }
	secure_key = {
	  enabled = true
      key	  = "mysupersecretkey"
      type	  = 0
	}
    static_response_headers = {
      enabled = true
      value = [
        {
          name = "First-Header"
          value = ["Header1"]
          always = true
        },
        {
          name = "Second-Header"
          value = ["Header2"]
          always = false
        }
      ]
    }
  }
  origin_group        = vkcs_cdn_origin_group.base.id
  origin_protocol     = "MATCH"
  secondary_hostnames = ["cdn1.vk.com", "cdn2.vk.com"]
  ssl_certificate = {
    type = "own"
    id   = vkcs_cdn_ssl_certificate.base.id
  }
  shielding = {
    pop_id  = data.vkcs_cdn_shielding_pop.base.id
    enabled = true
  }
}
`

const testAccCDNResourceResourceUpdateOld = `
{{ .TestAccCDNResourceResourceOriginGroupBase }}
{{ .TestAccCDNResourceResourceShieldingPopBase }}

resource "vkcs_cdn_resource" "update" {
  cname  = "{{ .Cname }}"
  active = false
  options = {
    allowed_http_methods = {
      value   = ["GET", "OPTIONS", "PUT", "PATCH"]
      enabled = true
    }
    brotli_compression = {
      value = ["text/html", "text/css", "application/javascript"]
    }
    browser_cache_settings = {
      value   = "3600s"
      enabled = false
    }
    cors = {
      value = ["app1.vk.com", "app2.vk.com"]
    }
    edge_cache_settings = {
      value = "10m"
      custom_values = {
        "200" : "60s",
        "404" : "30m"
      }
    }
    fetch_compressed = false
    force_return = {
      code = 301
      body = "https://vk.com/redirect/old"
    }
    forward_host_header = true
    gzip_on             = true
    ignore_cookie       = false
    query_params_blacklist = {
      value = ["some", "query"]
    }
    country_acl = {
      policy_type     = "allow"
      excepted_values = ["GB", "DE"]
    }
    referrer_acl = {
      policy_type     = "deny"
      excepted_values = ["example1.com", "*.example2.net"]
    }
    ip_address_acl = {
      policy_type     = "allow"
      excepted_values = ["192.168.1.100/32"]
    }
    slice = false
    stale = {
      value = ["http_403", "http_404"]
    }
    static_request_headers = {
      value = {
        "Header-One" : "Old Value 1",
        "Header-Two" : "Old Value 2"
      }
    }
	secure_key = {
	  enabled = true
      key	  = "mysupersecretkey"
      type	  = 0
	}
    static_response_headers = {
      enabled = true
      value = [
        {
          name = "First-Header"
          value = ["Header1"]
          always = true
        },
        {
          name = "Second-Header"
          value = ["Header2"]
          always = false
        }
      ]
    }
  }
  origin_group    = vkcs_cdn_origin_group.base.id
  origin_protocol = "HTTP"
  shielding = {
    pop_id  = data.vkcs_cdn_shielding_pop.base.id
    enabled = true
  }
}
`

const testAccCDNResourceResourceUpdateNew = `
{{ .TestAccCDNResourceResourceOriginGroupBase }}

resource "vkcs_cdn_resource" "update" {
  cname        = "{{ .Cname }}"
  active       = true
  options = {
    allowed_http_methods = {
      value   = ["GET", "HEAD", "OPTIONS"]
      enabled = true
    }
    brotli_compression = {
      enabled = false
    }
    browser_cache_settings = {
      value   = "5m"
      enabled = true
    }
    cors = {
      value = ["app3.vk.com", "app1.vk.com"]
    }
    edge_cache_settings = {
      default = "10m"
    }
    fetch_compressed = true
    force_return = {
      enabled = false
    }
    forward_host_header = false
    gzip_on             = false
	  host_header         = {
	    value = "host.com"
	  }
    ignore_cookie       = true
    query_params_blacklist = {
      enabled = false
    }
    query_params_whitelist = {
      value = ["some", "other", "query"]
    }
    country_acl = {
      policy_type     = "deny"
      excepted_values = ["BI", "JE"]
    }
    referrer_acl = {
      enabled = false
    }
    ip_address_acl = {
      policy_type     = "allow"
      excepted_values = ["192.168.1.100/32", "192.168.1.200/32"]
    }
    stale = {
	    enabled = false
    }
    static_request_headers = {
	    value = {
        "Header-Two" : "Value 2",
		    "Header-Three" : "Value 3" 
	    }
    }
	secure_key = {
	  enabled = false,
      key     = "mysimplekey"
      type    = 2
	}
    static_response_headers = {
      enabled = true
      value = [
        {
          name = "Second-Header"
          value = [
            "Header1",
            "Header2"
          ]
          always = true
        }
      ]
    }
  }
  origin_group    = vkcs_cdn_origin_group.base.id
  origin_protocol = "HTTPS"
  shielding = {
    enabled = false 
  }
}
`

const testAccCDNResourceResourceUpdateSSLCertificateNotUsed = `
{{ .TestAccCDNResourceResourceOriginGroupBase }}

resource "vkcs_cdn_resource" "update" {
  active       = true
  cname        = "{{ .Cname }}"
  origin_group = vkcs_cdn_origin_group.base.id
  ssl_certificate = {
    type = "not_used"
  }
}
`

const testAccCDNResourceResourceUpdateSSLCertificateOwnOld = `
{{ .TestAccCDNResourceResourceOriginGroupBase }}

resource "vkcs_cdn_ssl_certificate" "old" {
  name        = "{{ .OldCertName }}"
  certificate = <<EOT
{{ .Certificate }}
EOT

  private_key = <<EOT
{{ .PrivateKey }}
EOT
}

resource "vkcs_cdn_resource" "update" {
  active       = true
  cname        = "{{ .Cname }}"
  origin_group = vkcs_cdn_origin_group.base.id
  ssl_certificate = {
    type = "own"
    id   = vkcs_cdn_ssl_certificate.old.id
  }
}
`

const testAccCDNResourceResourceUpdateSSLCertificateOwnNew = `
{{ .TestAccCDNResourceResourceOriginGroupBase }}

resource "vkcs_cdn_ssl_certificate" "old" {
  name        = "{{ .OldCertName }}"
  certificate = <<EOT
{{ .Certificate }}
EOT

  private_key = <<EOT
{{ .PrivateKey }}
EOT

  # Controls the order of delete operations, so the certificate
  # will not be deleted before the resource if update failed
  depends_on = [vkcs_cdn_resource.update]
}

resource "vkcs_cdn_ssl_certificate" "new" {
  name        = "{{ .NewCertName }}"
  certificate = <<EOT
{{ .Certificate }}
EOT

  private_key = <<EOT
{{ .PrivateKey }}
EOT
}

resource "vkcs_cdn_resource" "update" {
  active       = true
  cname        = "{{ .Cname }}"
  origin_group = vkcs_cdn_origin_group.base.id
  ssl_certificate = {
    type = "own"
    id   = vkcs_cdn_ssl_certificate.new.id
  }
}
`

const testAccCDNResourceResourceUpdateSSLCertificateLetsEncrypt = `
{{ .TestAccCDNResourceResourceOriginGroupBase }}

resource "vkcs_cdn_ssl_certificate" "new" {
  name        = "{{ .NewCertName }}"
  certificate = <<EOT
{{ .Certificate }}
EOT

  private_key = <<EOT
{{ .PrivateKey }}
EOT

  # Controls the order of delete operations, so the certificate
  # will not be deleted before the resource if update failed
  depends_on = [vkcs_cdn_resource.update]
}

resource "vkcs_cdn_resource" "update" {
  active       = true
  cname        = "{{ .Cname }}"
  origin_group = vkcs_cdn_origin_group.base.id
  ssl_certificate = {
    type = "lets_encrypt"
  }
}
`
