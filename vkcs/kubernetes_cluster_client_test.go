package vkcs

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/stretchr/testify/assert"
)

func TestClusterCreateOpts(t *testing.T) {

	labels := map[string]string{
		"cluster_node_volume_type": "ms1",
		"container_infra_prefix":   "registry.infra.mail.ru:5010/",
	}

	mcount := 2

	createOpts := clusterCreateOpts{
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
	}

	b, _ := createOpts.Map()

	assert.IsType(t, map[string]interface{}{}, b["labels"])
	assert.Len(t, b["labels"], 2)
	assert.Len(t, b, 11)
}

func TestPatchOpts(t *testing.T) {

	patchOpts := nodeGroupClusterPatchOpts{
		{
			Path:  "/max_nodes",
			Value: 10,
			Op:    "replace",
		},
		{
			Path:  "/min_nodes",
			Value: 2,
			Op:    "replace",
		},
		{
			Path:  "/autoscaling_enabled",
			Value: false,
			Op:    "replace",
		},
	}

	b, _ := patchOpts.PatchMap()

	assert.IsType(t, []map[string]interface{}{}, b)
	assert.Len(t, b, 3)
}

func TestAddBatchOpts(t *testing.T) {

	addGroups := []nodeGroup{
		{
			Name:     "test1",
			FlavorID: "95663bae-6763-4a53-9424-831975285cc1",
		},
		{
			Name:     "test2",
			FlavorID: "95663bae-6763-4a53-9424-831975285cc1",
		},
	}

	addbatchOpts := nodeGroupBatchAddParams{
		Action:  "batch_add_ng",
		Payload: addGroups,
	}

	b, _ := addbatchOpts.Map()

	assert.IsType(t, []interface{}{}, b["payload"])
	assert.Len(t, b["payload"], 2)
	assert.Len(t, b, 2)
}

func k8sconfigFixture(t *testing.T, id string) {
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

	k8sconfigFixture(t, "123")

	serviceClient := fake.ServiceClient()
	config, err := k8sConfigGet(serviceClient, "123")
	assert.NoError(t, err)
	assert.EqualValues(t, "example", config)
}

func TestK8sConfigGetError(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	k8sconfigFixture(t, "notfound")

	serviceClient := fake.ServiceClient()
	_, err := k8sConfigGet(serviceClient, "notfound")
	assert.Error(t, err)
}
