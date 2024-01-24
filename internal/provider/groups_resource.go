package provider

import (
	"context"
	"fmt"
	"terraform-provider-kmi/internal/kmi"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &groupsResource{}
	_ resource.ResourceWithConfigure = &groupsResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewGroupsResource() resource.Resource {
	return &groupsResource{}
}

// groupsResource is the resource implementation.
type groupsResource struct {
	client *kmi.KMIRestClient
}

// Metadata returns the resource type name.
func (r *groupsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema defines the schema for the resource.
func (r *groupsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"group_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the group to create. ",
			},

			"account_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the account that KMI has been enabled for. ",
			},
			"adders": schema.StringAttribute{
				Computed:    true,
				Description: "The list of adders for the group. ",
			},
			"modifiers": schema.StringAttribute{
				Computed:    true,
				Description: "The list of modifiers for the group. ",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The last time the group was updated. ",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *groupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Create Groups Request")

	err := r.client.CreateGroup(plan.AccountName.ValueString(), plan.GroupName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Group",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	groupInfo, err := r.client.GetGroup(plan.GroupName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Groups",
			"Could not read Groups "+plan.GroupName.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Adders = types.StringValue(groupInfo.Adders)
	plan.Modifiers = types.StringValue(groupInfo.Modifiers)
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *groupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from KMI Group
	groupInfo, err := r.client.GetGroup(state.GroupName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Groups",
			"Could not read Groups "+state.GroupName.ValueString()+": "+err.Error(),
		)
		return
	}

	state.AccountName = types.StringValue(groupInfo.Account)
	state.Adders = types.StringValue(groupInfo.Adders)
	state.Modifiers = types.StringValue(groupInfo.Modifiers)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateGroup(plan.AccountName.ValueString(), plan.GroupName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Group",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	groupInfo, err := r.client.GetGroup(plan.GroupName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Groups",
			"Could not read Groups "+plan.GroupName.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Adders = types.StringValue(groupInfo.Adders)
	plan.Modifiers = types.StringValue(groupInfo.Modifiers)
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *groupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteGroup(state.GroupName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting KMI Group",
			"Could not group, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *groupsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kmi.KMIRestClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *KMIRestClient., got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

type groupResourceModel struct {
	AccountName types.String `tfsdk:"account_name"`
	GroupName   types.String `tfsdk:"group_name"`
	Adders      types.String `tfsdk:"adders"`
	Modifiers   types.String `tfsdk:"modifiers"`
	LastUpdated types.String `tfsdk:"last_updated"`
}
