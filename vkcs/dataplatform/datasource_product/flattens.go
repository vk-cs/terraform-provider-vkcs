package datasource_product

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/products"
)

func FlattenConfigs(ctx context.Context, o *products.ProductConfig) (ConfigsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o == nil {
		return NewConfigsValueNull(), nil
	}

	connections, d := FlattenConfigsConnections(ctx, o.Connections)
	diags.Append(d...)
	if diags.HasError() {
		return NewConfigsValueUnknown(), diags
	}

	settings, d := FlattenConfigsSettings(ctx, o.Settings)
	diags.Append(d...)
	if diags.HasError() {
		return NewConfigsValueUnknown(), diags
	}

	userRoles, d := FlattenConfigsUserRoles(ctx, o.UserRoles)
	diags.Append(d...)
	if diags.HasError() {
		return NewConfigsValueUnknown(), diags
	}

	crontabs, d := FlattenConfigsCrontabs(ctx, o.Crontabs)
	diags.Append(d...)
	if diags.HasError() {
		return NewConfigsValueUnknown(), diags
	}

	configsV := ConfigsValue{
		Connections: connections,
		Settings:    settings,
		UserRoles:   userRoles,
		Crontabs:    crontabs,
		state:       attr.ValueStateKnown,
	}

	return configsV, nil
}

func FlattenConfigsConnections(ctx context.Context, o []products.ProductConfigConnection) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(ConfigsConnectionsValue{}.Type(ctx)), nil
	}

	connectionsV := make([]attr.Value, len(o))
	for i, c := range o {
		settings, d := FlattenConfigsConnectionsSettings(ctx, c.Settings)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListUnknown(ConfigsConnectionsValue{}.Type(ctx)), diags
		}
		connectionsV[i] = ConfigsConnectionsValue{
			IsRequired:    types.BoolValue(c.IsRequired),
			Plug:          types.StringValue(c.Plug),
			Position:      types.Int64Value(int64(c.Position)),
			RequiredGroup: types.StringValue(c.RequiredGroup),
			Settings:      settings,
			state:         attr.ValueStateKnown,
		}
	}

	result, d := types.ListValue(ConfigsConnectionsValue{}.Type(ctx), connectionsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(ConfigsConnectionsValue{}.Type(ctx)), diags
	}
	return result, nil
}

func FlattenConfigsConnectionsSettings(ctx context.Context, o []products.ProductConfigConnectionSetting) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(ConfigsConnectionsSettingsValue{}.Type(ctx)), nil
	}

	settingsV := make([]attr.Value, len(o))
	for i, s := range o {
		stringVariation, d := FlattenConfigsConnectionsSettingsStringVariation(ctx, s.StringVariation)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListUnknown(ConfigsConnectionsSettingsValue{}.Type(ctx)), diags
		}

		settingsV[i] = ConfigsConnectionsSettingsValue{
			Alias:           types.StringValue(s.Alias),
			DefaultValue:    types.StringValue(s.DefaultValue),
			IsRequire:       types.BoolValue(s.IsRequired),
			IsSensitive:     types.BoolValue(s.IsSensitive),
			Regexp:          types.StringValue(s.RegExp),
			StringVariation: stringVariation,
			state:           attr.ValueStateKnown,
		}
	}
	result, d := types.ListValue(ConfigsConnectionsSettingsValue{}.Type(ctx), settingsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(ConfigsConnectionsSettingsValue{}.Type(ctx)), diags
	}
	return result, nil
}

func FlattenConfigsConnectionsSettingsStringVariation(ctx context.Context, o []string) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(types.StringType), nil
	}

	settingsV := make([]attr.Value, len(o))
	for i, s := range o {
		settingsV[i] = types.StringValue(s)
	}
	result, d := types.ListValue(types.StringType, settingsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(types.StringType), diags
	}
	return result, nil
}

func FlattenConfigsSettings(ctx context.Context, o []products.ProductConfigSetting) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(ConfigsSettingsValue{}.Type(ctx)), nil
	}

	settingsV := make([]attr.Value, len(o))
	for i, s := range o {
		stringVariation, d := FlattenConfigsSettingsStringVariation(ctx, s.StringVariation)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListUnknown(ConfigsSettingsValue{}.Type(ctx)), diags
		}

		settingsV[i] = ConfigsSettingsValue{
			Alias:           types.StringValue(s.Alias),
			DefaultValue:    types.StringValue(s.DefaultValue),
			IsRequire:       types.BoolValue(s.IsRequired),
			IsSensitive:     types.BoolValue(s.IsSensitive),
			Regexp:          types.StringValue(s.RegExp),
			StringVariation: stringVariation,
			state:           attr.ValueStateKnown,
		}
	}
	result, d := types.ListValue(ConfigsSettingsValue{}.Type(ctx), settingsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(ConfigsSettingsValue{}.Type(ctx)), diags
	}
	return result, nil
}

func FlattenConfigsSettingsStringVariation(ctx context.Context, o []string) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(types.StringType), nil
	}

	settingsV := make([]attr.Value, len(o))
	for i, s := range o {
		settingsV[i] = types.StringValue(s)
	}
	result, d := types.ListValue(types.StringType, settingsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(types.StringType), diags
	}
	return result, nil
}

func FlattenConfigsUserRoles(ctx context.Context, o []products.ProductConfigUserRole) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(ConfigsUserRolesValue{}.Type(ctx)), nil
	}

	userRolesV := make([]attr.Value, len(o))
	for i, s := range o {
		userRolesV[i] = ConfigsUserRolesValue{
			Name:  types.StringValue(s.Name),
			state: attr.ValueStateKnown,
		}
	}
	result, d := types.ListValue(ConfigsUserRolesValue{}.Type(ctx), userRolesV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(ConfigsUserRolesValue{}.Type(ctx)), diags
	}
	return result, nil
}

func FlattenConfigsCrontabs(ctx context.Context, o []products.ProductConfigCrontabs) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(ConfigsCrontabsValue{}.Type(ctx)), nil
	}

	crontabsV := make([]attr.Value, len(o))
	for i, c := range o {
		settings, d := FlattenConfigsCrontabsSettings(ctx, c.Settings)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListUnknown(ConfigsCrontabsValue{}.Type(ctx)), diags
		}
		crontabsV[i] = ConfigsCrontabsValue{
			Name:     types.StringValue(c.Name),
			Required: types.BoolValue(c.Required),
			Start:    types.StringValue(c.Start),
			Settings: settings,
			state:    attr.ValueStateKnown,
		}
	}
	result, d := types.ListValue(ConfigsCrontabsValue{}.Type(ctx), crontabsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(ConfigsCrontabsValue{}.Type(ctx)), diags
	}
	return result, nil
}

func FlattenConfigsCrontabsSettings(ctx context.Context, o []products.ProductConfigCrontabsSettings) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(ConfigsCrontabsSettingsValue{}.Type(ctx)), nil
	}

	settingsV := make([]attr.Value, len(o))
	for i, s := range o {
		stringVariation, d := FlattenConfigsCrontabsSettingsStringVariation(ctx, s.StringVariation)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListUnknown(ConfigsCrontabsSettingsValue{}.Type(ctx)), diags
		}

		settingsV[i] = ConfigsCrontabsSettingsValue{
			Alias:           types.StringValue(s.Alias),
			DefaultValue:    types.StringValue(s.DefaultValue),
			IsRequire:       types.BoolValue(s.IsRequired),
			IsSensitive:     types.BoolValue(s.IsSensitive),
			Regexp:          types.StringValue(s.RegExp),
			StringVariation: stringVariation,
			state:           attr.ValueStateKnown,
		}
	}
	result, d := types.ListValue(ConfigsCrontabsSettingsValue{}.Type(ctx), settingsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(ConfigsCrontabsSettingsValue{}.Type(ctx)), diags
	}
	return result, nil
}

func FlattenConfigsCrontabsSettingsStringVariation(ctx context.Context, o []string) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(types.StringType), nil
	}

	settingsV := make([]attr.Value, len(o))
	for i, s := range o {
		settingsV[i] = types.StringValue(s)
	}
	result, d := types.ListValue(types.StringType, settingsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(types.StringType), diags
	}
	return result, nil
}
