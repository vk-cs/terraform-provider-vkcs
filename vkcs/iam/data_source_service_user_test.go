package iam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccIAMServiceUserDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccIAMServiceUserDataSourceBasic, map[string]string{"TestAccIAMServiceUserDataSourceBase": testAccIAMServiceUserDataSourceBase}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.vkcs_iam_service_user.basic", "name", "vkcs_iam_service_user.base", "name"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_service_user.basic", "role_names", "vkcs_iam_service_user.base", "role_names"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_service_user.basic", "description", "vkcs_iam_service_user.base", "description"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_service_user.basic", "creator_name", "vkcs_iam_service_user.base", "creator_name"),
				),
			},
		},
	})
}

func TestAccIAMServiceUserDataSource_queryByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccIAMServiceUserDataSourceQueryByName, map[string]string{"TestAccIAMServiceUserDataSourceQueryByNameBase": testAccIAMServiceUserDataSourceQueryByNameBase}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.vkcs_iam_service_user.basic", "id", "vkcs_iam_service_user.k8s_admin", "id"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_service_user.basic", "name", "vkcs_iam_service_user.k8s_admin", "name"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_service_user.basic", "role_names", "vkcs_iam_service_user.k8s_admin", "role_names"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_service_user.basic", "description", "vkcs_iam_service_user.k8s_admin", "description"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_service_user.basic", "creator_name", "vkcs_iam_service_user.k8s_admin", "creator_name"),
				),
			},
		},
	})
}

const testAccIAMServiceUserDataSourceBase = `
resource "vkcs_iam_service_user" "base" {
  name       = "tfacc-service-user-base"
  role_names = ["mcs_k8s_admin"]
}
`

const testAccIAMServiceUserDataSourceBasic = `
{{ .TestAccIAMServiceUserDataSourceBase }}

data "vkcs_iam_service_user" "basic" {
  id = vkcs_iam_service_user.base.id
}
`

const testAccIAMServiceUserDataSourceQueryByNameBase = `
resource "vkcs_iam_service_user" "k8s_admin" {
  name = "tfacc-k8s-admin"
  role_names = ["mcs_k8s_admin"]
}

resource "vkcs_iam_service_user" "k8s_viewer" {
  name = "tfacc-k8s-viewer"
  role_names = ["mcs_k8s_viewer"]
}
`

const testAccIAMServiceUserDataSourceQueryByName = `
{{ .TestAccIAMServiceUserDataSourceQueryByNameBase }}

data "vkcs_iam_service_user" "basic" {
  name = vkcs_iam_service_user.k8s_admin.name
}
`
