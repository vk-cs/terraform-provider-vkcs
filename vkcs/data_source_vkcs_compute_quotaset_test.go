package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeQuotasetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeQuotasetDataSourceSource(),
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

func testAccComputeQuotasetDataSourceSource() string {
	return fmt.Sprintf(`
data "vkcs_compute_quotaset" "source" {
  project_id = "%s"
}
`, osProjectID)
}
