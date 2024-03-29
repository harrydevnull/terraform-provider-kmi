package provider

import (
	"context"
	"encoding/xml"
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
				Required:    true,
				Description: "The group name of the admins who will manage the collection permissions. This can be set to the KMI account admin group. ",
			},
			"modifiers": schema.StringAttribute{
				Required:    true,
				Description: "The group name of the admins who will manage the collection permissions. This can be set to the KMI account admin group. ",
			},
			"readers": schema.StringAttribute{
				Required:    true,
				Description: "The group name of the admins who will read the collection  ",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the collection to create. ",
			},
			"account_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the account that KMI has been enabled for. ",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"distributed_date": schema.StringAttribute{
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
	_, err := r.client.GetGroup(plan.Readers.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Reader Group on Collection Create",
			"Could not get reader group, unexpected error: "+err.Error(),
		)
		return
	}

	collection := kmi.CollectionRequest{
		Adders:    plan.Adders.ValueString(),
		Modifiers: plan.Modifiers.ValueString(),
		Readers:   plan.Readers.ValueString(),
	}

	out, err := xml.MarshalIndent(collection, "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Marshalling Collection",
			"Could not marshal collection, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating Collection: %s", string(out)))
	err = r.client.CreateCollection(plan.AccountName.ValueString(), plan.CollectionName.ValueString(), collection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Collection",
			"Could not create Collection, unexpected error: "+err.Error(),
		)
		return
	}

	var duration_Minute time.Duration = 2 * time.Minute

	// response, err = r.client.GetCollection(plan.CollectionName.ValueString())
	response, err := retry(5, duration_Minute, func() (*kmi.Collection, error) { return r.client.GetCollection(plan.CollectionName.ValueString()) })

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Collection",
			"Could not get collection, unexpected error: "+err.Error(),
		)
		return
	}

	plan.DistributedDate = types.StringValue(response.DistributedDate)
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

	// Get refreshed order value from KMI
	kmicollection, err := r.client.GetCollection(state.CollectionName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Collection",
			"Could not read Collection "+state.CollectionName.ValueString()+": "+err.Error(),
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
	Adders          types.String `tfsdk:"adders"`
	Modifiers       types.String `tfsdk:"modifiers"`
	Readers         types.String `tfsdk:"readers"`
	CollectionName  types.String `tfsdk:"name"`
	AccountName     types.String `tfsdk:"account_name"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	DistributedDate types.String `tfsdk:"distributed_date"`
}
