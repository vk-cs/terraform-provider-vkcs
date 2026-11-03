package restores

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestCreateRequestMapsParametersByResourceKey(t *testing.T) {
	enableFIP := false
	volumeAutoresizeEnabled := int64(0)
	cloudMonitoringEnabled := true
	request := RestoreCreate{
		Restore: &RestoreCreateOpts{
			ProviderID:   "nova-provider-id",
			CheckpointID: "checkpoint-id",
			Parameters: map[string]RestoreParameters{
				"OS::Nova::Server#server-id": {
					RestoreName:            "restored-server",
					Destination:            "new",
					RestoreNetID:           "network-id",
					RestoreKeypair:         "keypair",
					Flavor:                 "flavor-id",
					SecurityGroups:         []string{"default"},
					EnableFIP:              &enableFIP,
					VolumeAutoresizeEnabled: &volumeAutoresizeEnabled,
					CloudMonitoringEnabled: &cloudMonitoringEnabled,
				},
			},
		},
	}

	body, err := request.Map()
	if err != nil {
		t.Fatalf("map restore request: %v", err)
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal restore request: %v", err)
	}

	var decoded struct {
		Restore struct {
			ProviderID   string `json:"provider_id"`
			CheckpointID string `json:"checkpoint_id"`
			Parameters   map[string]struct {
				RestoreName    string   `json:"restore_name"`
				Destination    string   `json:"destination"`
				RestoreNetID   string   `json:"restore_net_id"`
				SecurityGroups []string `json:"security_groups"`
				EnableFIP      *bool    `json:"enable_fip"`
				AutoResize     *int64   `json:"volume_autoresize_enabled"`
			} `json:"parameters"`
		} `json:"restore"`
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal restore request: %v", err)
	}
	if decoded.Restore.ProviderID != "nova-provider-id" || decoded.Restore.CheckpointID != "checkpoint-id" {
		t.Fatalf("unexpected restore request: %#v", decoded.Restore)
	}
	parameters, ok := decoded.Restore.Parameters["OS::Nova::Server#server-id"]
	if !ok {
		t.Fatalf("restore parameters missing VM resource key: %#v", decoded.Restore.Parameters)
	}
	if parameters.RestoreName != "restored-server" || parameters.Destination != "new" || parameters.RestoreNetID != "network-id" {
		t.Fatalf("unexpected VM restore parameters: %#v", parameters)
	}
	if parameters.EnableFIP == nil || *parameters.EnableFIP {
		t.Fatalf("explicit false enable_fip was not preserved: %#v", parameters.EnableFIP)
	}
	if parameters.AutoResize == nil || *parameters.AutoResize != 0 {
		t.Fatalf("explicit zero volume_autoresize_enabled was not preserved: %#v", parameters.AutoResize)
	}
}

func TestRestoreResponseExtractsPerResourceStatuses(t *testing.T) {
	var response RestoreAPIObject
	err := json.Unmarshal([]byte(`{"restore":{"id":"restore-id","checkpoint_id":"checkpoint-id","provider_id":"nova-provider-id","parameters":{"OS::Nova::Server#server-id":{"restore_name":"restored-server"}},"resources_status":{"server-id":"error"},"resources_reason":{"server-id":"quota exceeded"}}}`), &response)
	if err != nil {
		t.Fatalf("unmarshal restore response: %v", err)
	}
	if response.Restore.ID != "restore-id" {
		t.Fatalf("unexpected restore ID: %q", response.Restore.ID)
	}
	if response.Restore.ResourcesStatus["server-id"] != "error" {
		t.Fatalf("unexpected resources_status: %#v", response.Restore.ResourcesStatus)
	}
	if response.Restore.ResourcesReason["server-id"] != "quota exceeded" {
		t.Fatalf("unexpected resources_reason: %#v", response.Restore.ResourcesReason)
	}
}

func TestListRestoreResponseIncludesMatchingParameters(t *testing.T) {
	var response ListRestoreResponse
	err := json.Unmarshal([]byte(`{"restores":[{"id":"restore-id","checkpoint_id":"checkpoint-id","provider_id":"nova-provider-id","parameters":{"OS::Nova::Server#server-id":{"restore_name":"restored-server","destination":"new"}},"resources_status":{"server-id":"error"}}]}`), &response)
	if err != nil {
		t.Fatalf("unmarshal list restore response: %v", err)
	}
	if len(response.Restores) != 1 || response.Restores[0].ID != "restore-id" {
		t.Fatalf("unexpected list restore response: %#v", response.Restores)
	}
	if response.Restores[0].Parameters["OS::Nova::Server#server-id"].RestoreName != "restored-server" {
		t.Fatalf("restore parameters missing from list response: %#v", response.Restores[0].Parameters)
	}
}

func TestListOptsEncodesPagination(t *testing.T) {
	query, err := (ListOpts{Limit: 100, Offset: 200, Sort: "desc"}).ToListQuery()
	if err != nil {
		t.Fatalf("build restore list query: %v", err)
	}
	parsed, err := url.Parse("https://example.com" + query)
	if err != nil {
		t.Fatalf("parse restore list query: %v", err)
	}
	values := parsed.Query()
	if values.Get("limit") != "100" || values.Get("offset") != "200" || values.Get("sort") != "desc" {
		t.Fatalf("unexpected restore list query: %s", query)
	}
}
