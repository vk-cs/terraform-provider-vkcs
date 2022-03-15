package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKeyManagerContainerDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_container.container_1", "id",
						"vkcs_keymanager_container.container_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_container.container_1", "secret_refs",
						"vkcs_keymanager_container.container_1", "secret_refs"),
					resource.TestCheckResourceAttr(
						"data.vkcs_keymanager_container.container_1", "secret_refs.#", "3"),
				),
			},
		},
	})
}

func TestAccKeyManagerContainerDataSource_acls(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerDataSourceAcls,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_container.container_1", "id",
						"vkcs_keymanager_container.container_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_container.container_1", "secret_refs",
						"vkcs_keymanager_container.container_1", "secret_refs"),
					resource.TestCheckResourceAttr(
						"data.vkcs_keymanager_container.container_1", "secret_refs.#", "3"),
					resource.TestCheckResourceAttr("data.vkcs_keymanager_container.container_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("data.vkcs_keymanager_container.container_1", "acl.0.read.0.users.#", "2"),
				),
			},
		},
	})
}

const testAccKeyManagerContainerDataSourceBasic = `
resource "vkcs_keymanager_secret" "certificate_1" {
  name                 = "certificate"
  payload              = "certificate"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "private_key_1" {
  name                 = "private_key"
  payload              = "private_key"
  secret_type          = "private"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "intermediate_1" {
  name                 = "intermediate"
  payload              = "intermediate"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_container" "container_1" {
  name = "generic"
  type = "generic"

  secret_refs {
    name       = "certificate"
    secret_ref = "${vkcs_keymanager_secret.certificate_1.secret_ref}"
  }

  secret_refs {
    name       = "private_key"
    secret_ref = "${vkcs_keymanager_secret.private_key_1.secret_ref}"
  }

  secret_refs {
    name       = "intermediate"
    secret_ref = "${vkcs_keymanager_secret.intermediate_1.secret_ref}"
  }
}

data "vkcs_keymanager_container" "container_1" {
  name = "${vkcs_keymanager_container.container_1.name}"
}
`

const testAccKeyManagerContainerDataSourceAcls = `
resource "vkcs_keymanager_secret" "certificate_1" {
  name                 = "certificate"
  payload              = "certificate"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "private_key_1" {
  name                 = "private_key"
  payload              = "private_key"
  secret_type          = "private"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "intermediate_1" {
  name                 = "intermediate"
  payload              = "intermediate"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_container" "container_1" {
  name = "generic"
  type = "generic"

  secret_refs {
    name       = "certificate"
    secret_ref = "${vkcs_keymanager_secret.certificate_1.secret_ref}"
  }

  secret_refs {
    name       = "private_key"
    secret_ref = "${vkcs_keymanager_secret.private_key_1.secret_ref}"
  }

  secret_refs {
    name       = "intermediate"
    secret_ref = "${vkcs_keymanager_secret.intermediate_1.secret_ref}"
  }

  acl {
    read {
      project_access = false
      users = [
        "96b3ebddf275996285eae440e71227ba47c651be18391b0f2ebf1032ebae5dca",
        "619e2ad074321cf246b03a89e95afee95fb26bb0b2d1fc7ba3bd30fcca25588a",
      ]
    }
  }
}

data "vkcs_keymanager_container" "container_1" {
  name = "${vkcs_keymanager_container.container_1.name}"
}
`
