package kubernetes_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesVersionsV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesVersionsV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.vkcs_kubernetes_versions_v2.versions",
						"k8s_versions.#",
						regexp.MustCompile(`[1-9]\d*`),
					),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_kubernetes_versions_v2.versions",
						"region",
					),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_kubernetes_versions_v2.versions",
						"id",
					),
					resource.TestCheckResourceAttr(
						"data.vkcs_kubernetes_versions_v2.versions",
						"id",
						"kubernetes_versions",
					),
				),
			},
		},
	})
}

const testAccKubernetesVersionsV2Config = `
data "vkcs_kubernetes_versions_v2" "versions" {}
`
