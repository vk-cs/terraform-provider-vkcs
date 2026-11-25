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

func TestValidateStaticHeadersConflict(t *testing.T) {
	ctx := context.Background()
	legacy := types.ObjectValueMust(resource_resource.StaticHeadersValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"enabled": types.BoolValue(false),
		"value":   types.MapValueMust(types.StringType, map[string]attr.Value{}),
	})
	current := types.ObjectValueMust(resource_resource.StaticResponseHeadersValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"enabled": types.BoolValue(false),
		"value":   types.ListValueMust(resource_resource.ValueValue{}.Type(ctx), []attr.Value{}),
	})

	diags := validateStaticHeadersConflict(legacy, current)
	if !diags.HasError() {
		t.Fatal("expected a diagnostic when both options are configured, even when disabled")
	}
}

func TestValidateStaticHeadersConflictAllowsOneOptionOrUnknown(t *testing.T) {
	ctx := context.Background()
	legacyTypes := resource_resource.StaticHeadersValue{}.AttributeTypes(ctx)
	currentTypes := resource_resource.StaticResponseHeadersValue{}.AttributeTypes(ctx)

	for _, tc := range []struct {
		name    string
		legacy  types.Object
		current types.Object
	}{
		{
			name:    "legacy only",
			legacy:  types.ObjectValueMust(legacyTypes, map[string]attr.Value{"enabled": types.BoolValue(true), "value": types.MapValueMust(types.StringType, map[string]attr.Value{})}),
			current: types.ObjectNull(currentTypes),
		},
		{
			name:    "new option unknown",
			legacy:  types.ObjectNull(legacyTypes),
			current: types.ObjectUnknown(currentTypes),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if diags := validateStaticHeadersConflict(tc.legacy, tc.current); diags.HasError() {
				t.Fatalf("unexpected conflict diagnostic: %v", diags)
			}
		})
	}
}

func TestValidateStaticHeadersConflictRejectsBothUnknownOptions(t *testing.T) {
	ctx := context.Background()
	legacy := types.ObjectUnknown(resource_resource.StaticHeadersValue{}.AttributeTypes(ctx))
	current := types.ObjectUnknown(resource_resource.StaticResponseHeadersValue{}.AttributeTypes(ctx))

	if diags := validateStaticHeadersConflict(legacy, current); !diags.HasError() {
		t.Fatal("expected a conflict diagnostic when both configured options are unknown")
	}
}

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
	diags = selectConfiguredResponseHeaderOption(ctx, &plan, newConfig, &legacyState)
	if diags.HasError() {
		t.Fatalf("selectConfiguredResponseHeaderOption returned diagnostics: %v", diags)
	}
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
	if diags := selectConfiguredResponseHeaderOption(ctx, &plan, legacyOnly, nil); diags.HasError() {
		t.Fatalf("selectConfiguredResponseHeaderOption returned diagnostics: %v", diags)
	}
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
	if diags := selectConfiguredResponseHeaderOption(ctx, &noOpPlan, legacyOnly, &readBack); diags.HasError() {
		t.Fatalf("selectConfiguredResponseHeaderOption on no-op returned diagnostics: %v", diags)
	}
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
	diags = selectConfiguredResponseHeaderOption(ctx, &plan, legacyConfig, &newState)
	if diags.HasError() {
		t.Fatalf("selectConfiguredResponseHeaderOption returned diagnostics: %v", diags)
	}
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

	requestOptions, diags := responseHeaderOptionsForAPIRequest(ctx, plan, legacyConfig)
	if diags.HasError() {
		t.Fatalf("responseHeaderOptionsForAPIRequest returned diagnostics: %v", diags)
	}
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

	requestOptions, diags := responseHeaderOptionsForAPIRequest(ctx, plan, newConfig)
	if diags.HasError() {
		t.Fatalf("responseHeaderOptionsForAPIRequest returned diagnostics: %v", diags)
	}
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

func TestSelectConfiguredResponseHeaderOptionRejectsAmbiguousUnknownOptions(t *testing.T) {
	ctx := context.Background()
	plan, diags := (resource_resource.OptionsValue{}).FromResourceOptions(ctx, &resources.ResourceOptions{
		StaticHeaders: &resources.ResourceOptionsStringMapOption{
			Enabled: true,
			Value:   map[string]string{"X-Legacy": "legacy"},
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
	unknownConfig := resource_resource.NewOptionsValueUnknown()

	diags = selectConfiguredResponseHeaderOption(ctx, &plan, unknownConfig, nil)
	if !diags.HasError() {
		t.Fatal("expected a conflict diagnostic for ambiguous planned header options")
	}
}

func TestValidateStaticResponseHeadersRequiresAlways(t *testing.T) {
	ctx := context.Background()
	header, diags := resource_resource.NewValueValue(resource_resource.ValueValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"name":   types.StringValue("X-Test"),
		"value":  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value")}),
		"always": types.BoolNull(),
	})
	if diags.HasError() {
		t.Fatalf("NewValueValue returned diagnostics: %v", diags)
	}
	values := types.ListValueMust(resource_resource.ValueValue{}.Type(ctx), []attr.Value{header})
	option := types.ObjectValueMust(resource_resource.StaticResponseHeadersValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"enabled": types.BoolValue(true),
		"value":   values,
	})

	diags = validateStaticResponseHeaders(ctx, option)
	if !diags.HasError() {
		t.Fatal("expected a diagnostic when always is not specified")
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
