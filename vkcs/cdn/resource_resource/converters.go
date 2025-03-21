package resource_resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/resources"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func (m *ResourceModel) UpdateFromResource(ctx context.Context, resource *resources.Resource) diag.Diagnostics {
	var diags diag.Diagnostics

	if resource == nil {
		return diags
	}

	m.Id = types.Int64Value(int64(resource.ID))
	m.Active = types.BoolValue(resource.Active)
	m.Cname = types.StringValue(resource.CNAME)

	options, d := OptionsValue{}.FromResourceOptions(ctx, &resource.Options)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.Options = options

	m.OriginGroup = types.Int64Value(int64(resource.OriginGroup))
	m.OriginProtocol = types.StringValue(string(resource.OriginProtocol))
	m.PresetApplied = types.BoolValue(resource.PresetApplied)

	secondaryHostnamesV, d := types.ListValueFrom(ctx, types.StringType, resource.SecondaryHostnames)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	m.SecondaryHostnames = secondaryHostnamesV
	m.SslCertificate = SslCertificateValue{}.FromSslOpts(ctx, &SslOpts{
		LeEnabled: resource.SSLLeEnabled,
		Enabled:   resource.SSLEnabled,
		Data:      resource.SSLData,
	})
	m.Status = types.StringValue(resource.Status)
	m.VpEnabled = types.BoolValue(resource.VPEnabled)

	return diags
}

func (v OptionsValue) ToResourceOptions(ctx context.Context) (*resources.ResourceOptions, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() || v.IsUnknown() {
		return nil, diags
	}

	result := &resources.ResourceOptions{}

	if o := v.AllowedHttpMethods; !o.IsUnknown() && !o.IsNull() {
		allowedHttpMethodsObjV, d := AllowedHttpMethodsType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		allowedHttpMethods := allowedHttpMethodsObjV.(AllowedHttpMethodsValue)
		value := make([]resources.ResourceAllowedHttpMethod, 0, len(allowedHttpMethods.Value.Elements()))
		if optV := allowedHttpMethods.Value; !optV.IsUnknown() && !optV.IsNull() {
			diags.Append(optV.ElementsAs(ctx, &value, true)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.AllowedHttpMethods = &resources.ResourceOptionsAllowedHttpMethodsOption{
			Enabled: allowedHttpMethods.Enabled.ValueBool(),
			Value:   value,
		}
	}

	if o := v.BrotliCompression; !o.IsUnknown() && !o.IsNull() {
		brotliCompressionObjV, d := BrotliCompressionType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		brotliCompression := brotliCompressionObjV.(BrotliCompressionValue)
		value := make([]string, 0, len(brotliCompression.Value.Elements()))
		if optV := brotliCompression.Value; !optV.IsUnknown() && !optV.IsNull() {
			diags.Append(optV.ElementsAs(ctx, &value, true)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.BrotliCompression = &resources.ResourceOptionsStringListOption{
			Enabled: brotliCompression.Enabled.ValueBool(),
			Value:   value,
		}
	}

	if o := v.BrowserCacheSettings; !o.IsUnknown() && !o.IsNull() {
		browserCacheSettingsObjV, d := BrowserCacheSettingsType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		browserCacheSettings := browserCacheSettingsObjV.(BrowserCacheSettingsValue)

		result.BrowserCacheSettings = &resources.ResourceOptionsStringOption{
			Enabled: browserCacheSettings.Enabled.ValueBool(),
			Value:   browserCacheSettings.Value.ValueString(),
		}
	}

	if o := v.Cors; !o.IsUnknown() && !o.IsNull() {
		corsObjV, d := CorsType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		cors := corsObjV.(CorsValue)
		value := make([]string, 0, len(cors.Value.Elements()))
		if optV := cors.Value; !optV.IsUnknown() && !optV.IsNull() {
			diags.Append(cors.Value.ElementsAs(ctx, &value, true)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.CORS = &resources.ResourceOptionsStringListOption{
			Enabled: cors.Enabled.ValueBool(),
			Value:   value,
		}
	}

	if o := v.CountryAcl; !o.IsUnknown() && !o.IsNull() {
		countryAclObjV, d := CountryAclType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		countryAcl := countryAclObjV.(CountryAclValue)
		exceptedValues := make([]string, 0, len(countryAcl.ExceptedValues.Elements()))
		if !countryAcl.ExceptedValues.IsUnknown() && !countryAcl.ExceptedValues.IsNull() {
			diags.Append(countryAcl.ExceptedValues.ElementsAs(ctx, &exceptedValues, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.CountryACL = &resources.ResourceOptionsACLOption{
			Enabled:        countryAcl.Enabled.ValueBool(),
			ExceptedValues: exceptedValues,
			PolicyType:     resources.ResourceACLPolicyType(countryAcl.PolicyType.ValueString()),
		}
	}

	if o := v.EdgeCacheSettings; !o.IsUnknown() && !o.IsNull() {
		edgeCacheSettingsObjV, d := EdgeCacheSettingsType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		edgeCacheSettings := edgeCacheSettingsObjV.(EdgeCacheSettingsValue)

		var customValues map[string]string
		if !edgeCacheSettings.CustomValues.IsUnknown() && !edgeCacheSettings.CustomValues.IsNull() {
			customValues = make(map[string]string, len(edgeCacheSettings.CustomValues.Elements()))
			diags.Append(edgeCacheSettings.CustomValues.ElementsAs(ctx, &customValues, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.EdgeCacheSettings = &resources.ResourceOptionsEdgeCacheSettingsOption{
			CustomValues: &customValues,
			Default:      edgeCacheSettings.Default.ValueString(),
			Enabled:      edgeCacheSettings.Enabled.ValueBool(),
			Value:        edgeCacheSettings.Value.ValueString(),
		}
	}

	result.FetchCompressed = expandBoolOption(v.FetchCompressed)

	if o := v.ForceReturn; !o.IsUnknown() && !o.IsNull() {
		forceReturnObjV, d := ForceReturnType{}.ValueFromObject(ctx, v.ForceReturn)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		forceReturn := forceReturnObjV.(ForceReturnValue)
		result.ForceReturn = &resources.ResourceOptionsForceReturnOption{
			Body:    forceReturn.Body.ValueString(),
			Code:    int(forceReturn.Code.ValueInt64()),
			Enabled: forceReturn.Enabled.ValueBool(),
		}
	}

	result.ForwardHostHeader = expandBoolOption(v.ForwardHostHeader)
	result.GzipOn = expandBoolOption(v.GzipOn)

	if o := v.HostHeader; !o.IsUnknown() && !o.IsNull() {
		hostHeaderObjV, d := HostHeaderType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		hostHeader := hostHeaderObjV.(HostHeaderValue)
		result.HostHeader = &resources.ResourceOptionsStringOption{
			Enabled: hostHeader.Enabled.ValueBool(),
			Value:   hostHeader.Value.ValueString(),
		}
	}

	result.IgnoreCookie = expandBoolOption(v.IgnoreCookie)
	result.IgnoreQueryString = expandBoolOption(v.IgnoreQueryString)

	if o := v.IpAddressAcl; !o.IsUnknown() && !o.IsNull() {
		ipAddressAclObjV, d := IpAddressAclType{}.ValueFromObject(ctx, v.IpAddressAcl)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		ipAddressAcl := ipAddressAclObjV.(IpAddressAclValue)
		exceptedValues := make([]string, 0, len(ipAddressAcl.ExceptedValues.Elements()))
		if !ipAddressAcl.ExceptedValues.IsUnknown() && !ipAddressAcl.ExceptedValues.IsNull() {
			diags.Append(ipAddressAcl.ExceptedValues.ElementsAs(ctx, &exceptedValues, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.IpAddressACL = &resources.ResourceOptionsACLOption{
			Enabled:        ipAddressAcl.Enabled.ValueBool(),
			ExceptedValues: exceptedValues,
			PolicyType:     resources.ResourceACLPolicyType(ipAddressAcl.PolicyType.ValueString()),
		}
	}

	if o := v.QueryParamsBlacklist; !o.IsUnknown() && !o.IsNull() {
		queryParamsBlacklistObjV, d := QueryParamsBlacklistType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		queryParamsBlacklist := queryParamsBlacklistObjV.(QueryParamsBlacklistValue)
		value := make([]string, 0, len(queryParamsBlacklist.Value.Elements()))
		if !queryParamsBlacklist.Value.IsNull() && !queryParamsBlacklist.Value.IsUnknown() {
			diags.Append(queryParamsBlacklist.Value.ElementsAs(ctx, &value, true)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.QueryParamsBlacklist = &resources.ResourceOptionsStringListOption{
			Enabled: queryParamsBlacklist.Enabled.ValueBool(),
			Value:   value,
		}
	}

	if o := v.QueryParamsWhitelist; !o.IsUnknown() && !o.IsNull() {
		queryParamsWhitelistObjV, d := QueryParamsWhitelistType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		queryParamsWhitelist := queryParamsWhitelistObjV.(QueryParamsWhitelistValue)
		value := make([]string, 0, len(queryParamsWhitelist.Value.Elements()))
		if !queryParamsWhitelist.Value.IsNull() && !queryParamsWhitelist.Value.IsUnknown() {
			diags.Append(queryParamsWhitelist.Value.ElementsAs(ctx, &value, true)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.QueryParamsWhitelist = &resources.ResourceOptionsStringListOption{
			Enabled: queryParamsWhitelist.Enabled.ValueBool(),
			Value:   value,
		}
	}

	if o := v.ReferrerAcl; !o.IsUnknown() && !o.IsNull() {
		referrerAclObjV, d := ReferrerAclType{}.ValueFromObject(ctx, v.ReferrerAcl)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		referrerAcl := referrerAclObjV.(ReferrerAclValue)
		exceptedValues := make([]string, 0, len(referrerAcl.ExceptedValues.Elements()))
		if !referrerAcl.ExceptedValues.IsNull() && !referrerAcl.ExceptedValues.IsUnknown() {
			diags.Append(referrerAcl.ExceptedValues.ElementsAs(ctx, &exceptedValues, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}
		result.ReferrerACL = &resources.ResourceOptionsACLOption{
			Enabled:        referrerAcl.Enabled.ValueBool(),
			ExceptedValues: exceptedValues,
			PolicyType:     resources.ResourceACLPolicyType(referrerAcl.PolicyType.ValueString()),
		}
	}

	if o := v.Slice; !o.IsUnknown() && !o.IsNull() {
		result.Slice = &resources.ResourceOptionsBoolOption{
			Enabled: v.Slice.ValueBool(),
			Value:   true,
		}
	}

	if o := v.Stale; !o.IsUnknown() && !o.IsNull() {
		staleObjV, d := StaleType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		stale := staleObjV.(StaleValue)
		value := make([]string, 0, len(stale.Value.Elements()))
		if !stale.Value.IsUnknown() && !stale.Value.IsNull() {
			diags.Append(stale.Value.ElementsAs(ctx, &value, true)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.Stale = &resources.ResourceOptionsStringListOption{
			Enabled: stale.Enabled.ValueBool(),
			Value:   value,
		}
	}

	if o := v.StaticHeaders; !o.IsUnknown() && !o.IsNull() {
		staticHeadersObjV, d := StaticHeadersType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		staticHeaders := staticHeadersObjV.(StaticHeadersValue)
		value := make(map[string]string, len(staticHeaders.Value.Elements()))
		if !staticHeaders.Value.IsUnknown() && !staticHeaders.Value.IsNull() {
			diags.Append(staticHeaders.Value.ElementsAs(ctx, &value, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.StaticHeaders = &resources.ResourceOptionsStringMapOption{
			Enabled: staticHeaders.Enabled.ValueBool(),
			Value:   value,
		}
	}

	if o := v.StaticRequestHeaders; !o.IsUnknown() && !o.IsNull() {
		staticRequestHeadersObjV, d := StaticRequestHeadersType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		staticRequestHeaders := staticRequestHeadersObjV.(StaticRequestHeadersValue)
		value := make(map[string]string, len(staticRequestHeaders.Value.Elements()))
		if !staticRequestHeaders.Value.IsNull() && !staticRequestHeaders.Value.IsUnknown() {
			diags.Append(staticRequestHeaders.Value.ElementsAs(ctx, &value, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		result.StaticRequestHeaders = &resources.ResourceOptionsStringMapOption{
			Enabled: staticRequestHeaders.Enabled.ValueBool(),
			Value:   value,
		}
	}

	if o := v.SecureKey; !o.IsUnknown() && !o.IsNull() {
		secureKeyObjV, d := SecureKeyType{}.ValueFromObject(ctx, o)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		secureKey := secureKeyObjV.(SecureKeyValue)
		result.SecureKey = &resources.ResourceOptionSecureKeyOption{
			Enabled: secureKey.Enabled.ValueBool(),
			Key:     secureKey.Key.ValueString(),
			Type:    secureKey.SecureKeyType.ValueInt64(),
		}
	}

	return result, diags
}

func (v OptionsValue) FromResourceOptions(ctx context.Context, opts *resources.ResourceOptions) (OptionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	if opts == nil {
		return NewOptionsValueNull(), nil
	}

	var allowedHttpMethods types.Object
	if o := opts.AllowedHttpMethods; o != nil {
		value, d := types.ListValueFrom(ctx, types.StringType, o.Value)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		allowedHttpMethods, d = AllowedHttpMethodsValue{
			Value:   value,
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		allowedHttpMethods, d = NewAllowedHttpMethodsValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var brotliCompression types.Object
	if o := opts.BrotliCompression; o != nil {
		value, d := types.SetValueFrom(ctx, types.StringType, o.Value)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		brotliCompression, d = BrotliCompressionValue{
			Value:   value,
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		brotliCompression, d = NewBrotliCompressionValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var browserCacheSettings types.Object
	if o := opts.BrowserCacheSettings; o != nil {
		browserCacheSettings, d = BrowserCacheSettingsValue{
			Value:   types.StringValue(o.Value),
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		browserCacheSettings, d = NewBrowserCacheSettingsValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var cors types.Object
	if o := opts.CORS; o != nil {
		value, d := types.ListValueFrom(ctx, types.StringType, o.Value)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		cors, d = CorsValue{
			Value:   value,
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		cors, d = NewCorsValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var countryAcl types.Object
	if o := opts.CountryACL; o != nil {
		exceptedValues, d := types.ListValueFrom(ctx, types.StringType, o.ExceptedValues)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		countryAcl, d = CountryAclValue{
			ExceptedValues: exceptedValues,
			PolicyType:     types.StringValue(string(o.PolicyType)),
			Enabled:        types.BoolValue(o.Enabled),
			state:          attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		countryAcl, d = NewCountryAclValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var edgeCacheSettings types.Object
	if o := opts.EdgeCacheSettings; o != nil {
		customValues, d := types.MapValueFrom(ctx, types.StringType, o.CustomValues)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		edgeCacheSettings, d = EdgeCacheSettingsValue{
			CustomValues: customValues,
			Default:      types.StringValue(o.Default),
			Enabled:      types.BoolValue(o.Enabled),
			Value:        types.StringValue(o.Value),
			state:        attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		edgeCacheSettings, d = NewEdgeCacheSettingsValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var forceReturn types.Object
	if o := opts.ForceReturn; o != nil {
		var d diag.Diagnostics
		forceReturn, d = ForceReturnValue{
			Body:    types.StringValue(o.Body),
			Code:    types.Int64Value(int64(o.Code)),
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		forceReturn, d = NewForceReturnValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var hostHeader types.Object
	if o := opts.HostHeader; o != nil {
		hostHeader, d = HostHeaderValue{
			Value:   types.StringValue(o.Value),
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		hostHeader, d = NewHostHeaderValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var ipAddressAcl types.Object
	if o := opts.IpAddressACL; o != nil {
		exceptedValues, d := types.ListValueFrom(ctx, types.StringType, o.ExceptedValues)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		ipAddressAcl, d = IpAddressAclValue{
			ExceptedValues: exceptedValues,
			PolicyType:     types.StringValue(string(o.PolicyType)),
			Enabled:        types.BoolValue(o.Enabled),
			state:          attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		ipAddressAcl, d = NewIpAddressAclValueNull().ToObjectValue(ctx)
		if diags.HasError() {
			return NewOptionsValueUnknown(), d
		}
	}

	var queryParamsBlacklist types.Object
	if o := opts.QueryParamsBlacklist; o != nil {
		value, d := types.ListValueFrom(ctx, types.StringType, o.Value)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		queryParamsBlacklist, d = QueryParamsBlacklistValue{
			Value:   value,
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		queryParamsBlacklist, d = NewQueryParamsBlacklistValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var queryParamsWhitelist types.Object
	if o := opts.QueryParamsWhitelist; o != nil {
		value, d := types.ListValueFrom(ctx, types.StringType, o.Value)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		queryParamsWhitelist, d = QueryParamsWhitelistValue{
			Value:   value,
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		queryParamsWhitelist, d = NewQueryParamsWhitelistValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var referrerAcl types.Object
	if o := opts.ReferrerACL; o != nil {
		exceptedValues, d := types.ListValueFrom(ctx, types.StringType, o.ExceptedValues)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		referrerAcl, d = ReferrerAclValue{
			ExceptedValues: exceptedValues,
			PolicyType:     types.StringValue(string(o.PolicyType)),
			Enabled:        types.BoolValue(o.Enabled),
			state:          attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		referrerAcl, d = NewReferrerAclValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var stale types.Object
	if o := opts.Stale; o != nil {
		value, d := types.ListValueFrom(ctx, types.StringType, o.Value)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		stale, d = StaleValue{
			Value:   value,
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		stale, d = NewStaleValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var staticHeaders types.Object
	if o := opts.StaticHeaders; o != nil {
		value, d := types.MapValueFrom(ctx, types.StringType, o.Value)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		staticHeaders, d = StaticHeadersValue{
			Value:   value,
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		staticHeaders, d = NewStaticHeadersValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var staticRequestHeaders types.Object
	if o := opts.StaticRequestHeaders; o != nil {
		value, d := types.MapValueFrom(ctx, types.StringType, o.Value)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
		staticRequestHeaders, d = StaticRequestHeadersValue{
			Value:   value,
			Enabled: types.BoolValue(o.Enabled),
			state:   attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		staticRequestHeaders, d = NewStaticRequestHeadersValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	var secureKey types.Object
	if o := opts.SecureKey; o != nil {
		secureKey, d = SecureKeyValue{
			Enabled:       types.BoolValue(o.Enabled),
			Key:           types.StringValue(o.Key),
			SecureKeyType: types.Int64Value(o.Type),
			state:         attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	} else {
		secureKey, d = NewSecureKeyValueNull().ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return NewOptionsValueUnknown(), diags
		}
	}

	return OptionsValue{
		AllowedHttpMethods:   allowedHttpMethods,
		BrotliCompression:    brotliCompression,
		BrowserCacheSettings: browserCacheSettings,
		Cors:                 cors,
		CountryAcl:           countryAcl,
		EdgeCacheSettings:    edgeCacheSettings,
		FetchCompressed:      flattenBoolOption(opts.FetchCompressed),
		ForceReturn:          forceReturn,
		ForwardHostHeader:    flattenBoolOption(opts.ForwardHostHeader),
		GzipOn:               flattenBoolOption(opts.GzipOn),
		HostHeader:           hostHeader,
		IgnoreCookie:         flattenBoolOption(opts.IgnoreCookie),
		IgnoreQueryString:    flattenBoolOption(opts.IgnoreQueryString),
		IpAddressAcl:         ipAddressAcl,
		QueryParamsBlacklist: queryParamsBlacklist,
		QueryParamsWhitelist: queryParamsWhitelist,
		ReferrerAcl:          referrerAcl,
		Slice:                flattenBoolOption(opts.Slice),
		Stale:                stale,
		StaticHeaders:        staticHeaders,
		StaticRequestHeaders: staticRequestHeaders,
		SecureKey:            secureKey,
		state:                attr.ValueStateKnown,
	}, diags
}

func (v ShieldingValue) FromResourceShielding(resourceShielding *resources.ResourceShielding) ShieldingValue {
	if resourceShielding == nil {
		return NewShieldingValueNull()
	}

	var enabled bool
	var popId int64

	if resourceShielding.ShieldingPop != nil {
		enabled = true
		popId = int64(*resourceShielding.ShieldingPop)
	}

	return ShieldingValue{
		Enabled: types.BoolValue(enabled),
		PopId:   types.Int64Value(popId),
		state:   attr.ValueStateKnown,
	}
}

func (v ShieldingValue) ToUpdateShieldingOpts() *resources.UpdateShieldingOpts {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	if !v.Enabled.ValueBool() {
		return &resources.UpdateShieldingOpts{}
	}

	return &resources.UpdateShieldingOpts{
		ShieldingPop: util.PointerOf(int(v.PopId.ValueInt64())),
	}
}

type SslOpts struct {
	LeEnabled bool
	Enabled   bool
	Data      int
}

func (v SslCertificateValue) ToSslOpts() *SslOpts {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	switch v.SslCertificateType.ValueString() {
	case string(SslCertificateProviderTypeNotUsed):
		return &SslOpts{}
	case string(SslCertificateProviderTypeLetsEncrypt):
		return &SslOpts{}
	case string(SslCertificateProviderTypeOwn):
		return &SslOpts{
			Enabled: true,
			Data:    int(v.Id.ValueInt64()),
		}
	}

	return nil
}

func (v SslCertificateValue) FromSslOpts(ctx context.Context, opts *SslOpts) SslCertificateValue {
	if opts == nil {
		return NewSslCertificateValueNull()
	}

	var sslType string
	var status string
	switch {
	case opts.Enabled && opts.Data != 0:
		status = string(SslCertificateStatusReady)
		sslType = string(SslCertificateProviderTypeOwn)
	case opts.LeEnabled:
		status = string(SslCertificateStatusReady)
		sslType = string(SslCertificateProviderTypeLetsEncrypt)
	default:
		sslType = string(SslCertificateProviderTypeNotUsed)
	}

	return SslCertificateValue{
		Id:                 types.Int64Value(int64(opts.Data)),
		Status:             types.StringValue(status),
		SslCertificateType: types.StringValue(sslType),
		state:              attr.ValueStateKnown,
	}
}

func expandBoolOption(opt types.Bool) *resources.ResourceOptionsBoolOption {
	return &resources.ResourceOptionsBoolOption{
		Enabled: opt.ValueBool(),
		Value:   true,
	}
}

func flattenBoolOption(opt *resources.ResourceOptionsBoolOption) types.Bool {
	if opt == nil {
		return types.BoolNull()
	}
	return types.BoolValue(opt.Value && opt.Enabled)
}
