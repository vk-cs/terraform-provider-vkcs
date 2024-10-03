package cdn_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccCDNSslCertificateDataSource_basic(t *testing.T) {
	baseConfig := acctest.AccTestRenderConfig(testAccCDNSslCertificateDataSourceBase, map[string]string{"Certificate": sslCert, "PrivateKey": sslPrivateKey})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNSslCertificateDataSourceBasic, map[string]string{"TestAccCDNSslCertificateDataSourceBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_cdn_ssl_certificate.basic", "name", "tfacc-ssl-certificate-base"),
					resource.TestCheckResourceAttrPair("data.vkcs_cdn_ssl_certificate.basic", "issuer", "vkcs_cdn_ssl_certificate.base", "issuer"),
					resource.TestCheckResourceAttrPair("data.vkcs_cdn_ssl_certificate.basic", "subject_cn", "vkcs_cdn_ssl_certificate.base", "subject_cn"),
					resource.TestCheckResourceAttrPair("data.vkcs_cdn_ssl_certificate.basic", "validity_not_after", "vkcs_cdn_ssl_certificate.base", "validity_not_after"),
					resource.TestCheckResourceAttrPair("data.vkcs_cdn_ssl_certificate.basic", "validity_not_before", "vkcs_cdn_ssl_certificate.base", "validity_not_before"),
				),
			},
		},
	})
}

const testAccCDNSslCertificateDataSourceBase = `
resource "vkcs_cdn_ssl_certificate" "base" {
  name        = "tfacc-ssl-certificate-base"
  certificate = <<EOT
{{ .Certificate }}
EOT

  private_key = <<EOT
{{ .PrivateKey }}
EOT
}
`

const testAccCDNSslCertificateDataSourceBasic = `
{{ .TestAccCDNSslCertificateDataSourceBase }}

data "vkcs_cdn_ssl_certificate" "basic" {
  name = vkcs_cdn_ssl_certificate.base.name
}
`
