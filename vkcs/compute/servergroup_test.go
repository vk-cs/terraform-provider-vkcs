package compute_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/compute"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	th "github.com/gophercloud/gophercloud/testhelper"
	thclient "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestComputeServerGroupCreateOpts(t *testing.T) {
	createOpts := compute.ComputeServerGroupCreateOpts{
		servergroups.CreateOpts{
			Name:     "foo",
			Policies: []string{"affinity"},
		},
		map[string]string{
			"foo": "bar",
		},
	}

	expected := map[string]interface{}{
		"server_group": map[string]interface{}{
			"name":     "foo",
			"policies": []interface{}{"affinity"},
			"foo":      "bar",
		},
	}

	actual, err := createOpts.ToServerGroupCreateMap()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExpandComputeServerGroupPoliciesMicroversions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	raw := []interface{}{
		"affinity",
		"soft-anti-affinity",
		"soft-affinity",
		"custom-policy",
	}
	client := thclient.ServiceClient()

	expectedPolicies := []string{
		"affinity",
		"soft-anti-affinity",
		"soft-affinity",
		"custom-policy",
	}
	expectedMicroversion := "2.15"

	actualPolicies := compute.ExpandComputeServerGroupPolicies(client, raw)
	actualMicroversion := client.Microversion

	assert.Equal(t, expectedMicroversion, actualMicroversion)
	assert.Equal(t, expectedPolicies, actualPolicies)
}

func TestExpandComputeServerGroupPoliciesMicroversionsLegacy(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	raw := []interface{}{
		"anti-affinity",
		"affinity",
	}
	client := thclient.ServiceClient()

	expectedPolicies := []string{
		"anti-affinity",
		"affinity",
	}
	expectedMicroversion := ""

	actualPolicies := compute.ExpandComputeServerGroupPolicies(client, raw)
	actualMicroversion := client.Microversion

	assert.Equal(t, expectedMicroversion, actualMicroversion)
	assert.Equal(t, expectedPolicies, actualPolicies)
}
