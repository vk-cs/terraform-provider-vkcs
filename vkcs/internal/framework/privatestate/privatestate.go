package privatestate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type PrivateState interface {
	SetKey(ctx context.Context, key string, data []byte) diag.Diagnostics
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

func ReadInto(ctx context.Context, privateState PrivateState, key string, target interface{}) diag.Diagnostics {
	data, diags := privateState.GetKey(ctx, key)
	if diags.HasError() {
		return diags
	}

	if data == nil {
		return nil
	}

	err := json.Unmarshal(data, target)
	if err != nil {
		diags.AddError(fmt.Sprintf("error reading private state data for %q", key), err.Error())
	}

	return diags
}

func WriteFrom(ctx context.Context, privateState PrivateState, key string, from interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	data, err := json.Marshal(from)
	if err != nil {
		diags.AddError(fmt.Sprintf("error writing private state data for %q", key), err.Error())
		return diags
	}

	return privateState.SetKey(ctx, key, data)
}
