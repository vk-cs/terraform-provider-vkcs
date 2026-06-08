package kubernetes_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesAddonsV2DataSource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAddonsV2DataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.vkcs_kubernetes_addons_v2.addons",
						"addons.#",
						regexp.MustCompile(`[1-9]\d*`),
					),

					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_addons_v2.addons", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_addons_v2.addons", "region"),
				),
			},
		},
	})
}

const testAccKubernetesAddonsV2DataSourceConfig = `
data "vkcs_kubernetes_addons_v2" "addons" {}
`
