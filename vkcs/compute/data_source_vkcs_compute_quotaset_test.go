package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccComputeQuotasetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeQuotasetDataSourceSource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeQuotasetDataSourceID("data.vkcs_compute_quotaset.source"),
					resource.TestCheckResourceAttrSet("data.vkcs_compute_quotaset.source", "key_pairs"),
					resource.TestCheckResourceAttrSet("data.vkcs_compute_quotaset.source", "metadata_items"),
					resource.TestCheckResourceAttrSet("data.vkcs_compute_quotaset.source", "ram"),
					resource.TestCheckResourceAttrSet("data.vkcs_compute_quotaset.source", "cores"),
					resource.TestCheckResourceAttrSet("data.vkcs_compute_quotaset.source", "instances"),
					resource.TestCheckResourceAttrSet("data.vkcs_compute_quotaset.source", "server_groups"),
					resource.TestCheckResourceAttrSet("data.vkcs_compute_quotaset.source", "server_group_members"),
				),
			},
		},
	})
}

func testAccCheckComputeQuotasetDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find compute quotaset data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Compute quotaset data source ID not set")
		}

		return nil
	}
}

const testAccComputeQuotasetDataSourceSource = `
data "vkcs_compute_quotaset" "source" {
  project_id = "{{.ProjectID}}"
}
`
