package kubernetes_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesVolumeTypesV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesVolumeTypesV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.vkcs_kubernetes_volume_types_v2.types",
						"volume_types.#",
						regexp.MustCompile(`[1-9]\d*`),
					),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_kubernetes_volume_types_v2.types",
						"region",
					),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_kubernetes_volume_types_v2.types",
						"id",
					),
					resource.TestCheckResourceAttr(
						"data.vkcs_kubernetes_volume_types_v2.types",
						"id",
						"volume_types",
					),
				),
			},
		},
	})
}

const testAccKubernetesVolumeTypesV2Config = `data "vkcs_kubernetes_volume_types_v2" "types" {}`
