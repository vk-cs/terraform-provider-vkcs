package vkcs

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

const clusterResourceFixture = `
		resource "vkcs_kubernetes_cluster" "%[1]s" {
		  name = "%[1]s"
		  cluster_template_id = "%s"
		  master_flavor       = "%s"
		  master_count        =  "%d"
		  keypair = "%s"
          network_id = "%s"
          subnet_id = "%s"
          floating_ip_enabled = false
          availability_zone = "%s"
		}
`

func clusterFixture(name, clusterTemplateID, flavorID, keypair,
	networkID, subnetID, az string, masterCount int) *clusterCreateOpts {
	return &clusterCreateOpts{
		Name:              name,
		MasterCount:       masterCount,
		ClusterTemplateID: clusterTemplateID,
		MasterFlavorID:    flavorID,
		Keypair:           keypair,
		NetworkID:         networkID,
		SubnetID:          subnetID,
		AvailabilityZone:  az,
	}
}

func checkClusterAttrs(resourceName string, cluster *clusterCreateOpts) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.Empty() == true {
			return fmt.Errorf("state not updated")
		}

		checksStore := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "name", cluster.Name),
			resource.TestCheckResourceAttr(resourceName, "master_count", strconv.Itoa(cluster.MasterCount)),
			resource.TestCheckResourceAttr(resourceName, "cluster_template_id", cluster.ClusterTemplateID),
		}

		return resource.ComposeTestCheckFunc(checksStore...)(s)
	}
}

func TestMockAccKubernetesCluster_basic(t *testing.T) {
	clientFixture := &ContainerClientFixture{}
	clusterUUID := uuid.NewV4().String()

	// Mock config methods
	dummyConfig := &dummyConfig{}
	dummyConfig.On("LoadAndValidate").Return(nil)
	dummyConfig.On("ContainerInfraV1Client", "").Return(clientFixture, nil)
	dummyConfig.On("getRegion").Return("")

	// Create cluster fixtures
	clusterName := "testcluster" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "vkcs_kubernetes_cluster." + clusterName

	createClusterFixture := clusterFixture(clusterName, clusterTemplateID, osFlavorID,
		osKeypairName, osNetworkID, osSubnetworkID, "MS1", 1)
	jsonClusterFixture, _ := createClusterFixture.Map()

	scaleFlavorClusterFixture := clusterFixture(clusterName, clusterTemplateID, osNewFlavorID,
		osKeypairName, osNetworkID, osSubnetworkID, "MS1", 1)
	scaleRequestFixture := map[string]interface{}{"action": "resize_masters", "payload": map[string]interface{}{"flavor": scaleFlavorClusterFixture.MasterFlavorID}}
	jsonClusterScaleFixture, _ := scaleFlavorClusterFixture.Map()

	// Mock API calls
	clientFixture.On("ServiceURL", []string{"clusters"}).Return(testAccURL)
	clientFixture.On("ServiceURL", []string{"clusters", clusterUUID}).Return(testAccURL)
	clientFixture.On("ServiceURL", []string{"clusters", clusterUUID, "actions"}).Return(testAccURL)
	// Create cluster
	clientFixture.On("Post", testAccURL+"/clusters", jsonClusterFixture, mock.Anything, getRequestOpts(202)).Return(makeClusterCreateResponseFixture(clusterUUID), nil)
	// Check it's status
	clientFixture.On("Get", testAccURL+"/clusters/"+clusterUUID, mock.Anything, getRequestOpts(200)).Return(makeClusterGetResponseFixture(jsonClusterFixture, clusterUUID, clusterStatusRunning), nil).Times(6)
	// Update cluster
	clientFixture.On("Post", testAccURL+"/clusters/"+clusterUUID+"/actions", scaleRequestFixture, mock.Anything, getRequestOpts(200, 202)).Return(makeClusterGetResponseFixture(jsonClusterScaleFixture, clusterUUID, clusterStatusRunning), nil)
	// Check it's status
	clientFixture.On("Get", testAccURL+"/clusters/"+clusterUUID, mock.Anything, getRequestOpts(200)).Return(makeClusterGetResponseFixture(jsonClusterScaleFixture, clusterUUID, clusterStatusRunning), nil).Times(5)
	// Delete cluster
	clientFixture.On("Delete", testAccURL+"/clusters/"+clusterUUID, getRequestOpts()).Return(makeClusterDeleteResponseFixture(), nil)
	// Check deleted
	clientFixture.On("Get", testAccURL+"/clusters/"+clusterUUID, mock.Anything, getRequestOpts(200)).Return(gophercloud.ErrDefault404{}).Twice()

	var cluster cluster

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckKubernetes(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterBasic(createClusterFixture),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterExists(resourceName, &cluster),
					checkClusterAttrs(resourceName, createClusterFixture),
				),
			},
			{
				Config: testAccKubernetesClusterBasic(scaleFlavorClusterFixture),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "master_flavor", scaleFlavorClusterFixture.MasterFlavorID),
					testAccCheckKubernetesClusterScaled(resourceName),
				),
			},
		},
	})
}

func testAccCheckKubernetesClusterExists(n string, cluster *cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, found, err := getClusterAndResource(n, s)
		if err != nil {
			return err
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		if found.UUID != rs.Primary.ID {
			return fmt.Errorf("cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCheckKubernetesClusterScaled(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, found, err := getClusterAndResource(n, s)
		if err != nil {
			return err
		}

		if found.MasterFlavorID != rs.Primary.Attributes["master_flavor"] {
			return fmt.Errorf("cluster flavor not changed")
		}
		return nil
	}
}

func getClusterAndResource(n string, s *terraform.State) (*terraform.ResourceState, *cluster, error) {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return nil, nil, fmt.Errorf("cluster not found: %s", n)
	}

	config := testAccProvider.Meta().(configer)
	containerInfraClient, err := config.ContainerInfraV1Client(osRegionName)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating container infra client: %s", err)
	}

	found, err := clusterGet(containerInfraClient, rs.Primary.ID).Extract()
	if err != nil {
		return nil, nil, err
	}
	return rs, found, nil
}

func testAccCheckKubernetesClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	containerInfraClient, err := config.ContainerInfraV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("error creating container infra client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_kubernetes_cluster" {
			continue
		}

		_, err := clusterGet(containerInfraClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("сluster still exists")
		}
	}

	return nil
}

func testAccKubernetesClusterBasic(createOpts *clusterCreateOpts) string {

	return fmt.Sprintf(
		clusterResourceFixture,
		createOpts.Name,
		createOpts.ClusterTemplateID,
		createOpts.MasterFlavorID,
		createOpts.MasterCount,
		createOpts.Keypair,
		createOpts.NetworkID,
		createOpts.SubnetID,
		createOpts.AvailabilityZone,
	)
}
