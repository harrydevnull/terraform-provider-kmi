package provider

import (
	"context"
	"terraform-provider-kmi/internal/kmi"

	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &accountDataSource{}
	_ datasource.DataSourceWithConfigure = &accountDataSource{}
)

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

// accountDataSource is the data source implementation.
type accountDataSource struct {
	client *kmi.KMIRestClient
}

// Metadata returns the data source type name.
func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

// Schema defines the schema for the data source.
func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the account that KMI has been enabled for. ",
			},
			"engines": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the engine. ",
						},

						"cloud": schema.StringAttribute{
							Computed:    true,
							Description: "The cloud type of the engine.  azure, gcp, or linode ",
						},
						"adders": schema.StringAttribute{
							Computed:    true,
							Description: "The group name of the admins who will manage the engine permissions. This can be set to the KMI account admin group. ",
						},
						"modifiers": schema.StringAttribute{
							Computed:    true,
							Description: "The group name of the admins who will manage the engine permissions. This can be set to the KMI account admin group. ",
						},
						"modified": schema.Int64Attribute{
							Computed:    true,
							Description: "The last time the engine was modified. ",
						},
						"source": schema.StringAttribute{
							Computed: true,
						},
						"published": schema.StringAttribute{
							Computed: true,
						},
						"published_location": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"collections": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the collection. ",
						},
						"source": schema.StringAttribute{
							Computed: true,
						},
						"readers": schema.StringAttribute{
							Computed:    true,
							Description: "The group name of the admins who will read the collection  ",
						},
						"adders": schema.StringAttribute{
							Computed:    true,
							Description: "The group name of the admins who will manage the collection permissions. This can be set to the KMI account admin group. ",
						},
						"modifiers": schema.StringAttribute{
							Computed:    true,
							Description: "The group name of the admins who will manage the collection permissions. This can be set to the KMI account admin group. ",
						},
						"modified": schema.Int64Attribute{
							Computed:    true,
							Description: "The last time the collection was modified. ",
						},
						"distributed": schema.Int64Attribute{
							Computed:    true,
							Description: "The last time the collection was modified. ",
						},
						"distributed_date": schema.StringAttribute{
							Computed:    true,
							Description: "The last time the collection was modified. ",
						},
						"keyspace": schema.StringAttribute{
							Computed:    true,
							Description: "The last time the collection was modified. ",
						},
						"account": schema.StringAttribute{
							Computed:    true,
							Description: "The account the collection belongs to. ",
						},
					},
				},
			},
			"groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the group. ",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the group. ",
						},
						"source": schema.StringAttribute{
							Computed: true,
						},
						"account": schema.StringAttribute{
							Computed:    true,
							Description: "The account the group belongs to. ",
						},
						"engine": schema.StringAttribute{
							Computed:    true,
							Description: "The engine the group belongs to. ",
						},
						"projection": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the identity projection",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountDataSourceModel
	tflog.Debug(ctx, "Preparing to read item data source")

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	engines, err := d.client.GetAccountDetails(state.AccountName)
	ctx = tflog.SetField(ctx, "Read engines", engines)
	tflog.Debug(ctx, "Reading   AccountDataSource")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Engines under account",
			err.Error(),
		)
		return
	}

	for _, kmiEngine := range engines.Engine {
		tfEngineState := engineModel{
			Name:              types.StringValue(kmiEngine.Name),
			CloudType:         types.StringValue(kmiEngine.Cloud),
			Adders:            types.StringValue(kmiEngine.Adders),
			Modifiers:         types.StringValue(kmiEngine.Modifiers),
			Modified:          types.Int64Value(kmiEngine.Modified),
			Source:            types.StringValue(kmiEngine.Source),
			Published:         types.StringValue(kmiEngine.Published),
			PublishedLocation: types.StringValue(kmiEngine.PublishedLocation),
		}

		state.Engines = append(state.Engines, tfEngineState)

	}
	for _, kmicollection := range engines.Collection {
		collectionState := collectionModel{
			Name:            types.StringValue(kmicollection.Name),
			Source:          types.StringValue(kmicollection.Source),
			Readers:         types.StringValue(kmicollection.Readers),
			Adders:          types.StringValue(kmicollection.Adders),
			Modified:        types.Int64Value(kmicollection.Modified),
			Modifiers:       types.StringValue(kmicollection.Modifiers),
			Distributed:     types.Int64Value(kmicollection.Distributed),
			Keyspace:        types.StringValue(kmicollection.Keyspace),
			DistributedDate: types.StringValue(kmicollection.DistributedDate),
			Account:         types.StringValue(kmicollection.Account),
		}
		state.Collections = append(state.Collections, collectionState)
	}

	for _, kmiGroup := range engines.Group {
		groupState := groupModel{
			Name:       types.StringValue(kmiGroup.Name),
			Type:       types.StringValue(kmiGroup.Type),
			Source:     types.StringValue(kmiGroup.Source),
			Account:    types.StringValue(kmiGroup.Account),
			Engine:     types.StringValue(kmiGroup.Engine),
			Projection: types.StringValue(kmiGroup.Projection),
		}
		state.Groups = append(state.Groups, groupState)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *accountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	ctx = tflog.SetField(ctx, "ProviderData ", req.ProviderData)

	tflog.Debug(ctx, "Configuring  AccountDataSource")

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

// coffeesDataSourceModel maps the data source schema data.
type accountDataSourceModel struct {
	AccountName string            `tfsdk:"account_name"`
	Engines     []engineModel     `tfsdk:"engines"`
	Collections []collectionModel `tfsdk:"collections"`
	Groups      []groupModel      `tfsdk:"groups"`
}

// coffeesModel maps coffees schema data.
type engineModel struct {
	Name              types.String `tfsdk:"name"`
	CloudType         types.String `tfsdk:"cloud"`
	Adders            types.String `tfsdk:"adders"`
	Modifiers         types.String `tfsdk:"modifiers"`
	Modified          types.Int64  `tfsdk:"modified"`
	Source            types.String `tfsdk:"source"`
	Published         types.String `tfsdk:"published"`
	PublishedLocation types.String `tfsdk:"published_location"`
}
type collectionModel struct {
	Name            types.String `tfsdk:"name"`
	Source          types.String `tfsdk:"source"`
	Readers         types.String `tfsdk:"readers"`
	Adders          types.String `tfsdk:"adders"`
	Modifiers       types.String `tfsdk:"modifiers"`
	Modified        types.Int64  `tfsdk:"modified"`
	Distributed     types.Int64  `tfsdk:"distributed"`
	DistributedDate types.String `tfsdk:"distributed_date"`
	Keyspace        types.String `tfsdk:"keyspace"`
	Account         types.String `tfsdk:"account"`
}

type groupModel struct {
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Source     types.String `tfsdk:"source"`
	Account    types.String `tfsdk:"account"`
	Engine     types.String `tfsdk:"engine"`
	Projection types.String `tfsdk:"projection"`
}
