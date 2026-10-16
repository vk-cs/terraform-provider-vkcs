package provider_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/provider"
)

func TestKMSFrameworkMigration(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(provider.Provider())()
	schemas, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, schemas.Diagnostics)
	sdk := provider.SDKProviderBase()
	for _, name := range []string{"vkcs_kms_key", "vkcs_kms_secret"} {
		require.Contains(t, schemas.ResourceSchemas, name)
		require.NotContains(t, sdk.ResourcesMap, name)
	}
	for _, name := range []string{"vkcs_kms_key", "vkcs_kms_key_encrypt", "vkcs_kms_key_decrypt", "vkcs_kms_keys_list", "vkcs_kms_secret", "vkcs_kms_secrets_list"} {
		require.Contains(t, schemas.DataSourceSchemas, name)
		require.NotContains(t, sdk.DataSourcesMap, name)
	}
	for _, tc := range []struct{ name, state string }{
		{"vkcs_kms_key", `{"id":"existing-key","name":"existing-key","type":"aes256-gcm96","deletion_allowed":false,"timeouts":null}`},
		{"vkcs_kms_secret", `{"id":"existing-secret","path":"existing-secret","data_json":"{\"password\":\"value\"}","created_time":"2026-09-01 10:00:00 +0000 UTC","version":2,"delete_protection":true,"timeouts":{"default":"2m"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := server.UpgradeResourceState(ctx, &tfprotov6.UpgradeResourceStateRequest{TypeName: tc.name, Version: 0, RawState: &tfprotov6.RawState{JSON: []byte(tc.state)}})
			require.NoError(t, err)
			require.Empty(t, resp.Diagnostics)
			require.NotNil(t, resp.UpgradedState)
			schema := schemas.ResourceSchemas[tc.name]
			value, err := resp.UpgradedState.Unmarshal(schema.ValueType())
			require.NoError(t, err)
			// The removed timeout setting is discarded when upgrading SDK state.
			var expected map[string]any
			require.NoError(t, json.Unmarshal([]byte(tc.state), &expected))
			delete(expected, "timeouts")
			expectedJSON, err := json.Marshal(expected)
			require.NoError(t, err)
			original, err := (&tfprotov6.DynamicValue{JSON: expectedJSON}).Unmarshal(schema.ValueType())
			require.NoError(t, err)
			require.True(t, original.Equal(value), "SDK state changed during migration")
		})
	}
}

func TestKMSFrameworkPlanning(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(provider.Provider())()
	schemas, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, schemas.Diagnostics)
	for _, tc := range []struct {
		name, config, prior, proposed string
		defaults                      map[string]any
		replaces                      bool
	}{
		{name: "vkcs_kms_key", config: `{"id":null,"name":"key","type":null,"deletion_allowed":null}`, defaults: map[string]any{"type": "aes256-gcm96", "deletion_allowed": false}},
		{name: "vkcs_kms_secret", config: `{"id":null,"path":"secret","data_json":"{}","created_time":null,"version":null,"delete_protection":null}`, defaults: map[string]any{"delete_protection": true}},
		{name: "vkcs_kms_key", config: `{"id":null,"name":"key","type":"aes256-gcm96","deletion_allowed":true}`, prior: `{"id":"key","name":"key","type":"aes256-gcm96","deletion_allowed":false}`, proposed: `{"id":"key","name":"key","type":"aes256-gcm96","deletion_allowed":true}`, defaults: map[string]any{"id": "key", "deletion_allowed": true}},
		{name: "vkcs_kms_key", config: `{"id":null,"name":"renamed","type":"aes256-gcm96","deletion_allowed":false}`, prior: `{"id":"key","name":"key","type":"aes256-gcm96","deletion_allowed":false}`, proposed: `{"id":"key","name":"renamed","type":"aes256-gcm96","deletion_allowed":false}`, replaces: true},
	} {
		t.Run(tc.name+tc.config, func(t *testing.T) {
			prior := tc.prior
			if prior == "" {
				prior = "null"
			}
			proposed := tc.proposed
			if proposed == "" {
				proposed = tc.config
			}
			resp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{TypeName: tc.name, Config: &tfprotov6.DynamicValue{JSON: []byte(tc.config)}, PriorState: &tfprotov6.DynamicValue{JSON: []byte(prior)}, ProposedNewState: &tfprotov6.DynamicValue{JSON: []byte(proposed)}})
			require.NoError(t, err)
			require.Empty(t, resp.Diagnostics)
			require.Equal(t, tc.replaces, len(resp.RequiresReplace) > 0)
			value, err := resp.PlannedState.Unmarshal(schemas.ResourceSchemas[tc.name].ValueType())
			require.NoError(t, err)
			var attributes map[string]tftypes.Value
			require.NoError(t, value.As(&attributes))
			for name, want := range tc.defaults {
				switch want := want.(type) {
				case string:
					var got string
					require.NoError(t, attributes[name].As(&got))
					require.Equal(t, want, got)
				case bool:
					var got bool
					require.NoError(t, attributes[name].As(&got))
					require.Equal(t, want, got)
				}
			}
		})
	}
}
