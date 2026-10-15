package kms_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

// The HTTP behavior tests share this harness, which builds real Framework
// requests and state without depending on SDK ResourceData.
type frameworkTestEndpoint struct {
	resource   resource.Resource
	datasource datasource.DataSource
	state      tfsdk.State
	defaults   map[string]any
}
type frameworkTestData struct {
	t     *testing.T
	state tfsdk.State
}

func newFrameworkTestEndpoint(t *testing.T, implementation any) *frameworkTestEndpoint {
	t.Helper()
	ctx := context.Background()
	e := &frameworkTestEndpoint{defaults: map[string]any{}}
	switch impl := implementation.(type) {
	case resource.Resource:
		e.resource = impl
		var resp resource.SchemaResponse
		impl.Schema(ctx, resource.SchemaRequest{}, &resp)
		require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
		e.state.Schema = resp.Schema
		for name, a := range resp.Schema.Attributes {
			switch a := a.(type) {
			case rschema.StringAttribute:
				if a.Default != nil {
					var d defaults.StringResponse
					a.Default.DefaultString(ctx, defaults.StringRequest{}, &d)
					e.defaults[name] = d.PlanValue.ValueString()
				}
			case rschema.BoolAttribute:
				if a.Default != nil {
					var d defaults.BoolResponse
					a.Default.DefaultBool(ctx, defaults.BoolRequest{}, &d)
					e.defaults[name] = d.PlanValue.ValueBool()
				}
			}
		}
	case datasource.DataSource:
		e.datasource = impl
		var resp datasource.SchemaResponse
		impl.Schema(ctx, datasource.SchemaRequest{}, &resp)
		require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
		e.state.Schema = resp.Schema
	default:
		t.Fatalf("unsupported implementation %T", implementation)
	}
	return e
}

func newFrameworkTestData(t *testing.T, e *frameworkTestEndpoint, input map[string]interface{}) *frameworkTestData {
	t.Helper()
	ctx := context.Background()
	object := e.state.Schema.Type().TerraformType(ctx).(tftypes.Object)
	values := map[string]tftypes.Value{}
	for name, typ := range object.AttributeTypes {
		value, ok := input[name]
		if !ok {
			value = e.defaults[name]
		}
		if list, ok := typ.(tftypes.List); ok && value != nil {
			var elements []tftypes.Value
			for _, v := range value.([]interface{}) {
				elements = append(elements, tftypes.NewValue(list.ElementType, v))
			}
			value = elements
		}
		values[name] = tftypes.NewValue(typ, value)
	}
	state := e.state
	state.Raw = tftypes.NewValue(object, values)
	return &frameworkTestData{t: t, state: state}
}
func (d *frameworkTestData) SetId(id string) {
	d.t.Helper()
	require.False(d.t, d.state.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), id).HasError())
}
func (d *frameworkTestData) Id() string {
	v := d.Get(kms.EntityFieldID)
	if v == nil {
		return ""
	}
	return v.(string)
}
func (d *frameworkTestData) Get(name string) any {
	d.t.Helper()
	ctx := context.Background()
	var v attr.Value
	diags := d.state.GetAttribute(ctx, path.Root(name), &v)
	require.False(d.t, diags.HasError(), "%v", diags)
	if v.IsNull() {
		switch v.(type) {
		case types.String:
			return ""
		case types.Int64:
			return 0
		case types.Bool:
			return false
		case types.List:
			return []interface{}{}
		}
	}
	switch v := v.(type) {
	case types.String:
		return v.ValueString()
	case types.Int64:
		return int(v.ValueInt64())
	case types.Bool:
		return v.ValueBool()
	case types.List:
		result := []interface{}{}
		for _, x := range v.Elements() {
			result = append(result, x.(types.String).ValueString())
		}
		return result
	}
	return nil
}
func (e *frameworkTestEndpoint) configure(ctx context.Context, config any) diag.Diagnostics {
	if e.resource != nil {
		var resp resource.ConfigureResponse
		e.resource.(resource.ResourceWithConfigure).Configure(ctx, resource.ConfigureRequest{ProviderData: config}, &resp)
		return resp.Diagnostics
	}
	var resp datasource.ConfigureResponse
	e.datasource.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: config}, &resp)
	return resp.Diagnostics
}
func (e *frameworkTestEndpoint) ReadContext(ctx context.Context, d *frameworkTestData, config any) diag.Diagnostics {
	if diags := e.configure(ctx, config); diags.HasError() {
		return diags
	}
	if e.resource != nil {
		resp := resource.ReadResponse{State: d.state}
		e.resource.Read(ctx, resource.ReadRequest{State: d.state}, &resp)
		d.state = resp.State
		return resp.Diagnostics
	}
	resp := datasource.ReadResponse{State: d.state}
	e.datasource.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: d.state.Schema, Raw: d.state.Raw}}, &resp)
	d.state = resp.State
	return resp.Diagnostics
}
func (e *frameworkTestEndpoint) CreateContext(ctx context.Context, d *frameworkTestData, config any) diag.Diagnostics {
	if diags := e.configure(ctx, config); diags.HasError() {
		return diags
	}
	resp := resource.CreateResponse{State: d.state}
	e.resource.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: d.state.Schema, Raw: d.state.Raw}}, &resp)
	d.state = resp.State
	return resp.Diagnostics
}
func (e *frameworkTestEndpoint) UpdateContext(ctx context.Context, d *frameworkTestData, config any) diag.Diagnostics {
	if diags := e.configure(ctx, config); diags.HasError() {
		return diags
	}
	resp := resource.UpdateResponse{State: d.state}
	e.resource.Update(ctx, resource.UpdateRequest{State: d.state, Plan: tfsdk.Plan{Schema: d.state.Schema, Raw: d.state.Raw}}, &resp)
	d.state = resp.State
	return resp.Diagnostics
}
func (e *frameworkTestEndpoint) DeleteContext(ctx context.Context, d *frameworkTestData, config any) diag.Diagnostics {
	if diags := e.configure(ctx, config); diags.HasError() {
		return diags
	}
	resp := resource.DeleteResponse{State: d.state}
	e.resource.Delete(ctx, resource.DeleteRequest{State: d.state}, &resp)
	return resp.Diagnostics
}
