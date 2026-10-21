package kms_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// newTestState fills unspecified attributes with null; callers supply planned defaults explicitly.
func newTestState(t *testing.T, state tfsdk.State, input map[string]interface{}) tfsdk.State {
	t.Helper()
	object := state.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	values := map[string]tftypes.Value{}
	for name, typ := range object.AttributeTypes {
		value := input[name]
		if list, ok := typ.(tftypes.List); ok && value != nil {
			elements := []tftypes.Value{}
			for _, v := range value.([]interface{}) {
				elements = append(elements, tftypes.NewValue(list.ElementType, v))
			}
			value = elements
		}
		values[name] = tftypes.NewValue(typ, value)
	}
	state.Raw = tftypes.NewValue(object, values)
	return state
}
