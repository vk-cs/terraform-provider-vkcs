package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccKubernetesClusterTemplateDataSource_basic(t *testing.T) {

	version := "1.20.4"
	name := "Kubernetes-centos-v1.20.4-mcs.1"
	uuid := clusterTemplateID
	resourceName := "data.vkcs_kubernetes_clustertemplate.template"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckKubernetes(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config:   testAccKubernetesClusterTemplateBasicByVersion(version),
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesClusterTemplateDataSourceID(resourceName, name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccKubernetesClusterTemplateBasicByName(name),
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesClusterTemplateDataSourceID(resourceName, uuid),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccKubernetesClusterTemplateBasicByUUID(uuid),
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesClusterTemplateDataSourceID(resourceName, uuid),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func testAccKubernetesClusterTemplateDataSourceID(resourceName string, id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ct, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find cluster template data source: %s", resourceName)
		}

		if ct.Primary.ID != id {
			return fmt.Errorf("cluster template data source ID is not set")
		}

		return nil
	}
}

func testAccKubernetesClusterTemplateBasicByVersion(version string) string {
	template := `
data "vkcs_kubernetes_clustertemplate" "template"{
  version = "` + version + `"
}
`
	return template
}

func testAccKubernetesClusterTemplateBasicByName(name string) string {
	template := `
data "vkcs_kubernetes_clustertemplate" "template" {
  name = "` + name + `"
}
`
	return template
}

func testAccKubernetesClusterTemplateBasicByUUID(uuid string) string {
	template := `
data "vkcs_kubernetes_clustertemplate" "template"{
  cluster_template_uuid = "` + uuid + `"
}
`
	return template
}
