package provider

import (
	"context"
	"fmt"
	"terraform-provider-kmi/internal/kmi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &collectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &collectionsDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewCollectionsDataSource() datasource.DataSource {
	return &collectionsDataSource{}
}

// collectionsDataSource is the data source implementation.
type collectionsDataSource struct {
	client *kmi.KMIRestClient
}

// Metadata returns the data source type name.
func (d *collectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collections"
}

// Schema defines the schema for the data source.
func (d *collectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"adders": schema.StringAttribute{
				Computed:    true,
				Description: "The group name of the admins who will manage the collection permissions. This can be set to the KMI account admin group. ",
			},
			"modifiers": schema.StringAttribute{
				Computed:    true,
				Description: "The group name of the admins who will manage the collection permissions. This can be set to the KMI account admin group. ",
			},
			"readers": schema.StringAttribute{
				Computed:    true,
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
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *collectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state collectionResourceModel

	// Get refreshed order value from KMI
	kmicollection, err := d.client.GetCollection(state.CollectionName.ValueString())
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

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *collectionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}
