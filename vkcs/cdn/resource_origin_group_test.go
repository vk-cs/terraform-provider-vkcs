package cdn_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccCDNOriginGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNOriginGroupResourceBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.basic", "name", "tfacc-origin-group-basic"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.basic", "origins.#", "1"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.basic", "origins.0.source", "example.com"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.basic", "origins.0.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.basic", "origins.0.backup", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.basic", "use_next", "false"),
				),
			},
			acctest.ImportStep("vkcs_cdn_origin_group.basic"),
		},
	})
}

func TestAccCDNOriginGroupResource_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNOriginGroupResourceFull),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.full", "name", "tfacc-origin-group-full"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.full", "origins.#", "2"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.full", "origins.0.source", "example1.com"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.full", "origins.0.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.full", "origins.0.backup", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.full", "origins.1.source", "example2.com"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.full", "origins.1.enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.full", "origins.1.backup", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.full", "use_next", "true"),
				),
			},
			acctest.ImportStep("vkcs_cdn_origin_group.full"),
		},
	})
}

func TestAccCDNOriginGroupResource_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNOriginGroupResourceUpdateOld),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "name", "tfacc-origin-group-update-old"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.#", "2"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.0.source", "example1.com"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.0.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.0.backup", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.1.source", "example2.com"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.1.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.1.backup", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "use_next", "true"),
				),
			},
			acctest.ImportStep("vkcs_cdn_origin_group.update"),
			{
				Config: acctest.AccTestRenderConfig(testAccCDNOriginGroupResourceUpdateNew),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "name", "tfacc-origin-group-update-new"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.#", "1"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.0.source", "example.com"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.0.enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "origins.0.backup", "false"),
					resource.TestCheckResourceAttr("vkcs_cdn_origin_group.update", "use_next", "false"),
				),
			},
			acctest.ImportStep("vkcs_cdn_origin_group.update"),
		},
	})
}

const testAccCDNOriginGroupResourceBasic = `
resource "vkcs_cdn_origin_group" "basic" {
  name = "tfacc-origin-group-basic"
  origins = [
    {
      source = "example.com"
    }
  ]
}
`

const testAccCDNOriginGroupResourceFull = `
resource "vkcs_cdn_origin_group" "full" {
  name = "tfacc-origin-group-full"
  origins = [
    {
      source  = "example1.com"
      enabled = true
      backup  = false
    },
    {
      source  = "example2.com"
      enabled = false
      backup  = true
    }
  ]
  use_next = true
}
`

const testAccCDNOriginGroupResourceUpdateOld = `
resource "vkcs_cdn_origin_group" "update" {
  name = "tfacc-origin-group-update-old"
  origins = [
    {
      source  = "example1.com"
    },
	{
      source  = "example2.com"
      backup  = true
    },
  ]
  use_next = true
}
`

const testAccCDNOriginGroupResourceUpdateNew = `
resource "vkcs_cdn_origin_group" "update" {
  name = "tfacc-origin-group-update-new"
  origins = [
    {
      source  = "example.com"
      enabled = true
    }
  ]
  use_next = false
}
`
