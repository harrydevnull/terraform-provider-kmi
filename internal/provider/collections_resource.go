package provider

import (
	"context"
	"fmt"
	"terraform-provider-kmi/internal/kmi"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &collectionsResource{}
	_ resource.ResourceWithConfigure = &collectionsResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewCollectionsResource() resource.Resource {
	return &collectionsResource{}
}

// collectionsResource is the resource implementation.
type collectionsResource struct {
	client *kmi.KMIRestClient
}

// Metadata returns the resource type name.
func (r *collectionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collections"
}

// Schema defines the schema for the resource.
func (r *collectionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"adders": schema.StringAttribute{
				Required: true,
			},
			"modifiers": schema.StringAttribute{
				Required: true,
			},
			"readers": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"account_name": schema.StringAttribute{
				Required: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *collectionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan collectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	collection := kmi.CollectionRequest{
		Adders:    plan.Adders.ValueString(),
		Modifiers: plan.Modifiers.ValueString(),
		Readers:   plan.Readers.ValueString(),
	}

	err := r.client.CreateCollection(plan.AccountName.ValueString(), plan.CollectionName.ValueString(), collection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Collection",
			"Could not create Collection, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.client.GetCollection(plan.CollectionName.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Collection",
			"Could not get collection, unexpected error: "+err.Error(),
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
func (r *collectionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state collectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from HashiCups
	kmicollection, err := r.client.GetCollection(state.CollectionName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading HashiCups Order",
			"Could not read HashiCups order ID "+state.CollectionName.ValueString()+": "+err.Error(),
		)
		return
	}

	state = collectionResourceModel{
		Adders:         types.StringValue(kmicollection.Adders),
		Modifiers:      types.StringValue(kmicollection.Modifiers),
		Readers:        types.StringValue(kmicollection.Readers),
		CollectionName: types.StringValue(kmicollection.Name),
		AccountName:    types.StringValue(kmicollection.Account),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *collectionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan collectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	collection := kmi.CollectionRequest{
		Adders:    plan.Adders.ValueString(),
		Modifiers: plan.Modifiers.ValueString(),
		Readers:   plan.Readers.ValueString(),
	}

	err := r.client.CreateCollection(plan.AccountName.ValueString(), plan.CollectionName.ValueString(), collection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Collection",
			"Could not create Collection, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.client.GetCollection(plan.CollectionName.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Collection",
			"Could not get collection, unexpected error: "+err.Error(),
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
func (r *collectionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state collectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing collection
	err := r.client.DeleteCollection(state.CollectionName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Collections ",
			"Could not delete collections, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *collectionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kmi.KMIRestClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *KMIRestClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

type collectionResourceModel struct {
	Adders         types.String `tfsdk:"adders"`
	Modifiers      types.String `tfsdk:"modifiers"`
	Readers        types.String `tfsdk:"readers"`
	CollectionName types.String `tfsdk:"name"`
	AccountName    types.String `tfsdk:"account_name"`
	LastUpdated    types.String `tfsdk:"last_updated"`
}
