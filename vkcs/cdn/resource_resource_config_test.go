package cdn

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/resource_resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/resources"
)

func TestSelectConfiguredResponseHeaderOptionAllowsCanonicalLegacyReadBack(t *testing.T) {
	ctx := context.Background()
	legacyState, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		StaticHeaders: &resources.ResourceOptionsStringMapOption{
			Enabled: true,
			Value:   map[string]string{"X-Old": "old"},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating legacy state returned diagnostics: %v", diags)
	}
	newConfig, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		StaticResponseHeaders: &resources.ResourceOptionsStaticResponseHeadersOption{
			Enabled: true,
			Value: []resources.StaticResponseHeaderValue{
				{Name: "X-New", Value: []string{"new"}, Always: true},
			},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating new configuration returned diagnostics: %v", diags)
	}

	plan := legacyState
	plan.StaticResponseHeaders = newConfig.StaticResponseHeaders
	selectConfiguredResponseHeaderOption(ctx, &plan, newConfig, &legacyState)
	if !plan.StaticHeaders.IsUnknown() {
		t.Fatalf("planned static_headers = %s, want unknown while API canonicalizes the new option", plan.StaticHeaders)
	}
	apiOptions, diags := plan.ToResourceOptions(ctx)
	if diags.HasError() {
		t.Fatalf("ToResourceOptions returned diagnostics: %v", diags)
	}
	if apiOptions.StaticHeaders != nil {
		t.Fatalf("legacy API option = %#v, want omitted while its alias is unknown", apiOptions.StaticHeaders)
	}
	if apiOptions.StaticResponseHeaders == nil || !apiOptions.StaticResponseHeaders.Enabled || len(apiOptions.StaticResponseHeaders.Value) != 1 || apiOptions.StaticResponseHeaders.Value[0].Name != "X-New" {
		t.Fatalf("new API option = %#v, want the configured X-New header", apiOptions.StaticResponseHeaders)
	}
}

func TestSelectConfiguredLegacyOptionAllowsCanonicalAPIReadBack(t *testing.T) {
	ctx := context.Background()
	legacyOnly, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		StaticHeaders: &resources.ResourceOptionsStringMapOption{
			Enabled: true,
			Value:   map[string]string{"X-Legacy": "legacy"},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating legacy-only options returned diagnostics: %v", diags)
	}
	plan := legacyOnly
	selectConfiguredResponseHeaderOption(ctx, &plan, legacyOnly, nil)
	if !plan.StaticResponseHeaders.IsUnknown() {
		t.Fatalf("create plan static_response_headers = %s, want unknown for API canonicalization", plan.StaticResponseHeaders)
	}

	readBack, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		StaticHeaders: &resources.ResourceOptionsStringMapOption{
			Enabled: true,
			Value:   map[string]string{"X-Legacy": "legacy"},
		},
		StaticResponseHeaders: &resources.ResourceOptionsStaticResponseHeadersOption{
			Enabled: true,
			Value: []resources.StaticResponseHeaderValue{
				{Name: "X-Legacy", Value: []string{"legacy"}, Always: true},
			},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating API read-back returned diagnostics: %v", diags)
	}
	noOpPlan := readBack
	selectConfiguredResponseHeaderOption(ctx, &noOpPlan, legacyOnly, &readBack)
	if !noOpPlan.StaticResponseHeaders.Equal(readBack.StaticResponseHeaders) {
		t.Fatalf("no-op plan static_response_headers = %s, want retained API read-back %s", noOpPlan.StaticResponseHeaders, readBack.StaticResponseHeaders)
	}
}

func TestSelectConfiguredResponseHeaderOptionMakesChangedComputedOptionUnknown(t *testing.T) {
	ctx := context.Background()
	newState, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		StaticResponseHeaders: &resources.ResourceOptionsStaticResponseHeadersOption{
			Enabled: true,
			Value: []resources.StaticResponseHeaderValue{
				{Name: "X-Old", Value: []string{"old"}, Always: true},
			},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating new state returned diagnostics: %v", diags)
	}
	legacyConfig, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		StaticHeaders: &resources.ResourceOptionsStringMapOption{
			Enabled: true,
			Value:   map[string]string{"X-New": "new"},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating legacy configuration returned diagnostics: %v", diags)
	}

	plan := newState
	plan.StaticHeaders = legacyConfig.StaticHeaders
	selectConfiguredResponseHeaderOption(ctx, &plan, legacyConfig, &newState)
	if !plan.StaticResponseHeaders.IsUnknown() {
		t.Fatalf("planned static_response_headers = %s, want unknown during legacy option change", plan.StaticResponseHeaders)
	}
	apiOptions, diags := plan.ToResourceOptions(ctx)
	if diags.HasError() {
		t.Fatalf("ToResourceOptions returned diagnostics: %v", diags)
	}
	if apiOptions.StaticResponseHeaders != nil {
		t.Fatalf("new API option = %#v, want omitted while its value is unknown", apiOptions.StaticResponseHeaders)
	}
	if apiOptions.StaticHeaders == nil || !apiOptions.StaticHeaders.Enabled || apiOptions.StaticHeaders.Value["X-New"] != "new" {
		t.Fatalf("legacy API option = %#v, want configured X-New header", apiOptions.StaticHeaders)
	}
}

func TestResponseHeaderAPIRequestOmitsCanonicalOptionForLegacyConfig(t *testing.T) {
	ctx := context.Background()
	plan, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		StaticHeaders: &resources.ResourceOptionsStringMapOption{
			Enabled: true,
			Value:   map[string]string{"X-Legacy": "configured"},
		},
		StaticResponseHeaders: &resources.ResourceOptionsStaticResponseHeadersOption{
			Enabled: true,
			Value: []resources.StaticResponseHeaderValue{
				{Name: "X-Legacy", Value: []string{"configured"}, Always: true},
			},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating plan returned diagnostics: %v", diags)
	}
	legacyConfig := plan
	legacyConfig.StaticResponseHeaders, diags = resource_resource.NewStaticResponseHeadersValueNull().ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("creating legacy-only configuration returned diagnostics: %v", diags)
	}

	requestOptions := responseHeaderOptionsForAPIRequest(ctx, plan, legacyConfig)
	apiOptions, diags := requestOptions.ToResourceOptions(ctx)
	if diags.HasError() {
		t.Fatalf("ToResourceOptions returned diagnostics: %v", diags)
	}
	if apiOptions.StaticHeaders == nil || !apiOptions.StaticHeaders.Enabled {
		t.Fatalf("legacy API option = %#v, want enabled", apiOptions.StaticHeaders)
	}
	if apiOptions.StaticResponseHeaders != nil {
		t.Fatalf("canonical API option = %#v, want omitted when legacy is configured", apiOptions.StaticResponseHeaders)
	}
}

func TestResponseHeaderAPIRequestOmitsLegacyAliasForNewConfig(t *testing.T) {
	ctx := context.Background()
	plan, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		StaticHeaders: &resources.ResourceOptionsStringMapOption{
			Enabled: false,
			Value:   map[string]string{},
		},
		StaticResponseHeaders: &resources.ResourceOptionsStaticResponseHeadersOption{
			Enabled: true,
			Value: []resources.StaticResponseHeaderValue{
				{Name: "X-New", Value: []string{"new"}, Always: true},
			},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating plan returned diagnostics: %v", diags)
	}
	newConfig := plan
	newConfig.StaticHeaders, diags = resource_resource.NewStaticHeadersValueNull().ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("creating new-only configuration returned diagnostics: %v", diags)
	}

	requestOptions := responseHeaderOptionsForAPIRequest(ctx, plan, newConfig)
	apiOptions, diags := requestOptions.ToResourceOptions(ctx)
	if diags.HasError() {
		t.Fatalf("ToResourceOptions returned diagnostics: %v", diags)
	}
	if apiOptions.StaticHeaders != nil {
		t.Fatalf("legacy API option = %#v, want omitted when new option is configured", apiOptions.StaticHeaders)
	}
	if apiOptions.StaticResponseHeaders == nil || !apiOptions.StaticResponseHeaders.Enabled || len(apiOptions.StaticResponseHeaders.Value) != 1 {
		t.Fatalf("canonical API option = %#v, want enabled X-New header", apiOptions.StaticResponseHeaders)
	}
}

func TestMergeResourceOptionsPreservesUnplannedOptionsAndPrefersPlan(t *testing.T) {
	prior := &resources.ResourceOptions{
		AllowedHttpMethods: &resources.ResourceOptionsAllowedHttpMethodsOption{
			Enabled: true,
			Value:   []resources.ResourceAllowedHttpMethod{resources.ResourceAllowedHttpMethodGet},
		},
		StaticHeaders: &resources.ResourceOptionsStringMapOption{
			Enabled: true,
			Value:   map[string]string{"X-Old": "old"},
		},
	}
	planned := &resources.ResourceOptions{
		StaticHeaders: &resources.ResourceOptionsStringMapOption{
			Enabled: false,
			Value:   map[string]string{},
		},
		StaticResponseHeaders: &resources.ResourceOptionsStaticResponseHeadersOption{
			Enabled: true,
			Value: []resources.StaticResponseHeaderValue{
				{Name: "X-New", Value: []string{"new"}, Always: true},
			},
		},
	}

	merged := mergeResourceOptions(planned, prior)
	if merged.AllowedHttpMethods == nil || !merged.AllowedHttpMethods.Enabled || len(merged.AllowedHttpMethods.Value) != 1 || merged.AllowedHttpMethods.Value[0] != resources.ResourceAllowedHttpMethodGet {
		t.Fatalf("AllowedHttpMethods = %#v, want prior value retained", merged.AllowedHttpMethods)
	}
	if merged.StaticHeaders == nil || merged.StaticHeaders.Enabled || len(merged.StaticHeaders.Value) != 0 {
		t.Fatalf("StaticHeaders = %#v, want the disabled planned value to replace prior state", merged.StaticHeaders)
	}
	if merged.StaticResponseHeaders == nil || !merged.StaticResponseHeaders.Enabled || len(merged.StaticResponseHeaders.Value) != 1 || merged.StaticResponseHeaders.Value[0].Name != "X-New" {
		t.Fatalf("StaticResponseHeaders = %#v, want planned value", merged.StaticResponseHeaders)
	}
}

func TestValidateStaticResponseHeadersRequiresValueUnlessDisabled(t *testing.T) {
	ctx := context.Background()
	attributeTypes := resource_resource.StaticResponseHeadersValue{}.AttributeTypes(ctx)
	value := types.ListNull(resource_resource.ValueValue{}.Type(ctx))

	for _, tc := range []struct {
		name      string
		enabled   types.Bool
		wantError bool
	}{
		{name: "enabled", enabled: types.BoolValue(true), wantError: true},
		{name: "enabled inferred", enabled: types.BoolNull(), wantError: true},
		{name: "enabled unknown", enabled: types.BoolUnknown(), wantError: false},
		{name: "disabled", enabled: types.BoolValue(false), wantError: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			option := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
				"enabled": tc.enabled,
				"value":   value,
			})
			diags := validateStaticResponseHeaders(ctx, option)
			if diags.HasError() != tc.wantError {
				t.Fatalf("HasError() = %t, want %t; diagnostics: %v", diags.HasError(), tc.wantError, diags)
			}
		})
	}
}

func TestValidateStaticResponseHeadersAllowsNullAndUnknown(t *testing.T) {
	ctx := context.Background()
	attributeTypes := resource_resource.StaticResponseHeadersValue{}.AttributeTypes(ctx)

	for _, value := range []types.Object{
		types.ObjectNull(attributeTypes),
		types.ObjectUnknown(attributeTypes),
	} {
		if diags := validateStaticResponseHeaders(ctx, value); diags.HasError() {
			t.Fatalf("unexpected validation diagnostics: %v", diags)
		}
	}
}

func TestPreserveUnconfiguredOptionsUsesPriorStateForUnknownOption(t *testing.T) {
	ctx := context.Background()
	previous, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		GzipCompression: &resources.ResourceOptionsStringListOption{
			Enabled: false,
			Value:   []string{"text/html", "application/json"},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating prior options returned diagnostics: %v", diags)
	}
	configured, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		AllowedHttpMethods: &resources.ResourceOptionsAllowedHttpMethodsOption{
			Enabled: true,
			Value:   []resources.ResourceAllowedHttpMethod{resources.ResourceAllowedHttpMethodGet},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating configured options returned diagnostics: %v", diags)
	}
	plan := configured
	plan.GzipCompression, diags = resource_resource.NewGzipCompressionValueUnknown().ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("creating unknown planned gzip option returned diagnostics: %v", diags)
	}
	plan.BrowserCacheSettings, diags = resource_resource.NewBrowserCacheSettingsValueUnknown().ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("creating unknown planned browser cache option returned diagnostics: %v", diags)
	}

	if diags := preserveUnconfiguredOptions(ctx, &plan, configured, previous); diags.HasError() {
		t.Fatalf("preserveUnconfiguredOptions returned diagnostics: %v", diags)
	}
	if !plan.GzipCompression.Equal(previous.GzipCompression) {
		t.Fatalf("planned gzip option = %s, want prior state %s", plan.GzipCompression, previous.GzipCompression)
	}
	if !plan.BrowserCacheSettings.IsNull() {
		t.Fatalf("planned browser cache option = %s, want the prior null value", plan.BrowserCacheSettings)
	}
}

func TestPreserveUnconfiguredOptionsUsesPriorStateWhenOptionsAreOmitted(t *testing.T) {
	ctx := context.Background()
	previous, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		AllowedHttpMethods: &resources.ResourceOptionsAllowedHttpMethodsOption{
			Enabled: true,
			Value:   []resources.ResourceAllowedHttpMethod{resources.ResourceAllowedHttpMethodGet},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating prior options returned diagnostics: %v", diags)
	}
	plan := resource_resource.NewOptionsValueUnknown()
	config := resource_resource.NewOptionsValueNull()

	if diags := preserveUnconfiguredOptions(ctx, &plan, config, previous); diags.HasError() {
		t.Fatalf("preserveUnconfiguredOptions returned diagnostics: %v", diags)
	}
	if !plan.Equal(previous) {
		t.Fatalf("planned options = %s, want prior state %s", plan, previous)
	}
}

func TestPlanDisabledShieldingUsesZeroPopID(t *testing.T) {
	ctx := context.Background()
	configurationObject, diags := types.ObjectValue(resource_resource.ShieldingValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"enabled": types.BoolValue(false),
		"pop_id":  types.Int64Unknown(),
	})
	if diags.HasError() {
		t.Fatalf("creating disabled shielding config returned diagnostics: %v", diags)
	}
	configurationValue, diags := resource_resource.ShieldingType{}.ValueFromObject(ctx, configurationObject)
	if diags.HasError() {
		t.Fatalf("decoding disabled shielding config returned diagnostics: %v", diags)
	}
	configuration := configurationValue.(resource_resource.ShieldingValue)
	plan := configuration

	if !planDisabledShieldingPopID(&plan, configuration) {
		t.Fatal("disabled shielding plan was not normalized")
	}
	if plan.PopId.ValueInt64() != 0 {
		t.Fatalf("planned shielding pop_id = %d, want 0", plan.PopId.ValueInt64())
	}
}

func TestConfiguredSslCertificateChangeTreatsUnknownIDAsChange(t *testing.T) {
	ctx := context.Background()
	attributeTypes := resource_resource.SslCertificateValue{}.AttributeTypes(ctx)
	previous := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"id":     types.Int64Value(123),
		"status": types.StringValue("ready"),
		"type":   types.StringValue("own"),
	})
	configWithUnknownID := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"id":     types.Int64Unknown(),
		"status": types.StringUnknown(),
		"type":   types.StringValue("own"),
	})
	configWithKnownID := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"id":     types.Int64Value(123),
		"status": types.StringUnknown(),
		"type":   types.StringValue("own"),
	})

	if !configuredObjectChanged(configWithUnknownID, previous, "status") {
		t.Fatal("unknown configured certificate ID should prevent prior state from being reused")
	}
	if configuredObjectChanged(configWithKnownID, previous, "status") {
		t.Fatal("unknown computed certificate status should not count as a configuration change")
	}
}

func TestPreserveEdgeCacheDefaultDoesNotRestoreConflictingPriorValue(t *testing.T) {
	ctx := context.Background()
	attributeTypes := resource_resource.EdgeCacheSettingsValue{}.AttributeTypes(ctx)
	config := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"custom_values": types.MapNull(types.StringType),
		"default":       types.StringValue("10m"),
		"enabled":       types.BoolValue(true),
		"value":         types.StringNull(),
	})
	plan := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"custom_values": types.MapUnknown(types.StringType),
		"default":       types.StringValue("10m"),
		"enabled":       types.BoolValue(true),
		"value":         types.StringUnknown(),
	})
	previous := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"custom_values": types.MapNull(types.StringType),
		"default":       types.StringNull(),
		"enabled":       types.BoolValue(true),
		"value":         types.StringValue("10m"),
	})

	got, _, diags := preserveEdgeCacheSettingsValue(ctx, config, plan, previous)
	if diags.HasError() {
		t.Fatalf("preserveEdgeCacheSettingsValue returned diagnostics: %v", diags)
	}
	if !got.(types.Object).Attributes()["value"].IsUnknown() {
		t.Fatalf("planned edge cache value = %s, want unknown instead of prior value", got.(types.Object).Attributes()["value"])
	}
}

func TestPreserveEdgeCacheValueRetainsEmptyPriorDefault(t *testing.T) {
	ctx := context.Background()
	attributeTypes := resource_resource.EdgeCacheSettingsValue{}.AttributeTypes(ctx)
	config := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"custom_values": types.MapNull(types.StringType),
		"default":       types.StringNull(),
		"enabled":       types.BoolValue(true),
		"value":         types.StringValue("10m"),
	})
	plan := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"custom_values": types.MapUnknown(types.StringType),
		"default":       types.StringUnknown(),
		"enabled":       types.BoolValue(true),
		"value":         types.StringValue("10m"),
	})
	previous := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"custom_values": types.MapNull(types.StringType),
		"default":       types.StringValue(""),
		"enabled":       types.BoolValue(true),
		"value":         types.StringValue("old-value"),
	})

	got, _, diags := preserveEdgeCacheSettingsValue(ctx, config, plan, previous)
	if diags.HasError() {
		t.Fatalf("preserveEdgeCacheSettingsValue returned diagnostics: %v", diags)
	}
	if gotDefault := got.(types.Object).Attributes()["default"].(types.String); gotDefault.IsNull() || gotDefault.ValueString() != "" {
		t.Fatalf("planned edge cache default = %s, want prior empty value", gotDefault)
	}
}

func TestResourceStateUpgradeV0ConvertsTlsVersionsListToSet(t *testing.T) {
	ctx := context.Background()
	versionValues, diags := types.ListValueFrom(ctx, types.StringType, []string{"TLSv1.1", "TLSv1.2"})
	if diags.HasError() {
		t.Fatalf("types.ListValueFrom returned diagnostics: %v", diags)
	}
	priorOptions := priorOptionsValueV0(t, ctx, versionValues)

	priorSchema := resourceResourceSchemaV0(ctx)
	optionsSchema := priorSchema.Attributes["options"].(schema.SingleNestedAttribute)
	if _, ok := optionsSchema.Attributes["static_response_headers"]; ok {
		t.Fatal("prior schema version 0 unexpectedly contains static_response_headers")
	}
	tlsVersionsSchema := optionsSchema.Attributes["tls_versions"].(schema.SingleNestedAttribute)
	if _, ok := tlsVersionsSchema.Attributes["value"].(schema.ListAttribute); !ok {
		t.Fatal("prior schema version 0 does not describe tls_versions.value as a list")
	}
	priorState := tfsdk.State{Schema: priorSchema}
	priorModel := resourceResourceStateV0{
		Active:             types.BoolValue(true),
		Cname:              types.StringValue("cdn.example.com"),
		Id:                 types.Int64Value(123),
		Options:            priorOptions,
		OriginGroup:        types.Int64Value(456),
		OriginProtocol:     types.StringValue("HTTPS"),
		PresetApplied:      types.BoolValue(false),
		Region:             types.StringValue("RegionOne"),
		SecondaryHostnames: types.ListNull(types.StringType),
		Shielding:          resource_resource.NewShieldingValueNull(),
		SslCertificate:     resource_resource.NewSslCertificateValueNull(),
		Status:             types.StringValue("active"),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType,
			"update": types.StringType,
		})},
		VpEnabled: types.BoolValue(false),
	}
	if diags := priorState.Set(ctx, &priorModel); diags.HasError() {
		t.Fatalf("setting prior state returned diagnostics: %v", diags)
	}

	upgrader, ok := (&resourceResource{}).UpgradeState(ctx)[0]
	if !ok {
		t.Fatal("no state upgrader registered for schema version 0")
	}
	currentSchema := resource_resource.ResourceResourceSchema(ctx)
	currentSchema.Version = 1
	response := resource.UpgradeStateResponse{State: tfsdk.State{Schema: currentSchema}}
	upgrader.StateUpgrader(ctx, resource.UpgradeStateRequest{State: &priorState}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("state upgrader returned diagnostics: %v", response.Diagnostics)
	}

	var got resource_resource.ResourceModel
	if diags := response.State.Get(ctx, &got); diags.HasError() {
		t.Fatalf("reading upgraded state returned diagnostics: %v", diags)
	}
	if got.Cname.ValueString() != priorModel.Cname.ValueString() || got.Id.ValueInt64() != priorModel.Id.ValueInt64() {
		t.Fatalf("upgraded resource identity = %q/%d, want %q/%d", got.Cname.ValueString(), got.Id.ValueInt64(), priorModel.Cname.ValueString(), priorModel.Id.ValueInt64())
	}
	allowedHttpMethods, diags := resource_resource.AllowedHttpMethodsType{}.ValueFromObject(ctx, got.Options.AllowedHttpMethods)
	if diags.HasError() {
		t.Fatalf("decoding preserved HTTP methods returned diagnostics: %v", diags)
	}
	var gotMethods []string
	if diags := allowedHttpMethods.(resource_resource.AllowedHttpMethodsValue).Value.ElementsAs(ctx, &gotMethods, true); diags.HasError() {
		t.Fatalf("reading preserved HTTP methods returned diagnostics: %v", diags)
	}
	if len(gotMethods) != 1 || gotMethods[0] != "GET" {
		t.Fatalf("upgraded HTTP methods = %v, want [GET]", gotMethods)
	}
	tlsVersions, diags := resource_resource.TlsVersionsType{}.ValueFromObject(ctx, got.Options.TlsVersions)
	if diags.HasError() {
		t.Fatalf("decoding upgraded TLS option returned diagnostics: %v", diags)
	}
	wantVersions, diags := types.SetValueFrom(ctx, types.StringType, []string{"TLSv1.1", "TLSv1.2"})
	if diags.HasError() {
		t.Fatalf("building expected TLS version set returned diagnostics: %v", diags)
	}
	if !tlsVersions.(resource_resource.TlsVersionsValue).Value.Equal(wantVersions) {
		t.Fatalf("upgraded TLS versions = %v, want %v", tlsVersions.(resource_resource.TlsVersionsValue).Value, wantVersions)
	}

	var currentSchemaResponse resource.SchemaResponse
	(&resourceResource{}).Schema(ctx, resource.SchemaRequest{}, &currentSchemaResponse)
	if currentSchemaResponse.Schema.Version != 1 {
		t.Fatalf("current resource schema version = %d, want 1", currentSchemaResponse.Schema.Version)
	}
	currentOptionsSchema := currentSchemaResponse.Schema.Attributes["options"].(schema.SingleNestedAttribute)
	currentTlsSchema := currentOptionsSchema.Attributes["tls_versions"].(schema.SingleNestedAttribute)
	if _, ok := currentTlsSchema.Attributes["value"].(schema.SetAttribute); !ok {
		t.Fatal("current schema does not describe tls_versions.value as a set")
	}
}

func TestUpgradeTLSVersionsListToSetPreservesNullAndUnknown(t *testing.T) {
	for _, tc := range []struct {
		name  string
		prior types.List
		check func(types.Set) bool
	}{
		{name: "null", prior: types.ListNull(types.StringType), check: func(value types.Set) bool { return value.IsNull() }},
		{name: "unknown", prior: types.ListUnknown(types.StringType), check: func(value types.Set) bool { return value.IsUnknown() }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, diags := upgradeTLSVersionsListToSet(tc.prior)
			if diags.HasError() {
				t.Fatalf("upgradeTLSVersionsListToSet returned diagnostics: %v", diags)
			}
			if !tc.check(got) {
				t.Fatalf("upgraded set = %v, want %s", got, tc.name)
			}
		})
	}
}

func priorOptionsValueV0(t *testing.T, ctx context.Context, tlsVersions types.List) types.Object {
	t.Helper()

	options, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		AllowedHttpMethods: &resources.ResourceOptionsAllowedHttpMethodsOption{
			Enabled: true,
			Value:   []resources.ResourceAllowedHttpMethod{resources.ResourceAllowedHttpMethodGet},
		},
	})
	if diags.HasError() {
		t.Fatalf("creating options fixture returned diagnostics: %v", diags)
	}
	optionsObject, diags := options.ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("converting options fixture returned diagnostics: %v", diags)
	}

	priorTlsTypes := map[string]attr.Type{
		"enabled": types.BoolType,
		"value":   types.ListType{ElemType: types.StringType},
	}
	priorTls, diags := types.ObjectValue(priorTlsTypes, map[string]attr.Value{
		"enabled": types.BoolValue(true),
		"value":   tlsVersions,
	})
	if diags.HasError() {
		t.Fatalf("creating prior TLS option returned diagnostics: %v", diags)
	}

	attributes := optionsObject.Attributes()
	delete(attributes, "static_response_headers")
	attributes["tls_versions"] = priorTls
	priorOptionTypes := resource_resource.OptionsValue{}.AttributeTypes(ctx)
	delete(priorOptionTypes, "static_response_headers")
	priorOptionTypes["tls_versions"] = types.ObjectType{AttrTypes: priorTlsTypes}
	priorOptions, diags := types.ObjectValue(priorOptionTypes, attributes)
	if diags.HasError() {
		t.Fatalf("creating prior options object returned diagnostics: %v", diags)
	}
	return priorOptions
}
