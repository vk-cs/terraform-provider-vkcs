package frameworkwrapper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/providerwrapper/framework/resource/customschema"
)

var (
	_ resource.Resource = &ResourceWrapper{}
)

func NewResourceWrapper(resource resource.Resource, resourceJSON jsonschema.ResourceJSON) *ResourceWrapper {
	return &ResourceWrapper{
		resource:     resource,
		resourceJSON: resourceJSON,
	}
}

type ResourceWrapper struct {
	resource     resource.Resource
	resourceJSON jsonschema.ResourceJSON
}

func (rw *ResourceWrapper) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	rw.resource.Metadata(ctx, req, resp)
}

func (rw *ResourceWrapper) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	rw.resource.Schema(ctx, req, resp)
	resp.Schema = customschema.CustomizeSchema(rw.resourceJSON, resp.Schema)
}

func (rw *ResourceWrapper) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	rw.resource.Create(ctx, req, resp)
}

func (rw *ResourceWrapper) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	rw.resource.Read(ctx, req, resp)
}

func (rw *ResourceWrapper) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	rw.resource.Update(ctx, req, resp)
}

func (rw *ResourceWrapper) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	rw.resource.Delete(ctx, req, resp)
}

func (rw *ResourceWrapper) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if rs, ok := rw.resource.(resource.ResourceWithConfigure); ok {
		rs.Configure(ctx, req, resp)
	}
}

func (rw *ResourceWrapper) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if rs, ok := rw.resource.(resource.ResourceWithConfigValidators); ok {
		return rs.ConfigValidators(ctx)
	}
	return nil
}

func (rw *ResourceWrapper) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if rs, ok := rw.resource.(resource.ResourceWithImportState); ok {
		rs.ImportState(ctx, req, resp)
	} else {
		resp.Diagnostics.AddError(
			"Resource Import Not Implemented",
			"This resource does not support import.",
		)
	}
}

func (rw *ResourceWrapper) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if rs, ok := rw.resource.(resource.ResourceWithModifyPlan); ok {
		rs.ModifyPlan(ctx, req, resp)
	}
}

func (rw *ResourceWrapper) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	if rs, ok := rw.resource.(resource.ResourceWithUpgradeState); ok {
		return rs.UpgradeState(ctx)
	}
	return make(map[int64]resource.StateUpgrader, 0)
}
