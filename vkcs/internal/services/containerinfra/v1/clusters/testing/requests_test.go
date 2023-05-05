package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/stretchr/testify/assert"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/clusters"
)

func TestClusterCreateOpts(t *testing.T) {

	labels := map[string]string{
		"cluster_node_volume_type": "ms1",
		"container_infra_prefix":   "registry.infra.mail.ru:5010/",
	}

	mcount := 2

	createOpts := clusters.CreateOpts{
		ClusterTemplateID:  "95663bae-6763-4a53-9424-831975285cc1",
		Keypair:            "default",
		Labels:             labels,
		MasterFlavorID:     "Basic-1-2-20",
		Name:               "k8s-cluster",
		AvailabilityZone:   "MS1",
		MasterCount:        mcount,
		NetworkID:          "95663bae-6763-4a53-9424-831975285cc1",
		SubnetID:           "95663bae-6763-4a53-9424-831975285cc1",
		FloatingIPEnabled:  false,
		InsecureRegistries: []string{"1.2.3.4", "6.7.8.9:1234"},
		DNSDomain:          "cluster.example",
	}

	b, _ := createOpts.Map()

	assert.IsType(t, map[string]interface{}{}, b["labels"])
	assert.Len(t, b["labels"], 2)
	assert.Len(t, b, 12)
}

func k8sConfigFixture(t *testing.T, id string) {
	switch id {
	case "notfound":
		th.Mux.HandleFunc(fmt.Sprintf("/clusters/%s/kube_config", id), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			w.WriteHeader(http.StatusBadRequest)
		})
	default:
		th.Mux.HandleFunc(fmt.Sprintf("/clusters/%s/kube_config", id), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("content-disposition", "attachment; filename='kubeconfig'")
			w.Header().Add("Content-Length", "7")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "example")
		})
	}
}

func TestK8sConfigGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	k8sConfigFixture(t, "123")

	serviceClient := fake.ServiceClient()
	config, err := clusters.KubeConfigGet(serviceClient, "123")
	assert.NoError(t, err)
	assert.EqualValues(t, "example", config)
}

func TestK8sConfigGetError(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	k8sConfigFixture(t, "notfound")

	serviceClient := fake.ServiceClient()
	_, err := clusters.KubeConfigGet(serviceClient, "notfound")
	assert.Error(t, err)
}
