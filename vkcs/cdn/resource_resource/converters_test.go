package resource_resource

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/resources"
)

func TestStaticResponseHeadersToResourceOptions(t *testing.T) {
	ctx := context.Background()
	option := newStaticResponseHeadersOption(t, ctx, true, []ValueValue{
		newHeaderValue(t, ctx, "X-Multi", []string{"first", "second"}, false),
		newHeaderValue(t, ctx, "X-Always", []string{"third"}, true),
	})
	options := optionsWithAttribute(t, ctx, "static_response_headers", option)

	apiOptions, diags := options.ToResourceOptions(ctx)
	if diags.HasError() {
		t.Fatalf("ToResourceOptions returned diagnostics: %v", diags)
	}
	if apiOptions.StaticResponseHeaders == nil {
		t.Fatal("StaticResponseHeaders is nil")
	}
	if !apiOptions.StaticResponseHeaders.Enabled {
		t.Fatal("StaticResponseHeaders.Enabled = false, want true")
	}
	want := []resources.StaticResponseHeaderValue{
		{Name: "X-Multi", Value: []string{"first", "second"}, Always: false},
		{Name: "X-Always", Value: []string{"third"}, Always: true},
	}
	if len(apiOptions.StaticResponseHeaders.Value) != len(want) {
		t.Fatalf("got %d headers, want %d", len(apiOptions.StaticResponseHeaders.Value), len(want))
	}
	for i := range want {
		if got := apiOptions.StaticResponseHeaders.Value[i]; got.Name != want[i].Name || got.Always != want[i].Always || !equalStrings(got.Value, want[i].Value) {
			t.Errorf("header %d = %#v, want %#v", i, got, want[i])
		}
	}

	body, err := json.Marshal(apiOptions)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}
	var request map[string]json.RawMessage
	if err := json.Unmarshal(body, &request); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}
	if _, ok := request["static_response_headers"]; !ok {
		t.Fatalf("API request is missing static_response_headers: %s", body)
	}
	if _, ok := request["staticResponseHeaders"]; ok {
		t.Fatalf("API request used the wrong field name: %s", body)
	}
}

func TestStaticResponseHeadersFromResourceOptions(t *testing.T) {
	ctx := context.Background()
	apiOptions := &resources.ResourceOptions{
		StaticResponseHeaders: &resources.ResourceOptionsStaticResponseHeadersOption{
			Enabled: true,
			Value: []resources.StaticResponseHeaderValue{
				{Name: "X-Multi", Value: []string{"first", "second"}, Always: false},
				{Name: "X-Always", Value: []string{"third"}, Always: true},
			},
		},
	}

	options, diags := OptionsValue{}.FromResourceOptions(ctx, apiOptions)
	if diags.HasError() {
		t.Fatalf("FromResourceOptions returned diagnostics: %v", diags)
	}
	got := StaticResponseHeadersTypeValue(t, ctx, options.StaticResponseHeaders)
	if !got.Enabled.ValueBool() {
		t.Fatal("StaticResponseHeaders.Enabled = false, want true")
	}
	var headers []ValueValue
	if diags := got.Value.ElementsAs(ctx, &headers, true); diags.HasError() {
		t.Fatalf("ElementsAs returned diagnostics: %v", diags)
	}
	if len(headers) != 2 {
		t.Fatalf("got %d headers, want 2", len(headers))
	}
	for i, want := range []struct {
		name   string
		values []string
		always bool
	}{
		{name: "X-Multi", values: []string{"first", "second"}, always: false},
		{name: "X-Always", values: []string{"third"}, always: true},
	} {
		var values []string
		if diags := headers[i].Value.ElementsAs(ctx, &values, true); diags.HasError() {
			t.Fatalf("header %d ElementsAs returned diagnostics: %v", i, diags)
		}
		if headers[i].Name.ValueString() != want.name || headers[i].Always.ValueBool() != want.always || !equalStrings(values, want.values) {
			t.Errorf("header %d = name %q, values %v, always %t; want %q, %v, %t", i, headers[i].Name.ValueString(), values, headers[i].Always.ValueBool(), want.name, want.values, want.always)
		}
	}
}

func TestStaticResponseHeadersNullUnknownAndDisabled(t *testing.T) {
	ctx := context.Background()
	attributeTypes := StaticResponseHeadersValue{}.AttributeTypes(ctx)

	for _, tc := range []struct {
		name   string
		option types.Object
	}{
		{name: "null", option: types.ObjectNull(attributeTypes)},
		{name: "unknown", option: types.ObjectUnknown(attributeTypes)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := optionsWithAttribute(t, ctx, "static_response_headers", tc.option)
			apiOptions, diags := options.ToResourceOptions(ctx)
			if diags.HasError() {
				t.Fatalf("ToResourceOptions returned diagnostics: %v", diags)
			}
			if apiOptions.StaticResponseHeaders != nil {
				t.Fatalf("StaticResponseHeaders = %#v, want nil", apiOptions.StaticResponseHeaders)
			}
		})
	}

	disabled := newStaticResponseHeadersOption(t, ctx, false, nil)
	optionValue := StaticResponseHeadersTypeValue(t, ctx, disabled)
	optionValue.Value = types.ListNull(ValueValue{}.Type(ctx))
	disabled, diags := optionValue.ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("ToObjectValue returned diagnostics: %v", diags)
	}
	apiOptions, diags := optionsWithAttribute(t, ctx, "static_response_headers", disabled).ToResourceOptions(ctx)
	if diags.HasError() {
		t.Fatalf("ToResourceOptions returned diagnostics for disabled option: %v", diags)
	}
	if apiOptions.StaticResponseHeaders == nil || apiOptions.StaticResponseHeaders.Enabled || apiOptions.StaticResponseHeaders.Value == nil || len(apiOptions.StaticResponseHeaders.Value) != 0 {
		t.Fatalf("disabled API option = %#v, want enabled=false and an empty value list", apiOptions.StaticResponseHeaders)
	}
}

func TestStaticResponseHeadersRejectsMissingRequiredAPIFields(t *testing.T) {
	ctx := context.Background()
	valueList, diags := types.ListValue(ValueValue{}.Type(ctx), []attr.Value{
		ValueValue{
			Name:   types.StringValue("X-Missing-Always"),
			Value:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("v")}),
			Always: types.BoolNull(),
			state:  attr.ValueStateKnown,
		},
	})
	if diags.HasError() {
		t.Fatalf("types.ListValue returned diagnostics: %v", diags)
	}
	option, diags := (StaticResponseHeadersValue{
		Enabled: types.BoolValue(true),
		Value:   valueList,
		state:   attr.ValueStateKnown,
	}).ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("ToObjectValue returned diagnostics: %v", diags)
	}

	apiOptions, diags := optionsWithAttribute(t, ctx, "static_response_headers", option).ToResourceOptions(ctx)
	if !diags.HasError() {
		t.Fatal("ToResourceOptions returned no diagnostic for missing always")
	}
	if apiOptions != nil {
		t.Fatalf("ToResourceOptions returned API options %#v despite invalid input", apiOptions)
	}
}

func TestStaticResponseHeadersUnknownValueList(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name      string
		enabled   bool
		wantError bool
	}{
		{name: "enabled", enabled: true, wantError: true},
		{name: "disabled", enabled: false, wantError: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			option := types.ObjectValueMust(StaticResponseHeadersValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"enabled": types.BoolValue(tc.enabled),
				"value":   types.ListUnknown(ValueValue{}.Type(ctx)),
			})

			apiOptions, diags := optionsWithAttribute(t, ctx, "static_response_headers", option).ToResourceOptions(ctx)
			if diags.HasError() != tc.wantError {
				t.Fatalf("HasError() = %t, want %t; diagnostics: %v", diags.HasError(), tc.wantError, diags)
			}
			if tc.wantError && apiOptions != nil {
				t.Fatalf("ToResourceOptions returned API options %#v despite unknown enabled values", apiOptions)
			}
			if !tc.wantError && (apiOptions.StaticResponseHeaders == nil || apiOptions.StaticResponseHeaders.Enabled || apiOptions.StaticResponseHeaders.Value == nil || len(apiOptions.StaticResponseHeaders.Value) != 0) {
				t.Fatalf("disabled API option = %#v, want an empty list", apiOptions.StaticResponseHeaders)
			}
		})
	}
}

func TestStaticHeadersRemainsSupported(t *testing.T) {
	ctx := context.Background()
	legacy, diags := (StaticHeadersValue{
		Enabled: types.BoolValue(true),
		Value:   types.MapValueMust(types.StringType, map[string]attr.Value{"X-Legacy": types.StringValue("legacy")}),
		state:   attr.ValueStateKnown,
	}).ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("legacy ToObjectValue returned diagnostics: %v", diags)
	}

	options := optionsWithAttribute(t, ctx, "static_headers", legacy)
	apiOptions, diags := options.ToResourceOptions(ctx)
	if diags.HasError() {
		t.Fatalf("legacy ToResourceOptions returned diagnostics: %v", diags)
	}
	if apiOptions.StaticHeaders == nil || apiOptions.StaticHeaders.Value["X-Legacy"] != "legacy" {
		t.Fatalf("legacy option = %#v, want original staticHeaders API payload", apiOptions.StaticHeaders)
	}
}

func TestTlsVersionsUseCanonicalAPIOrderInBothDirections(t *testing.T) {
	ctx := context.Background()
	configuredValues, diags := types.ListValueFrom(ctx, types.StringType, []string{"TLSv1.1", "TLSv1.2"})
	if diags.HasError() {
		t.Fatalf("types.ListValueFrom returned diagnostics: %v", diags)
	}
	tlsVersions, diags := (TlsVersionsValue{
		Value:   configuredValues,
		Enabled: types.BoolValue(true),
		state:   attr.ValueStateKnown,
	}).ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("TlsVersionsValue.ToObjectValue returned diagnostics: %v", diags)
	}
	options := optionsWithAttribute(t, ctx, "tls_versions", tlsVersions)
	apiOptions, diags := options.ToResourceOptions(ctx)
	if diags.HasError() {
		t.Fatalf("ToResourceOptions returned diagnostics: %v", diags)
	}
	if got := apiTlsVersionStrings(apiOptions.TlsVersions.Value); !equalStrings(got, []string{"TLSv1.2", "TLSv1.1"}) {
		t.Fatalf("request TLS versions = %v, want canonical API order [TLSv1.2 TLSv1.1]", got)
	}

	readOptions, diags := OptionsValue{}.FromResourceOptions(ctx, &resources.ResourceOptions{
		TlsVersions: &resources.ResourceOptionsTlsVersionsOption{
			Enabled: true,
			Value: []resources.ResourceTlsVersion{
				resources.ResourceTlsVersionTLSv1_1,
				resources.ResourceTlsVersionTLSv1_2,
			},
		},
	})
	if diags.HasError() {
		t.Fatalf("FromResourceOptions returned diagnostics: %v", diags)
	}
	tlsObject, diags := TlsVersionsType{}.ValueFromObject(ctx, readOptions.TlsVersions)
	if diags.HasError() {
		t.Fatalf("TlsVersionsType.ValueFromObject returned diagnostics: %v", diags)
	}
	var readValues []string
	if diags := tlsObject.(TlsVersionsValue).Value.ElementsAs(ctx, &readValues, true); diags.HasError() {
		t.Fatalf("reading TLS versions returned diagnostics: %v", diags)
	}
	if !equalStrings(readValues, []string{"TLSv1.2", "TLSv1.1"}) {
		t.Fatalf("state TLS versions = %v, want canonical API order [TLSv1.2 TLSv1.1]", readValues)
	}
}

func newStaticResponseHeadersOption(t *testing.T, ctx context.Context, enabled bool, headers []ValueValue) types.Object {
	t.Helper()
	values := make([]attr.Value, 0, len(headers))
	for _, header := range headers {
		values = append(values, header)
	}
	value, diags := types.ListValue(ValueValue{}.Type(ctx), values)
	if diags.HasError() {
		t.Fatalf("types.ListValue returned diagnostics: %v", diags)
	}
	option, diags := (StaticResponseHeadersValue{
		Enabled: types.BoolValue(enabled),
		Value:   value,
		state:   attr.ValueStateKnown,
	}).ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("ToObjectValue returned diagnostics: %v", diags)
	}
	return option
}

func newHeaderValue(t *testing.T, ctx context.Context, name string, values []string, always bool) ValueValue {
	t.Helper()
	value, diags := types.ListValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		t.Fatalf("types.ListValueFrom returned diagnostics: %v", diags)
	}
	return ValueValue{
		Name:   types.StringValue(name),
		Value:  value,
		Always: types.BoolValue(always),
		state:  attr.ValueStateKnown,
	}
}

func optionsWithAttribute(t *testing.T, ctx context.Context, name string, value types.Object) OptionsValue {
	t.Helper()
	options, diags := OptionsValue{}.FromResourceOptions(ctx, &resources.ResourceOptions{})
	if diags.HasError() {
		t.Fatalf("OptionsValue.FromResourceOptions returned diagnostics: %v", diags)
	}
	optionObject, diags := options.ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("OptionsValue.ToObjectValue returned diagnostics: %v", diags)
	}
	attributes := optionObject.Attributes()
	attributes[name] = value
	object, diags := types.ObjectValue(OptionsValue{}.AttributeTypes(ctx), attributes)
	if diags.HasError() {
		t.Fatalf("types.ObjectValue returned diagnostics: %v", diags)
	}
	result, diags := OptionsType{}.ValueFromObject(ctx, object)
	if diags.HasError() {
		t.Fatalf("OptionsType.ValueFromObject returned diagnostics: %v", diags)
	}
	return result.(OptionsValue)
}

func StaticResponseHeadersTypeValue(t *testing.T, ctx context.Context, object types.Object) StaticResponseHeadersValue {
	t.Helper()
	value, diags := StaticResponseHeadersType{}.ValueFromObject(ctx, object)
	if diags.HasError() {
		t.Fatalf("StaticResponseHeadersType.ValueFromObject returned diagnostics: %v", diags)
	}
	return value.(StaticResponseHeadersValue)
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func apiTlsVersionStrings(values []resources.ResourceTlsVersion) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = string(value)
	}
	return result
}
