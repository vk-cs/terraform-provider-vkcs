package resource_cluster

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// configsUsersAttribute walks the schema to extract configs.users as schema.ListNestedAttribute.
func configsUsersAttribute(t *testing.T, ctx context.Context) schema.ListNestedAttribute {
	t.Helper()
	s := ClusterResourceSchema(ctx)
	configsAttr, ok := s.Attributes["configs"].(schema.SingleNestedAttribute)
	require.True(t, ok, "configs must be SingleNestedAttribute")
	usersAttr, ok := configsAttr.Attributes["users"].(schema.ListNestedAttribute)
	require.True(t, ok, "configs.users must be ListNestedAttribute")
	return usersAttr
}

func configsWarehousesAttribute(t *testing.T, ctx context.Context) schema.ListNestedAttribute {
	t.Helper()
	s := ClusterResourceSchema(ctx)
	configsAttr, ok := s.Attributes["configs"].(schema.SingleNestedAttribute)
	require.True(t, ok)
	whAttr, ok := configsAttr.Attributes["warehouses"].(schema.ListNestedAttribute)
	require.True(t, ok, "configs.warehouses must be ListNestedAttribute")
	return whAttr
}

func TestSchema_Users_IsRequired(t *testing.T) {
	ctx := context.Background()
	usersAttr := configsUsersAttribute(t, ctx)

	assert.True(t, usersAttr.IsRequired(), "configs.users must be Required")
	assert.False(t, usersAttr.IsComputed(), "configs.users must not be Computed (incompatible with Required)")
	require.NotEmpty(t, usersAttr.Validators, "configs.users must have at least one validator")
}

func TestSchema_Warehouses_IsRequired(t *testing.T) {
	ctx := context.Background()
	whAttr := configsWarehousesAttribute(t, ctx)

	assert.True(t, whAttr.IsRequired(), "configs.warehouses must be Required")
	assert.False(t, whAttr.IsComputed())
	require.NotEmpty(t, whAttr.Validators, "configs.warehouses must have at least one validator")
}

func runListValidators(ctx context.Context, validators []validator.List, value types.List) bool {
	req := validator.ListRequest{
		Path:           path.Root("test"),
		PathExpression: path.MatchRoot("test"),
		ConfigValue:    value,
	}
	resp := validator.ListResponse{}
	for _, v := range validators {
		v.ValidateList(ctx, req, &resp)
	}
	return resp.Diagnostics.HasError()
}

func TestSchema_Users_EmptyListFailsValidation(t *testing.T) {
	ctx := context.Background()
	usersAttr := configsUsersAttribute(t, ctx)
	elemType := usersAttr.NestedObject.Type()

	t.Run("empty list errors at plan time", func(t *testing.T) {
		emptyList := types.ListValueMust(elemType, []attr.Value{})
		assert.True(t, runListValidators(ctx, usersAttr.Validators, emptyList),
			"SizeAtLeast(1) must return a diagnostic for an empty list")
	})

	t.Run("single user passes validation", func(t *testing.T) {
		one := newUser("vkdata", "Pass!", "dbOwner")
		lv, diags := types.ListValueFrom(ctx, elemType, []ConfigsUsersValue{one})
		require.False(t, diags.HasError())
		assert.False(t, runListValidators(ctx, usersAttr.Validators, lv),
			"a non-empty users list must pass validation")
	})
}

func runStringValidators(ctx context.Context, validators []validator.String, value types.String) bool {
	req := validator.StringRequest{
		Path:           path.Root("test"),
		PathExpression: path.MatchRoot("test"),
		ConfigValue:    value,
	}
	resp := validator.StringResponse{}
	for _, v := range validators {
		v.ValidateString(ctx, req, &resp)
	}
	return resp.Diagnostics.HasError()
}

func TestSchema_FloatingIPPool_OnlyAutoOrOmitted(t *testing.T) {
	ctx := context.Background()
	s := ClusterResourceSchema(ctx)
	fipAttr, ok := s.Attributes["floating_ip_pool"].(schema.StringAttribute)
	require.True(t, ok, "floating_ip_pool must be StringAttribute")
	require.NotEmpty(t, fipAttr.Validators, "floating_ip_pool must have a validator")

	cases := []struct {
		name      string
		value     types.String
		expectErr bool
	}{
		{"auto - ok", types.StringValue("auto"), false},
		{"null (attribute omitted in HCL) - ok", types.StringNull(), false},
		{"unknown - ok", types.StringUnknown(), false},
		{"empty string in HCL - error (only auto or omission allowed)", types.StringValue(""), true},
		{"arbitrary uuid - error", types.StringValue("11111111-2222-3333-4444-555555555555"), true},
		{"arbitrary string - error", types.StringValue("default"), true},
		{"uppercase AUTO - error (case-sensitive)", types.StringValue("AUTO"), true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotErr := runStringValidators(ctx, fipAttr.Validators, c.value)
			assert.Equal(t, c.expectErr, gotErr,
				"value=%q expectErr=%v gotErr=%v", c.value.ValueString(), c.expectErr, gotErr)
		})
	}
}

func TestSchema_Warehouses_SizeBetween1_1(t *testing.T) {
	ctx := context.Background()
	whAttr := configsWarehousesAttribute(t, ctx)
	elemType := whAttr.NestedObject.Type()

	t.Run("empty list errors", func(t *testing.T) {
		emptyList := types.ListValueMust(elemType, []attr.Value{})
		assert.True(t, runListValidators(ctx, whAttr.Validators, emptyList),
			"SizeBetween(1,1) must reject an empty list")
	})

	t.Run("exactly one warehouse passes validation", func(t *testing.T) {
		one := newWarehouse(ctx, "db_customer")
		lv, diags := types.ListValueFrom(ctx, elemType, []ConfigsWarehousesValue{one})
		require.False(t, diags.HasError())
		assert.False(t, runListValidators(ctx, whAttr.Validators, lv),
			"a single warehouse must pass validation")
	})

	t.Run("two warehouses error", func(t *testing.T) {
		two := []ConfigsWarehousesValue{
			newWarehouse(ctx, "wh1"),
			newWarehouse(ctx, "wh2"),
		}
		lv, diags := types.ListValueFrom(ctx, elemType, two)
		require.False(t, diags.HasError())
		assert.True(t, runListValidators(ctx, whAttr.Validators, lv),
			"SizeBetween(1,1) must reject two warehouses - current implementation supports only one")
	})
}
