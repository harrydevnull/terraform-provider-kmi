package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-kmi/internal/kmi"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &groupsMembershipResource{}
	_ resource.ResourceWithConfigure = &groupsMembershipResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewGroupsMembershipResource() resource.Resource {
	return &groupsMembershipResource{}
}

// groupsMembershipResource is the resource implementation.
type groupsMembershipResource struct {
	client *kmi.KMIRestClient
}

// Metadata returns the resource type name.
func (r *groupsMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_membership"
}

// Schema defines the schema for the resource.
func (r *groupsMembershipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"group_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the group to create. ",
			},

			"members": schema.SetAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *groupsMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupsMembershipModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Create Groups Request")

	var errstrings []string

	elements := make([]types.String, 0, len(plan.Members.Elements()))
	diags = plan.Members.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return
	}
	for _, member := range elements {
		err := r.client.CreateGroupMembership(plan.GroupName.ValueString(), member.ValueString())
		if err != nil {
			errstrings = append(errstrings, err.Error())
		}
	}

	if errstrings != nil {
		resp.Diagnostics.AddError(
			"Error creating Group Membership",
			"Could not create Group Membership, unexpected error: "+strings.Join(errstrings, "\n"),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *groupsMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupsMembershipModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from HashiCups
	_, err := r.client.GetGroup(state.GroupName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Groups",
			"Could not read Groups "+state.GroupName.ValueString()+": "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupsMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan groupsMembershipModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var errstrings []string

	elements := make([]types.String, 0, len(plan.Members.Elements()))
	diags = plan.Members.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return
	}
	for _, member := range elements {
		err := r.client.CreateGroupMembership(plan.GroupName.ValueString(), member.ValueString())
		if err != nil {
			errstrings = append(errstrings, err.Error())
		}
	}
	if errstrings != nil {
		resp.Diagnostics.AddError(
			"Error creating Group",
			"Could not create Group, unexpected error: "+strings.Join(errstrings, "\n"),
		)
		return
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *groupsMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state groupsMembershipModel
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

func (r *groupsMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type groupsMembershipModel struct {
	GroupName   types.String `tfsdk:"group_name"`
	Members     types.List   `tfsdk:"members"`
	LastUpdated types.String `tfsdk:"last_updated"`
}
