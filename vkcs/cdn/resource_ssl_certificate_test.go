package cdn_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccCDNSslCertificateResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNSslCertificateResourceBasic, map[string]string{"Certificate": sslCert, "PrivateKey": sslPrivateKey}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_ssl_certificate.basic", "name", "tfacc-ssl-certificate-basic"),
					resource.TestCheckResourceAttrSet("vkcs_cdn_ssl_certificate.basic", "certificate"),
					resource.TestCheckResourceAttrSet("vkcs_cdn_ssl_certificate.basic", "private_key"),
					resource.TestCheckResourceAttrSet("vkcs_cdn_ssl_certificate.basic", "issuer"),
					resource.TestCheckResourceAttrSet("vkcs_cdn_ssl_certificate.basic", "subject_cn"),
					resource.TestCheckResourceAttrSet("vkcs_cdn_ssl_certificate.basic", "validity_not_after"),
					resource.TestCheckResourceAttrSet("vkcs_cdn_ssl_certificate.basic", "validity_not_before"),
				),
			},
			acctest.ImportStep("vkcs_cdn_ssl_certificate.basic", "certificate", "private_key"),
		},
	})
}

func TestAccCDNSslCertificateResource_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNSslCertificateResourceUpdate, map[string]string{"Certificate": sslCert, "PrivateKey": sslPrivateKey,
					"Name": "tfacc-ssl-certificate-update-old"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_ssl_certificate.update", "name", "tfacc-ssl-certificate-update-old"),
				),
			},
			acctest.ImportStep("vkcs_cdn_ssl_certificate.update", "certificate", "private_key"),
			{
				Config: acctest.AccTestRenderConfig(testAccCDNSslCertificateResourceUpdate, map[string]string{"Certificate": sslCert, "PrivateKey": sslPrivateKey,
					"Name": "tfacc-ssl-certificate-update-new"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_ssl_certificate.update", "name", "tfacc-ssl-certificate-update-new"),
				),
			},
			acctest.ImportStep("vkcs_cdn_ssl_certificate.update", "certificate", "private_key"),
		},
	})
}

const sslCert = `-----BEGIN CERTIFICATE-----
MIIFKTCCAxGgAwIBAgIUJWBSairj3vo/zg2RBaSznNjbpO4wDQYJKoZIhvcNAQEL
BQAwIzELMAkGA1UEBhMCUlUxFDASBgNVBAMMC2V4YW1wbGUuY29tMCAXDTI0MDky
NDA5MTMzMloYDzIxMjQwODMxMDkxMzMyWjAjMQswCQYDVQQGEwJSVTEUMBIGA1UE
AwwLZXhhbXBsZS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCY
EYstwQlUWO7bgSpZIc50mSrmD3DbK7RWyumqTBZeCsQYapiCU6oFgzDJjDKOH9Dj
HtPo4kVfdKglDHz7TXXD0g25AHN91Bgbh88ssgJI/Q+ozlr1eVJ63vxidzCRZeyj
gFRX8d+mZhs2Q/019BamoOucRlCCYu/wppv9mbpefiMm5QrGmhWKVSjjUf6HVG9D
eW+BjAZUQNXuvEAgnsOHZAspPcmgKibZdXnGQhXJsNxuFvpdSC1/uVmj2ef97LoT
esGYyvYXI+XY0HhURocHMDN03uHdpwYUIt/192xsqLa4E8fxJA8p/mGyaj9TYKgJ
HPl+hFL5aFFBiRXTjWbzzZOC+kpSxvBTYNdoVo1y4X1fE2OhSz4WTmfbaS9HBWhY
No4xJxH0/9v1eVa4Xyslt0uSSKXhcxEuNzQLQ/7G3feDkii+BmudY2uv28Vp/8yC
5stRwMvtZdjphAW179TG508jJm7WULhvZ6c7OivpVW0YfNSHnPr2yT82ORDW6WxJ
IDHbWtYNyHGXhXb7pYa8uh9msb4CwbezQJ55d4+Y4m8AkYcjwzrGCyq4pbU3E3tW
MFtweGbaj0/87OLqpbpR+iocw2xxefq6G4dX1ywlaxT7MGIgfa4yq3ohrBX5+/iv
NgIRxo8x3bLeTjyH2BKgZJeLpaQKLmNbbKh1MecFwwIDAQABo1MwUTAdBgNVHQ4E
FgQUYIxAQrC6o+s7+FgTDuWnsxQM1rswHwYDVR0jBBgwFoAUYIxAQrC6o+s7+FgT
DuWnsxQM1rswDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAPNEE
nqE+/NWZElO+G0A8fPPmZJE3BtEnuOZIeFTc1JNKDilxnSGxBiXYV00Ej9T8NjPk
4Dgl24KKhnIHXrAO4Ps4UvqI/+PGgEFqyPoof8o5bUX+lgq6LI/ahW8JBR12XiBY
d+vK62ycoRsCtYHEvR1Eypf26psUjpCHYZYnQO5446xwoy25y9mzLekOgrjmezRJ
OfIH6Qt+vxoLaWYp/ScjTZ16mx3x/UsFPy5e5fmKp2Be/KmA0QdBiDB4XF7eZnnZ
nftJAWfxrbOxkWvJzI6aiUNHDkek3jKziJPbjtWSc6CjxkVwlEnFdc+iRrtez81E
HCc+xPZAY2mmGt4A8ApnKy+TZ5xPM7B/EszQH0wsIjcBOU1bnuy0El5MzQ2GxvRR
Nh9P7C94oHexDapgjrohbSLJO9X1vuFqVM6DndIzM+xm7ikrXl/3plQ0s+zhH4XP
91+0kW6D0JfZnomXdi6//UCrRarYhiYwhlSzSu9k3/6INijT1jEapP+kcgxDQ+X2
XXs/Wzsx8TK2dqYh1GdWlm4r98MUXEdjf3QA9BL56dpSfC18sfkuSAlL9uFJYpDT
ujP3aMd4Ki7Gn8CJGWHshHQfM25ksc5o2ElbuTCtE6q3dyWIg+S2fVA2ktNYwUUl
/0m5wBXHbi+Yq9pgKqxPOvzLo96J8SmLGxTLvGo=
-----END CERTIFICATE-----
`

const sslPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIJQQIBADANBgkqhkiG9w0BAQEFAASCCSswggknAgEAAoICAQCYEYstwQlUWO7b
gSpZIc50mSrmD3DbK7RWyumqTBZeCsQYapiCU6oFgzDJjDKOH9DjHtPo4kVfdKgl
DHz7TXXD0g25AHN91Bgbh88ssgJI/Q+ozlr1eVJ63vxidzCRZeyjgFRX8d+mZhs2
Q/019BamoOucRlCCYu/wppv9mbpefiMm5QrGmhWKVSjjUf6HVG9DeW+BjAZUQNXu
vEAgnsOHZAspPcmgKibZdXnGQhXJsNxuFvpdSC1/uVmj2ef97LoTesGYyvYXI+XY
0HhURocHMDN03uHdpwYUIt/192xsqLa4E8fxJA8p/mGyaj9TYKgJHPl+hFL5aFFB
iRXTjWbzzZOC+kpSxvBTYNdoVo1y4X1fE2OhSz4WTmfbaS9HBWhYNo4xJxH0/9v1
eVa4Xyslt0uSSKXhcxEuNzQLQ/7G3feDkii+BmudY2uv28Vp/8yC5stRwMvtZdjp
hAW179TG508jJm7WULhvZ6c7OivpVW0YfNSHnPr2yT82ORDW6WxJIDHbWtYNyHGX
hXb7pYa8uh9msb4CwbezQJ55d4+Y4m8AkYcjwzrGCyq4pbU3E3tWMFtweGbaj0/8
7OLqpbpR+iocw2xxefq6G4dX1ywlaxT7MGIgfa4yq3ohrBX5+/ivNgIRxo8x3bLe
TjyH2BKgZJeLpaQKLmNbbKh1MecFwwIDAQABAoICABxQ+/0nmmCh9Mxb93JEeMi+
cr4HNwkg0MJuo2cqJuoZEB3Jz59JC/pdzPpiyFEtvHxmU6hkZe2Z7+uCMU2sRVcS
6Ko/2sGd+mU5+0qD1SgZM07IKijWkBTAK/f74MfaVl+1uD7uE6rNDZkjvOVMj+E0
Sts9PqWg3bQOmjJ1az5IN6x47vI/Y+5v4B7AOGijwNosJSbW16Ddt9huJnTMi3VN
HETwM/jGkJhipyvTR3JYpBs93R38oDhN39LCc1AVwip5a85TUNLLRPQEEbwDrAYb
JCHJlP5sqRWbzt9i1MZA/lE1ocAV6lrz+uY5oZQZhgC4a/7yje7STXsqFy3fAUHw
Tat96KzDZCbtF1pPGn+3vrLKFoemdJSH3H1Vm31+j3CX/prV/4aEJds336dD9Pah
IY12s9oYrN+PVsaeLkAo6laWDherTy9dn7juz74X5h0WKPnTrY1pZ0mtiriVB4o9
CZQ+G74mTTwjZRbJ8qi/h8ZDG01UCAvj7FHJeg7xyl/qg5GuNl7YFqr8ZJqOGGEh
R3LhNsFRtrHnqj9lt3aqN9zLcSymVgb1Ufq5d0CgSVx3ANVyrDGkaMZSQT/AoJM0
9d6HEdBVOkKyclewqJfJ95rNAFI1tqF4I4Ma48DOPLnCSYSGgZSdVi5EQxM8yG5H
47PKMRpEbnn+sxmldSXhAoIBAQDK/msBW9IszWxq8nh7NT0RiLmSYZNDUP5XxcU6
2bVlTiGJGvq11A1/dRzyO2v0202AmtUZm/axZu1rh9fLaSjKR03I2FwbKxw01UgA
MW9J+mJN7zXZjV2FXMe34dpKR6R4unbj2/39INKJ/OHjByOOkpg3I/m0NmMtD0vu
3kiLSmaV37lzdkOl6RHB9c5xQ/OP5e9vDzFfGvEDfuoOBk23Or1SEKE43OUJBa3h
aswI2jQfawYS4hnMXOX3J+/Ysml0WJWESvvuTz6/uJfSiUwAs0TbMxFYaST1S+oC
/js9WQ3KG8orLfyE4ZP1Eq2s1QYXXXG2nwoOSFjvlZaA1ichAoIBAQC/xujCu4if
6C3kIYMGqCYobi3vqHm67v4L/OeuoFoQsMQ5pkvxNUAw8Pm9rrnfwAfqKQshkINw
+bNm9nO/TXqZrsXqRsqBOUuOiYOGxcArV+mGc15AjxqEyaR7GxOGlfiPHu4oPz5x
yrnVTRQqgSH9gplNDSv35+TeEBkodVGmt/b5jgek+qJs2J0GzH9pQGvNBPdTQ0k9
xoBsrSBLcmYpwZFEhOi7i4JvIpOFU7LjkcWvB8vQ+9i8mwgLUJoo5Tb/x9OQeYEP
uMiChVtriScgMlmGkrATDdYzBMAHAywnR7P8g2KCnQtfXelZYGeKM+DXs0upy7Yo
YETONoj3LWRjAoIBABWgQTol1CBdyj0ik99cbqMdk9eaeZvkVxR2x2pbo02lo0D+
FNOmQcHgcjMETZ4KdxlKZYWS7hc8RfL8x+qty2CxdAH/uuBSGXEvf7o1iguxlyOC
ZpRE0T/SAJ0AfMcJFuadxujDmS2Mf6GfxVjwe8NGrtzBAmtGmA5G4OoT2FqulHtH
GHTKlq5oRDILw/ChMqOT9Yw5bCMbta2PqdPQrBrnMSA7EVIDhosNhdbMD+ypgoAO
YNlGKUVyaDWKlazaZQ22Gke7zVc4LhEy00nkwqoYby+DI0ft+7f+XHHxL9J7WFK9
3y32ej6V4bNsSABvuXRnyiOQkfuvjXoIz73uEsECggEACodTFA3TrGPE0Td9yAnH
PoT+BKBNPpQMYoAaAB5Rk4UA4OFeXfm8cnNoYp+LGNukE5j5QXh7nuI2lTqGlEQe
rP1JZFlKmNmaalLmY6nLqRWEfpsq24R/wjaHzzJnWgY1xTW/gXonZXvpw+odZ0/7
m71lyTMl7NBQYlij6PK37t0+s+i2Rrpz3GHwDQWBITgmMvVI1stk4/1X45+FnF7F
eRllbkuVs2YvXQaa7sDvm0rPfZKCABEzjvc789MTA5fB8zz3QoFJMqHEcFX99ONs
wHnDLH41KHakAd0K5deovudS3FQiPmV80FmJjByc0puShoUTbFkAwSq33FJmJpvZ
XQKCAQAi0oCeDLIqaD+m1ykDmg/2QqCPXcDLqlb7mSfm7xxM07fhWzhMtKA4kSnm
aHt/6CK9EEO7OFaKrJAC4iK97I4rilB7afjfluPDnqABnyr638AZngQGJYQP6Mcy
rCanXZVIT+oymzUfbIKz+J5DK8pSOh7IT8xBjH1pMBHk6Hs3lq5M1U0ItOaj30Fm
N48+mmnBaEP7q+z8SOrYc2s7rrBpm+YkuIZwhhZxHdwLSlU0uC7LUpIr3MWyfcop
cnR4seWhJQIu0Nk4dsStQbJpjRBdDkqKFupwpAS74Gx3HgHcDwoJwMJv2rH1AoEu
WcDW/RNZpfcNMb0+ZdPgl7nCrCY7
-----END PRIVATE KEY-----
`

const testAccCDNSslCertificateResourceBasic = `
resource "vkcs_cdn_ssl_certificate" "basic" {
  name        = "tfacc-ssl-certificate-basic"
  certificate = <<EOT
{{ .Certificate }}
EOT

  private_key = <<EOT
{{ .PrivateKey }}
EOT
}
`

const testAccCDNSslCertificateResourceUpdate = `
resource "vkcs_cdn_ssl_certificate" "update" {
  name        = "{{ .Name }}"
  certificate = <<EOT
{{ .Certificate }}
EOT

  private_key = <<EOT
{{ .PrivateKey }}
EOT
}
`
