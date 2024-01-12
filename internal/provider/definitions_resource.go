package provider

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"reflect"
	"terraform-provider-kmi/internal/kmi"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &definitionsResource{}
	_ resource.ResourceWithConfigure = &definitionsResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDefinitionsResource() resource.Resource {
	return &definitionsResource{}
}

// definitionsResource is the resource implementation.
type definitionsResource struct {
	client *kmi.KMIRestClient
}

// Metadata returns the resource type name.
func (r *definitionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_definitions"
}

// Schema defines the schema for the resource.
// Schema defines the schema for the resource.
func (r *definitionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"collection_name": schema.StringAttribute{
				Required: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"ssl_cert": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"auto_generate": schema.BoolAttribute{
						Required: true,
					},
				},
				Optional: true,
			},
			"azure_sp": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"auto_generate": schema.BoolAttribute{
						Required: true,
					},
				},
				Optional: true,
			},
			"option": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"value": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"opaque": schema.StringAttribute{

				Optional: true,
			},
			"secret_indexes": schema.StringAttribute{

				Computed: true,
			},
			"symmetric_key": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"auto_generate": schema.BoolAttribute{
						Required: true,
					},
					"expire_period": schema.Int64Attribute{
						Required: true,
					},
					"refresh_period": schema.Int64Attribute{
						Required: true,
					},
					"key_size_bytes": schema.Int64Attribute{
						Optional: true,
					},
				},
				Optional: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *definitionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan definitionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !plan.SSLCert.IsEmpty() {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.SSLCert)
	}
	if !plan.SymetricKey.IsEmpty() {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.SymetricKey)
	}
	if !plan.AzureSP.IsEmpty() {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.AzureSP)
	}

	if !plan.Opaque.IsEmpty() {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.Opaque)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Definition",
			"Could not create Definition, unexpected error: "+err.Error(),
		)
		return
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	definitionDetails, err := r.client.GetDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading definitions details",
			"Could not read definitions "+plan.DefinitionName.ValueString()+": "+err.Error(),
		)
		return
	}
	plan.Options = []DefinitionOption{}

	for _, optionfromKmi := range definitionDetails.Option {

		plan.Options = append(plan.Options, DefinitionOption{
			Name:  types.StringValue(optionfromKmi.Name),
			Value: types.StringValue(optionfromKmi.Text),
		})
	}
	var secretsIndex bytes.Buffer
	for _, secret := range definitionDetails.Secret {
		secretsIndex.WriteString(fmt.Sprintf("%s,", secret.Index))
	}
	plan.SecretIndexes = types.StringValue(secretsIndex.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *definitionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state definitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	definitionDetails, err := r.client.GetDefinition(state.CollectionName.ValueString(), state.DefinitionName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading definitions details",
			"Could not read definitions "+state.DefinitionName.ValueString()+": "+err.Error(),
		)
		return
	}
	state.Options = []DefinitionOption{}

	for _, optionfromKmi := range definitionDetails.Option {

		state.Options = append(state.Options, DefinitionOption{
			Name:  types.StringValue(optionfromKmi.Name),
			Value: types.StringValue(optionfromKmi.Text),
		})
	}
	var secretsIndex bytes.Buffer
	for _, secret := range definitionDetails.Secret {
		secretsIndex.WriteString(fmt.Sprintf("%s,", secret.Index))
	}
	state.SecretIndexes = types.StringValue(secretsIndex.String())
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *definitionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan definitionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !plan.SSLCert.IsEmpty() {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.SSLCert)
	}
	if !plan.SymetricKey.IsEmpty() {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.SymetricKey)
	}
	if !plan.AzureSP.IsEmpty() {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.AzureSP)
	}

	if !plan.Opaque.IsEmpty() {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.Opaque)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Definition",
			"Could not create Definition, unexpected error: "+err.Error(),
		)
		return
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	definitionDetails, err := r.client.GetDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading definitions details",
			"Could not read definitions "+plan.DefinitionName.ValueString()+": "+err.Error(),
		)
		return
	}
	plan.Options = []DefinitionOption{}

	for _, optionfromKmi := range definitionDetails.Option {

		plan.Options = append(plan.Options, DefinitionOption{
			Name:  types.StringValue(optionfromKmi.Name),
			Value: types.StringValue(optionfromKmi.Text),
		})
	}
	var secretsIndex bytes.Buffer
	for _, secret := range definitionDetails.Secret {
		secretsIndex.WriteString(fmt.Sprintf("%s,", secret.Index))
	}
	plan.SecretIndexes = types.StringValue(secretsIndex.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *definitionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state definitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDefinition(state.CollectionName.ValueString(), state.DefinitionName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Definitions",
			"Could not delete Definitions, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *definitionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// I feel like throwing up writing this function
func boolStr(s bool) string {
	if s {
		return "True"
	} else {
		return "False"
	}
}

type definitionResourceModel struct {
	DefinitionName types.String       `tfsdk:"name"`
	CollectionName types.String       `tfsdk:"collection_name"`
	LastUpdated    types.String       `tfsdk:"last_updated"`
	SSLCert        SSLCert            `tfsdk:"ssl_cert"`
	AzureSP        AzureSP            `tfsdk:"azure_sp"`
	Opaque         Opaque             `tfsdk:"opaque"`
	SymetricKey    SymetricKey        `tfsdk:"symmetric_key"`
	Options        []DefinitionOption `tfsdk:"option"`
	SecretIndexes  types.String       `tfsdk:"secret_indexes"`
}

type DefinitionOption struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
type Opaque struct {
}

func (op Opaque) IsEmpty() bool {
	return reflect.DeepEqual(op, Opaque{})
}

func (op Opaque) RequestPayload() ([]byte, error) {
	defn := kmi.KMIDefinition{
		Type: "opaque",
	}
	return xml.MarshalIndent(defn, " ", "  ")
}

type SSLCert struct {
	AutoGenerate types.Bool `tfsdk:"auto_generate"`
}

func (s SSLCert) IsEmpty() bool {
	return reflect.DeepEqual(s, SSLCert{})
}

func (s SSLCert) RequestPayload() ([]byte, error) {

	defn := kmi.KMIDefinition{
		AutoGenerate: boolStr(s.AutoGenerate.ValueBool()),
		Type:         "ssl_cert",
	}
	return xml.MarshalIndent(defn, " ", "  ")

}

type AzureSP struct {
	AutoGenerate types.Bool `tfsdk:"auto_generate"`
}

func (sp AzureSP) RequestPayload() ([]byte, error) {
	defn := kmi.KMIDefinition{
		AutoGenerate: boolStr(sp.AutoGenerate.ValueBool()),
		Type:         "azure_sp",
	}
	return xml.MarshalIndent(defn, " ", "  ")

}

func (s AzureSP) IsEmpty() bool {
	return reflect.DeepEqual(s, AzureSP{})
}

type SymetricKey struct {
	AutoGenerate  types.Bool   `tfsdk:"auto_generate"`
	ExpiryPeriod  types.String `tfsdk:"expire_period"`
	RefreshPeriod types.String `tfsdk:"refresh_period"`
	KeySizeBytes  types.Int64  `tfsdk:"key_size_bytes"`
}

func (s SymetricKey) IsEmpty() bool {
	return reflect.DeepEqual(s, SymetricKey{})
}

func (sk SymetricKey) RequestPayload() ([]byte, error) {
	defn := kmi.KMIDefinition{
		AutoGenerate:  boolStr(sk.AutoGenerate.ValueBool()),
		Type:          "symmetric_key",
		ExpirePeriod:  sk.ExpiryPeriod.ValueString(),
		RefreshPeriod: sk.RefreshPeriod.ValueString(),
		Option: &kmi.KMIOption{
			Name: "key_size_bytes",
			Text: fmt.Sprintf("%d", sk.KeySizeBytes.ValueInt64()),
		},
	}
	return xml.MarshalIndent(defn, " ", "  ")

}
