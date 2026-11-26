package lb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/netip"
	"time"

	"github.com/gophercloud/gophercloud"
	octaviapools "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	ipools "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/lb/v2/pools"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	membersDefaultTimeout      = 10 * time.Minute
	privateExplicitEmptySubnet = "explicit_empty_subnet"
	privateRegion              = "region"
	putMembersRetryMinDelay    = 1 * time.Second
	putMembersRetryMaxDelay    = 10 * time.Second
)

var putMembersRetryableCodes = []int{409, 500, 502, 503, 504}

type privateGetter interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

type privateSetter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

type explicitEmptySubnets map[memberKey]bool

type explicitEmptySubnetKey struct {
	Address      string `json:"address"`
	ProtocolPort int64  `json:"protocol_port"`
}

func explicitEmptySubnetsFromConfig(configMembers []memberModel) explicitEmptySubnets {
	res := explicitEmptySubnets{}
	for _, m := range configMembers {
		k, ok := m.key()
		if ok && !m.SubnetID.IsNull() && !m.SubnetID.IsUnknown() && m.SubnetID.ValueString() == "" {
			res[k] = true
		}
	}
	return res
}

func loadExplicitEmptySubnets(ctx context.Context, private privateGetter) (explicitEmptySubnets, diag.Diagnostics) {
	res := explicitEmptySubnets{}
	if private == nil {
		return res, nil
	}
	raw, diags := private.GetKey(ctx, privateExplicitEmptySubnet)
	if diags.HasError() || len(raw) == 0 {
		return res, diags
	}
	var keys []explicitEmptySubnetKey
	if err := json.Unmarshal(raw, &keys); err != nil {
		diags.AddError("Error decoding vkcs_lb_members private state", err.Error())
		return res, diags
	}
	for _, k := range keys {
		res[memberKey{address: normalizeAddress(k.Address), protocolPort: k.ProtocolPort}] = true
	}
	return res, diags
}

func saveExplicitEmptySubnets(ctx context.Context, private privateSetter, subnets explicitEmptySubnets) diag.Diagnostics {
	keys := make([]explicitEmptySubnetKey, 0, len(subnets))
	for k := range subnets {
		keys = append(keys, explicitEmptySubnetKey{Address: k.address, ProtocolPort: k.protocolPort})
	}
	raw, err := json.Marshal(keys)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Error encoding vkcs_lb_members private state", err.Error())
		return diags
	}
	return private.SetKey(ctx, privateExplicitEmptySubnet, raw)
}

func configMembersFrom(ctx context.Context, config interface {
	GetAttribute(context.Context, path.Path, interface{}) diag.Diagnostics
}) ([]memberModel, diag.Diagnostics) {
	var configSet types.Set
	diags := config.GetAttribute(ctx, path.Root("member"), &configSet)
	if diags.HasError() {
		return nil, diags
	}
	members, d := membersFromSet(ctx, configSet)
	diags.Append(d...)
	return members, diags
}

var (
	_ resource.Resource                = (*membersResource)(nil)
	_ resource.ResourceWithConfigure   = (*membersResource)(nil)
	_ resource.ResourceWithImportState = (*membersResource)(nil)
	_ resource.ResourceWithModifyPlan  = (*membersResource)(nil)
)

func NewMembersResource() resource.Resource {
	return &membersResource{}
}

type membersResource struct {
	config clients.Config
}

type membersResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Region   types.String   `tfsdk:"region"`
	PoolID   types.String   `tfsdk:"pool_id"`
	Member   types.Set      `tfsdk:"member"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type memberModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Address      types.String `tfsdk:"address"`
	ProtocolPort types.Int64  `tfsdk:"protocol_port"`
	Weight       types.Int64  `tfsdk:"weight"`
	SubnetID     types.String `tfsdk:"subnet_id"`
	Backup       types.Bool   `tfsdk:"backup"`
	AdminStateUp types.Bool   `tfsdk:"admin_state_up"`
}

type memberKey struct {
	address      string
	protocolPort int64
}

func (m memberModel) key() (memberKey, bool) {
	if m.Address.IsUnknown() || m.Address.IsNull() || m.ProtocolPort.IsUnknown() || m.ProtocolPort.IsNull() {
		return memberKey{}, false
	}
	return memberKey{address: normalizeAddress(m.Address.ValueString()), protocolPort: m.ProtocolPort.ValueInt64()}, true
}

func normalizeAddress(address string) string {
	if ip, err := netip.ParseAddr(address); err == nil {
		return ip.String()
	}
	return address
}

var memberAttrTypes = map[string]attr.Type{
	"id":             types.StringType,
	"name":           types.StringType,
	"address":        types.StringType,
	"protocol_port":  types.Int64Type,
	"weight":         types.Int64Type,
	"subnet_id":      types.StringType,
	"backup":         types.BoolType,
	"admin_state_up": types.BoolType,
}

var memberObjectType = types.ObjectType{AttrTypes: memberAttrTypes}

func (r *membersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lb_members"
}

func (r *membersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					regionPlanModifier{r: r},
				},
				Description: "The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new members resource.",
			},

			"pool_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The id of the pool that members will be assigned to. Changing this creates a new members resource.",
			},
		},
		Blocks: map[string]schema.Block{
			"member": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The unique ID for the member.",
						},

						"name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Human-readable name for the member.",
						},

						"address": schema.StringAttribute{
							Required:    true,
							Description: "The IP address of the members to receive traffic from the load balancer.",
						},

						"protocol_port": schema.Int64Attribute{
							Required: true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
							Description: "The port on which to listen for client traffic.",
						},

						"weight": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							Validators: []validator.Int64{
								int64validator.Between(0, 256),
							},
							Description: "A positive integer value that indicates the relative portion of traffic that this members should receive from the pool. For example, a member with a weight of 10 receives five times as much traffic as a member with a weight of 2. Defaults to 1.",
						},

						"subnet_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The subnet in which to access the member.",
						},

						"backup": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "A bool that indicates whether the member is backup.",
						},

						"admin_state_up": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The administrative state of the member. A valid value is true (UP) or false (DOWN). Defaults to true.",
						},
					},
				},
				Description: "A set of dictionaries containing member parameters. The structure is described below.",
			},
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func (r *membersResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *membersResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan membersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Member.IsNull() || plan.Member.IsUnknown() {
		return
	}

	stateByKey := map[memberKey]memberModel{}
	if !req.State.Raw.IsNull() {
		var state membersResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
		stateMembers, diags := membersFromSet(ctx, state.Member)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, m := range stateMembers {
			if k, ok := m.key(); ok {
				stateByKey[k] = m
			}
		}
	}

	planMembers, diags := membersFromSet(ctx, plan.Member)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	configMembers, diags := configMembersFrom(ctx, req.Config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	explicitEmpty, diags := loadExplicitEmptySubnets(ctx, req.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i := range planMembers {
		p := &planMembers[i]

		var candidates []memberModel
		for _, c := range configMembers {
			if c.Address.Equal(p.Address) && c.ProtocolPort.Equal(p.ProtocolPort) {
				candidates = append(candidates, c)
			}
		}

		k, keyKnown := p.key()
		prior, hasPrior := stateByKey[k]
		hasPrior = keyKnown && hasPrior

		p.Name = planOrDefault(p.Name, candidates, func(m memberModel) types.String { return m.Name }, types.StringValue(""))
		if !isConfigured(p.SubnetID, candidates, func(m memberModel) types.String { return m.SubnetID }) {
			p.SubnetID = types.StringUnknown()
			if hasPrior && !prior.SubnetID.IsNull() && !prior.SubnetID.IsUnknown() &&
				(prior.SubnetID.ValueString() != "" || !explicitEmpty[k]) {
				p.SubnetID = prior.SubnetID
			}
		}
		p.Backup = planOrDefault(p.Backup, candidates, func(m memberModel) types.Bool { return m.Backup }, types.BoolValue(false))
		p.Weight = planOrDefault(p.Weight, candidates, func(m memberModel) types.Int64 { return m.Weight }, types.Int64Value(1))
		p.AdminStateUp = planOrDefault(p.AdminStateUp, candidates, func(m memberModel) types.Bool { return m.AdminStateUp }, types.BoolValue(true))

		p.ID = types.StringUnknown()
		if hasPrior && !prior.ID.IsNull() && !prior.ID.IsUnknown() {
			p.ID = prior.ID
		}
	}

	members, diags := types.SetValueFrom(ctx, memberObjectType, planMembers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("member"), members)...)
}

func isConfigured[T attr.Value](planned T, candidates []memberModel, get func(memberModel) T) bool {
	for _, c := range candidates {
		if v := get(c); !v.IsNull() && v.Equal(planned) {
			return true
		}
	}
	return len(candidates) == 0 && !planned.IsNull()
}

func planOrDefault[T attr.Value](planned T, candidates []memberModel, get func(memberModel) T, def T) T {
	if isConfigured(planned, candidates, get) {
		return planned
	}
	return def
}

func (r *membersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan membersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, membersDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	region, diags := r.resolveRegion(ctx, plan.Region, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	lbClient, err := r.lbClient(ctx, region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS loadbalancer client", err.Error())
		return
	}

	planMembers, diags := membersFromSet(ctx, plan.Member)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := plan.PoolID.ValueString()
	plan.ID = types.StringValue(poolID)
	plan.Region = regionValue(plan.Region, region)

	parentPool, err := putMembers(ctx, lbClient, poolID, expandMembers(planMembers), timeout)
	if err != nil {
		resp.Diagnostics.AddError("Error creating members", err.Error())
		return
	}
	resp.Diagnostics.Append(saveRegion(ctx, resp.Private, region)...)

	if err := waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout); err != nil {
		resp.Diagnostics.AddError("Error creating members", err.Error())
		plan.Member = types.SetNull(memberObjectType)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	configMembers, diags := configMembersFrom(ctx, req.Config)
	resp.Diagnostics.Append(diags...)
	explicitEmpty := explicitEmptySubnetsFromConfig(configMembers)
	resp.Diagnostics.Append(saveExplicitEmptySubnets(ctx, resp.Private, explicitEmpty)...)

	members, diags, err := readMembers(ctx, lbClient, poolID, planMembers, explicitEmpty)
	if err != nil {
		resp.Diagnostics.AddError("Error getting vkcs_lb_members", err.Error())
		members = types.SetNull(memberObjectType)
	}
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		members = types.SetNull(memberObjectType)
	}
	plan.Member = members

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *membersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state membersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region, diags := r.resolveRegion(ctx, state.Region, req.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	lbClient, err := r.lbClient(ctx, region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS loadbalancer client", err.Error())
		return
	}

	priorMembers, diags := membersFromSet(ctx, state.Member)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := state.ID.ValueString()
	explicitEmpty, diags := loadExplicitEmptySubnets(ctx, req.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	members, diags, err := readMembers(ctx, lbClient, poolID, priorMembers, explicitEmpty)
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error getting vkcs_lb_members", err.Error())
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.PoolID = types.StringValue(poolID)
	state.Region = regionValue(state.Region, region)
	state.Member = members

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *membersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state membersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, membersDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	configMembers, diags := configMembersFrom(ctx, req.Config)
	resp.Diagnostics.Append(diags...)
	explicitEmpty := explicitEmptySubnetsFromConfig(configMembers)

	region, diags := r.resolveRegion(ctx, plan.Region, req.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.Region = regionValue(plan.Region, region)

	if plan.Member.Equal(state.Member) {
		resp.Diagnostics.Append(saveExplicitEmptySubnets(ctx, resp.Private, explicitEmpty)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	lbClient, err := r.lbClient(ctx, region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS loadbalancer client", err.Error())
		return
	}

	planMembers, diags := membersFromSet(ctx, plan.Member)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := state.ID.ValueString()
	parentPool, err := putMembers(ctx, lbClient, poolID, expandMembers(planMembers), timeout)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to update member %s", poolID), err.Error())
		return
	}
	if err := waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to update member %s", poolID), err.Error())
		return
	}

	members, diags, err := readMembers(ctx, lbClient, poolID, planMembers, explicitEmpty)
	if err != nil {
		resp.Diagnostics.AddError("Error getting vkcs_lb_members", err.Error())
		return
	}
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	plan.Member = members

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	resp.Diagnostics.Append(saveExplicitEmptySubnets(ctx, resp.Private, explicitEmpty)...)
}

func (r *membersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state membersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, membersDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	region, diags := r.resolveRegion(ctx, state.Region, req.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	lbClient, err := r.lbClient(ctx, region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS loadbalancer client", err.Error())
		return
	}

	parentPool, err := putMembers(ctx, lbClient, state.ID.ValueString(), []octaviapools.BatchUpdateMemberOpts{}, timeout)
	if err == nil {
		err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	}
	if errutil.IsNotFound(err) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error deleting members", err.Error())
	}
}

func (r *membersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *membersResource) lbClient(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	client, err := r.config.LoadBalancerV2Client(region)
	if err != nil {
		return nil, err
	}
	return withRequestContext(ctx, client), nil
}

func (r *membersResource) resolveRegion(ctx context.Context, v types.String, private privateGetter) (string, diag.Diagnostics) {
	if region := v.ValueString(); region != "" {
		return region, nil
	}
	var diags diag.Diagnostics
	if private != nil {
		raw, d := private.GetKey(ctx, privateRegion)
		diags.Append(d...)
		if diags.HasError() {
			return "", diags
		}
		var region string
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &region); err != nil {
				diags.AddError("Error decoding vkcs_lb_members private state", err.Error())
				return "", diags
			}
		}
		if region != "" {
			return region, diags
		}
	}
	return r.config.GetRegion(), diags
}

func saveRegion(ctx context.Context, private privateSetter, region string) diag.Diagnostics {
	raw, err := json.Marshal(region)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Error encoding vkcs_lb_members private state", err.Error())
		return diags
	}
	return private.SetKey(ctx, privateRegion, raw)
}

type regionPlanModifier struct {
	r *membersResource
}

func (m regionPlanModifier) Description(context.Context) string {
	return "keeps the region the members were created in; an empty value means the provider region at creation time"
}

func (m regionPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m regionPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	if req.ConfigValue.IsUnknown() {
		resp.RequiresReplace = true
		return
	}
	if req.ConfigValue.IsNull() || req.ConfigValue.ValueString() == "" {
		resp.PlanValue = req.StateValue
		return
	}
	actual, diags := m.r.resolveRegion(ctx, req.StateValue, req.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if req.ConfigValue.ValueString() != actual {
		resp.RequiresReplace = true
	}
}

func regionValue(v types.String, region string) types.String {
	if !v.IsNull() && !v.IsUnknown() {
		return v
	}
	return types.StringValue(region)
}

func putMembers(ctx context.Context, lbClient *gophercloud.ServiceClient, poolID string, opts []octaviapools.BatchUpdateMemberOpts, timeout time.Duration) (*octaviapools.Pool, error) {
	parentPool, err := ipools.Get(lbClient, poolID).Extract()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve parent pool %s: %w", poolID, err)
	}

	if err := waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout); err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "Updating pool members", map[string]interface{}{"pool_id": poolID, "opts": fmt.Sprintf("%#v", opts)})
	delay := putMembersRetryMinDelay
	for {
		err := ipools.BatchUpdateMembers(lbClient, poolID, opts).ExtractErr()
		if err == nil {
			return parentPool, nil
		}
		if !errutil.Any(err, putMembersRetryableCodes) {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%w (last error: %w)", ctx.Err(), err)
		case <-time.After(delay):
		}
		delay = min(delay*2, putMembersRetryMaxDelay)
	}
}

func readMembers(ctx context.Context, lbClient *gophercloud.ServiceClient, poolID string, prior []memberModel, explicitEmpty explicitEmptySubnets) (types.Set, diag.Diagnostics, error) {
	apiMembers, err := listMembers(ctx, lbClient, poolID)
	if err != nil {
		return types.SetNull(memberObjectType), nil, err
	}

	tflog.Debug(ctx, "Retrieved pool members", map[string]interface{}{"pool_id": poolID, "members": fmt.Sprintf("%#v", apiMembers)})

	members, diags := types.SetValueFrom(ctx, memberObjectType, flattenMembers(apiMembers, prior, explicitEmpty))
	return members, diags, nil
}

func listMembers(ctx context.Context, lbClient *gophercloud.ServiceClient, poolID string) ([]octaviapools.Member, error) {
	allPages, err := octaviapools.ListMembers(lbClient, poolID, octaviapools.ListMembersOpts{}).AllPages()
	if err != nil {
		return nil, err
	}
	return octaviapools.ExtractMembers(allPages)
}

func withRequestContext(ctx context.Context, sc *gophercloud.ServiceClient) *gophercloud.ServiceClient {
	base := sc.ProviderClient
	pc := &gophercloud.ProviderClient{
		IdentityBase:      base.IdentityBase,
		IdentityEndpoint:  base.IdentityEndpoint,
		EndpointLocator:   base.EndpointLocator,
		HTTPClient:        base.HTTPClient,
		UserAgent:         base.UserAgent,
		Context:           ctx,
		RetryBackoffFunc:  base.RetryBackoffFunc,
		MaxBackoffRetries: base.MaxBackoffRetries,
		RetryFunc:         base.RetryFunc,
	}
	pc.CopyTokenFrom(base)
	if base.ReauthFunc != nil {
		pc.ReauthFunc = func() error {
			if err := ctx.Err(); err != nil {
				return err
			}
			prev := pc.Token()
			done := make(chan error, 1)
			go func() { done <- base.Reauthenticate(prev) }()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-done:
				if err != nil {
					return err
				}
			}
			pc.CopyTokenFrom(base)
			return nil
		}
	}

	scCopy := *sc
	scCopy.ProviderClient = pc
	return &scCopy
}

func membersFromSet(ctx context.Context, s types.Set) ([]memberModel, diag.Diagnostics) {
	var members []memberModel
	if s.IsNull() || s.IsUnknown() {
		return members, nil
	}
	diags := s.ElementsAs(ctx, &members, false)
	return members, diags
}

func expandMembers(members []memberModel) []octaviapools.BatchUpdateMemberOpts {
	opts := make([]octaviapools.BatchUpdateMemberOpts, 0, len(members))
	for _, m := range members {
		name := m.Name.ValueString()
		subnetID := m.SubnetID.ValueString()
		weight := int(m.Weight.ValueInt64())
		adminStateUp := m.AdminStateUp.ValueBool()
		backup := m.Backup.ValueBool()

		opts = append(opts, octaviapools.BatchUpdateMemberOpts{
			Address:      m.Address.ValueString(),
			ProtocolPort: int(m.ProtocolPort.ValueInt64()),
			Name:         &name,
			SubnetID:     &subnetID,
			Weight:       &weight,
			AdminStateUp: &adminStateUp,
			Backup:       &backup,
		})
	}
	return opts
}

func flattenMembers(apiMembers []octaviapools.Member, prior []memberModel, explicitEmpty explicitEmptySubnets) []memberModel {
	priorByKey := make(map[memberKey]memberModel, len(prior))
	for _, m := range prior {
		if k, ok := m.key(); ok {
			priorByKey[k] = m
		}
	}

	members := make([]memberModel, 0, len(apiMembers))
	for _, am := range apiMembers {
		m := memberModel{
			ID:           types.StringValue(am.ID),
			Name:         types.StringValue(am.Name),
			Address:      types.StringValue(am.Address),
			ProtocolPort: types.Int64Value(int64(am.ProtocolPort)),
			Weight:       types.Int64Value(int64(am.Weight)),
			SubnetID:     types.StringValue(am.SubnetID),
			Backup:       types.BoolValue(am.Backup),
			AdminStateUp: types.BoolValue(am.AdminStateUp),
		}

		k := memberKey{address: normalizeAddress(am.Address), protocolPort: int64(am.ProtocolPort)}
		if p, ok := priorByKey[k]; ok {
			m.Address = p.Address
		}
		if explicitEmpty[k] {
			m.SubnetID = types.StringValue("")
		}

		members = append(members, m)
	}
	return members
}
