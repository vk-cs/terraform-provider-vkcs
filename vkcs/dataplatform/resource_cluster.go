package dataplatform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/dataplatform/resource_cluster"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

var (
	_ resource.Resource = (*clusterResource)(nil)
)

func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

type clusterResource struct {
	config clients.Config
}

func (r *clusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataplatform_cluster"
}

func (r *clusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_cluster.ClusterResourceSchema(ctx)
}

func (r *clusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}
