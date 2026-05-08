package resource_cluster

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/templates"
)

func makeUsersListValue(t *testing.T, ctx context.Context, users []ConfigsUsersValue) basetypes.ListValue {
	t.Helper()
	lv, diags := types.ListValueFrom(ctx, ConfigsUsersValue{}.Type(ctx), users)
	require.False(t, diags.HasError(), "ListValueFrom users: %v", diags)
	return lv
}

func makeWarehousesListValue(t *testing.T, ctx context.Context, warehouses []ConfigsWarehousesValue) basetypes.ListValue {
	t.Helper()
	lv, diags := types.ListValueFrom(ctx, ConfigsWarehousesValue{}.Type(ctx), warehouses)
	require.False(t, diags.HasError(), "ListValueFrom warehouses: %v", diags)
	return lv
}

func newWarehouse(ctx context.Context, name string) ConfigsWarehousesValue {
	return ConfigsWarehousesValue{
		Connections: types.ListNull(ConfigsWarehousesConnectionsValue{}.Type(ctx)),
		Id:          types.StringNull(),
		Name:        types.StringValue(name),
		state:       attr.ValueStateKnown,
	}
}

func newUser(username, password, role string) ConfigsUsersValue {
	return ConfigsUsersValue{
		CreatedAt: types.StringNull(),
		Id:        types.StringNull(),
		Username:  types.StringValue(username),
		Password:  types.StringValue(password),
		Role:      types.StringValue(role),
		state:     attr.ValueStateKnown,
	}
}

type podGroupSpec struct {
	name  string
	count int64
}

func makePodGroupsListValueNoVolumes(t *testing.T, ctx context.Context, specs []podGroupSpec) basetypes.ListValue {
	t.Helper()
	pgValues := make([]PodGroupsValue, 0, len(specs))
	for _, s := range specs {
		pgValues = append(pgValues, PodGroupsValue{
			Alias:            types.StringNull(),
			AvailabilityZone: types.StringNull(),
			Count:            types.Int64Value(s.count),
			FloatingIpPool:   types.StringNull(),
			Id:               types.StringNull(),
			Name:             types.StringValue(s.name),
			Resource:         types.ObjectNull(PodGroupsResourceValue{}.AttributeTypes(ctx)),
			Volumes:          types.MapNull(PodGroupsVolumesValue{}.Type(ctx)),
			state:            attr.ValueStateKnown,
		})
	}
	lv, diags := types.ListValueFrom(ctx, PodGroupsValue{}.Type(ctx), pgValues)
	require.False(t, diags.HasError(), "ListValueFrom pod_groups: %v", diags)
	return lv
}

func newConfigsValue(ctx context.Context, users basetypes.ListValue, warehouses basetypes.ListValue) ConfigsValue {
	return ConfigsValue{
		Maintenance: types.ObjectNull(ConfigsMaintenanceValue{}.AttributeTypes(ctx)),
		Settings:    types.ListNull(ConfigsSettingsValue{}.Type(ctx)),
		Users:       users,
		Warehouses:  warehouses,
		state:       attr.ValueStateKnown,
	}
}

func TestExpandClusterConfigsUsers_AddsConnectionStoreCreateFalse(t *testing.T) {
	ctx := context.Background()

	users := []ConfigsUsersValue{
		newUser("vkdata", "Pass#1!", "dbOwner"),
		newUser("vkdata1", "Pass#2!", "readWrite"),
	}
	listVal := makeUsersListValue(t, ctx, users)

	result, diags := ExpandClusterConfigsUsers(ctx, listVal)
	require.False(t, diags.HasError(), "ExpandClusterConfigsUsers: %v", diags)
	require.Len(t, result, 2)

	for i, u := range result {
		require.NotNilf(t, u.ConnectionStore, "user %d (%s): ConnectionStore must be set", i, u.Username)
		assert.Falsef(t, u.ConnectionStore.Create, "user %d (%s): ConnectionStore.Create must be false", i, u.Username)
	}

	assert.Equal(t, "vkdata", result[0].Username)
	assert.Equal(t, "Pass#1!", result[0].Password)
	assert.Equal(t, "dbOwner", result[0].Role)
	assert.Equal(t, "vkdata1", result[1].Username)
	assert.Equal(t, "readWrite", result[1].Role)
}

func TestExpandClusterConfigsUsers_EmptyList(t *testing.T) {
	ctx := context.Background()

	listVal := makeUsersListValue(t, ctx, []ConfigsUsersValue{})

	result, diags := ExpandClusterConfigsUsers(ctx, listVal)
	require.False(t, diags.HasError(), "ExpandClusterConfigsUsers: %v", diags)
	assert.Empty(t, result)
}

func TestExpandClusterConfigs_PropagatesUsernamesToSingleWarehouse(t *testing.T) {
	ctx := context.Background()

	users := makeUsersListValue(t, ctx, []ConfigsUsersValue{
		newUser("vkdata", "Pass#1!", "dbOwner"),
		newUser("vkdata1", "Pass#2!", "readWrite"),
	})
	warehouses := makeWarehousesListValue(t, ctx, []ConfigsWarehousesValue{
		newWarehouse(ctx, "db_customer"),
	})

	result, diags := ExpandClusterConfigs(ctx, newConfigsValue(ctx, users, warehouses))
	require.False(t, diags.HasError(), "ExpandClusterConfigs: %v", diags)
	require.NotNil(t, result)

	require.Len(t, result.Warehouses, 1)
	assert.Equal(t, "db_customer", result.Warehouses[0].Name)
	assert.Equal(t, []string{"vkdata", "vkdata1"}, result.Warehouses[0].Users)

	require.Len(t, result.Users, 2)
	for _, u := range result.Users {
		require.NotNil(t, u.ConnectionStore)
		assert.False(t, u.ConnectionStore.Create)
	}
}

func TestExpandClusterConfigs_PropagatesUsernames_SingleUserSingleWarehouse(t *testing.T) {
	ctx := context.Background()

	users := makeUsersListValue(t, ctx, []ConfigsUsersValue{
		newUser("admin", "Strong#Pass!", "dbOwner"),
	})
	warehouses := makeWarehousesListValue(t, ctx, []ConfigsWarehousesValue{
		newWarehouse(ctx, "main"),
	})

	result, diags := ExpandClusterConfigs(ctx, newConfigsValue(ctx, users, warehouses))
	require.False(t, diags.HasError(), "ExpandClusterConfigs: %v", diags)

	require.Len(t, result.Warehouses, 1)
	assert.Equal(t, []string{"admin"}, result.Warehouses[0].Users)
}

func TestExpandClusterConfigs_NoUsers_WarehouseUsersStayEmpty(t *testing.T) {
	ctx := context.Background()

	users := types.ListNull(ConfigsUsersValue{}.Type(ctx))
	warehouses := makeWarehousesListValue(t, ctx, []ConfigsWarehousesValue{
		newWarehouse(ctx, "main"),
	})

	result, diags := ExpandClusterConfigs(ctx, newConfigsValue(ctx, users, warehouses))
	require.False(t, diags.HasError(), "ExpandClusterConfigs: %v", diags)
	require.NotNil(t, result)

	require.Len(t, result.Warehouses, 1)
	assert.Empty(t, result.Warehouses[0].Users)
	assert.Empty(t, result.Users)
}

func TestExpandClusterConfigs_NoWarehouses_NoPanic(t *testing.T) {
	ctx := context.Background()

	users := makeUsersListValue(t, ctx, []ConfigsUsersValue{
		newUser("vkdata", "Pass#1!", "dbOwner"),
	})
	warehouses := types.ListNull(ConfigsWarehousesValue{}.Type(ctx))

	result, diags := ExpandClusterConfigs(ctx, newConfigsValue(ctx, users, warehouses))
	require.False(t, diags.HasError(), "ExpandClusterConfigs: %v", diags)
	require.NotNil(t, result)

	assert.Empty(t, result.Warehouses)
	require.Len(t, result.Users, 1)
	require.NotNil(t, result.Users[0].ConnectionStore)
	assert.False(t, result.Users[0].ConnectionStore.Create)
}

func TestExpandClusterConfigsUsers_ConnectionStoreSerializesCreateFalse(t *testing.T) {
	user := clusters.ClusterCreateConfigUser{
		Username: "u",
		Password: "p",
		Role:     "dbOwner",
		Settings: []clusters.ClusterCreateConfigSetting{},
		ConnectionStore: &clusters.ClusterCreateConfigUserConnectionStore{
			Create: false,
		},
	}

	body, err := user.Map()
	require.NoError(t, err)

	cs, ok := body["connection_store"].(map[string]interface{})
	require.True(t, ok, "connection_store must be a map, got %T", body["connection_store"])
	createVal, ok := cs["create"]
	require.True(t, ok, "connection_store.create must be present in payload")
	assert.Equal(t, false, createVal, "connection_store.create must serialize as literal false (not omitted)")
}

func TestExpandClusterConfigsUsers_RoleRequired(t *testing.T) {
	// gophercloud BuildRequestBody must return an error when a required field is empty.
	user := clusters.ClusterCreateConfigUser{
		Username: "u",
		Password: "p",
		Role:     "",
	}

	_, err := user.Map()
	require.Error(t, err, "empty role must break BuildRequestBody")
}

func TestExpandClusterConfigsUsers_SettingsAlwaysEmptyArrayInJSON(t *testing.T) {
	// Atom requires users[].settings to be sent as [], not null or absent.
	ctx := context.Background()

	listVal := makeUsersListValue(t, ctx, []ConfigsUsersValue{
		newUser("vkdata", "Pass#1!", "dbOwner"),
	})

	expanded, diags := ExpandClusterConfigsUsers(ctx, listVal)
	require.False(t, diags.HasError())
	require.Len(t, expanded, 1)

	require.NotNil(t, expanded[0].Settings, "Settings must be initialized as empty slice, not nil")
	assert.Empty(t, expanded[0].Settings)

	raw, err := json.Marshal(expanded[0])
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"settings":[]`,
		"users[].settings must be [] in JSON, not null/absent; got: %s", string(raw))
	assert.NotContains(t, string(raw), `"settings":null`)
}

func TestExpandClusterConfigsWarehouses_ExtensionsAlwaysEmptyArrayInJSON(t *testing.T) {
	// Atom requires warehouses[].extensions to be sent as [], not null or absent.
	ctx := context.Background()

	listVal := makeWarehousesListValue(t, ctx, []ConfigsWarehousesValue{
		newWarehouse(ctx, "db_customer"),
	})

	expanded, diags := ExpandClusterConfigsWarehouses(ctx, listVal)
	require.False(t, diags.HasError())
	require.Len(t, expanded, 1)

	require.NotNil(t, expanded[0].Extensions, "Extensions must be initialized as empty slice, not nil")
	assert.Empty(t, expanded[0].Extensions)

	raw, err := json.Marshal(expanded[0])
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"extensions":[]`,
		"warehouses[].extensions must be [] in JSON; got: %s", string(raw))
	assert.NotContains(t, string(raw), `"extensions":null`)
}

func TestExpandClusterPodGroups_VolumesAlwaysEmptyMapInJSON(t *testing.T) {
	// When volumes are not specified in TF config, payload must have "volumes": {} (not null/absent).
	ctx := context.Background()

	template := &templates.ClusterTemplate{
		PodGroups: []templates.ClusterTemplatePodgroup{
			{ID: "tpl-1", Name: "broker"},
		},
	}

	pgList := makePodGroupsListValueNoVolumes(t, ctx, []podGroupSpec{
		{name: "broker", count: 1},
	})

	expanded, diags := ExpandClusterPodGroups(ctx, template, pgList)
	require.False(t, diags.HasError())
	require.Len(t, expanded, 1)

	require.NotNil(t, expanded[0].Volumes, "Volumes must be initialized as empty map, not nil")
	assert.Empty(t, expanded[0].Volumes)

	raw, err := json.Marshal(expanded[0])
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"volumes":{}`,
		"pod_groups[].volumes must be {} in JSON; got: %s", string(raw))
	assert.NotContains(t, string(raw), `"volumes":null`)
}

func TestClusterCreate_FullPayloadMatchesAtomShape(t *testing.T) {
	// Build a full request mirroring the Atom example: verify JSON shape at key spots.
	count := 1
	body := clusters.ClusterCreate{
		Name:              "tf-acc-cluster",
		ClusterTemplateID: "bce51fd9-89b6-4622-8f65-97569a6a848b",
		NetworkID:         "6ceca5c1-57dc-4ebb-b3b3-ce851106e75a",
		SubnetID:          "f8ed1b9c-5b31-4781-858b-cc325442a409",
		ProductName:       "clickhouse",
		ProductVersion:    "25.3.0",
		AvailabilityZone:  "MS1",
		FloatingIPPool:    "auto",
		SDN:               "SPRUT",
		Configs: &clusters.ClusterCreateConfig{
			Maintenance: &clusters.ClusterCreateConfigMaintenance{
				Start: "0 22 * * *",
			},
			Users: []clusters.ClusterCreateConfigUser{
				{
					Username: "vkdata",
					Password: "&I`nm<lx@Lb@A%8j",
					Role:     "dbOwner",
					Settings: []clusters.ClusterCreateConfigSetting{},
					ConnectionStore: &clusters.ClusterCreateConfigUserConnectionStore{
						Create: false,
					},
				},
			},
			Warehouses: []clusters.ClusterCreateConfigWarehouse{
				{
					Name:       "db_customer",
					Users:      []string{"vkdata"},
					Extensions: []clusters.ClusterCreateConfigWarehouseExtension{},
				},
			},
		},
		PodGroups: []clusters.ClusterCreatePodGroup{
			{
				Count:              &count,
				PodGroupTemplateID: "6ce0769a-fd7e-4925-80af-81cf97a556f5",
				Resource: &clusters.ClusterCreatePodGroupResource{
					CPURequest: "0.5",
					RAMRequest: "1",
				},
				Volumes: map[string]clusters.ClusterCreatePodGroupVolume{
					"data": {StorageClassName: "ceph-ssd", Storage: "5", Count: 1},
				},
			},
		},
	}

	raw, err := json.Marshal(body)
	require.NoError(t, err)
	s := string(raw)

	assert.Contains(t, s, `"role":"dbOwner"`)
	assert.Contains(t, s, `"settings":[]`)
	assert.Contains(t, s, `"connection_store":{"create":false}`)
	assert.Contains(t, s, `"users":["vkdata"]`)
	assert.Contains(t, s, `"extensions":[]`)
	assert.Contains(t, s, `"volumes":{"data":`)
	assert.Contains(t, s, `"sdn":"SPRUT"`)
	assert.Contains(t, s, `"floating_ip_pool":"auto"`)
	assert.NotContains(t, s, `"settings":null`)
	assert.NotContains(t, s, `"extensions":null`)
	assert.NotContains(t, s, `"volumes":null`)
}

func TestClusterCreate_SDNAlwaysPresentInJSON(t *testing.T) {
	// sdn is a required field in Atom payload, no omitempty.
	body := clusters.ClusterCreate{Name: "x", SDN: "SPRUT"}
	raw, err := json.Marshal(body)
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"sdn":"SPRUT"`)

	// Even with an empty value, the field must be present (no omitempty).
	emptyBody := clusters.ClusterCreate{Name: "x"}
	raw, err = json.Marshal(emptyBody)
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"sdn":""`,
		"sdn must be present in JSON even as an empty string; got: %s", string(raw))
}

func TestClusterCreate_FloatingIPPoolPayload(t *testing.T) {
	// floating_ip_pool: either "auto" in HCL → "floating_ip_pool":"auto" in JSON,
	// or omitted in HCL (= "" on the Go side) → field omitted from JSON (omitempty).
	t.Run("auto value reaches JSON", func(t *testing.T) {
		body := clusters.ClusterCreate{Name: "x", FloatingIPPool: "auto"}
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		assert.Contains(t, string(raw), `"floating_ip_pool":"auto"`)
	})

	t.Run("empty string is omitted from JSON", func(t *testing.T) {
		body := clusters.ClusterCreate{Name: "x", FloatingIPPool: ""}
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		assert.NotContains(t, string(raw), `"floating_ip_pool"`,
			"floating_ip_pool must not appear in JSON when value is empty; got: %s", string(raw))
	})
}
