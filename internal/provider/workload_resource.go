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
	_ resource.Resource              = &workloadResource{}
	_ resource.ResourceWithConfigure = &workloadResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewWorkloadResource() resource.Resource {
	return &workloadResource{}
}

// workloadResource is the resource implementation.
type workloadResource struct {
	client *kmi.KMIRestClient
}

// Metadata returns the resource type name.
func (r *workloadResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workload"
}

// orderItemModel maps order item data.
type WorkloadVMResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Account     types.String `tfsdk:"account"`
	Engine      types.String `tfsdk:"engine"`
	Region      types.String `tfsdk:"region"`
	LinodeLabel types.String `tfsdk:"linode_label"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Schema defines the schema for the resource.
func (r *workloadResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{

			"name": schema.StringAttribute{
				Required: true,
			},
			"account": schema.StringAttribute{
				Required: true,
			},
			"engine": schema.StringAttribute{
				Required: true,
			},
			"region": schema.StringAttribute{
				Required: true,
			},
			"linode_label": schema.StringAttribute{
				Optional: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}

}

// Create creates the resource and sets the initial Terraform state.
func (r *workloadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkloadVMResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	kmiworkload := &kmi.Workload{
		Projection: plan.Name.ValueString(),
		Region: struct {
			Text   string "xml:\",chardata\""
			Source string "xml:\"source,attr,omitempty\""
		}{
			Text: plan.Region.ValueString(),
		},
		LinodeLabel: &kmi.LinodeLabel{
			Text: plan.LinodeLabel.ValueString(),
		},
	}
	tflog.Info(ctx, "Create workload payload %v\n")

	_, err := r.client.CreateWorkloadDetails(plan.Account.ValueString(), plan.Engine.ValueString(), *kmiworkload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workload",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	kmiworkloadfromservice, err := r.client.GetWorkloadDetails(plan.Account.ValueString(), plan.Engine.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workload",
			"Could not create workload, unexpected error: "+err.Error(),
		)
		return
	}
	plan.Name = types.StringValue(kmiworkloadfromservice.Projection)
	plan.Region = types.StringValue(kmiworkloadfromservice.Region.Text)
	plan.LinodeLabel = types.StringValue(kmiworkloadfromservice.LinodeLabel.Text)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *workloadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkloadVMResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	kmiworkloadfromservice, err := r.client.GetWorkloadDetails(state.Account.ValueString(), state.Engine.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting workload",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}
	state.Name = types.StringValue(kmiworkloadfromservice.Projection)
	state.Region = types.StringValue(kmiworkloadfromservice.Region.Text)
	state.LinodeLabel = types.StringValue(kmiworkloadfromservice.LinodeLabel.Text)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *workloadResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *workloadResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkloadVMResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteWorkload(state.Account.ValueString(), state.Engine.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting KMI workloadResource",
			"Could not group, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *workloadResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kmi.KMIRestClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *kmi.KMIRestClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
