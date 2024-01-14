package provider

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"regexp"
	"terraform-provider-kmi/internal/kmi"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				Required:    true,
				Description: "The name of the definition to create. ",
			},
			"collection_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the collection to create. ",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"ssl_cert": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"auto_generate": schema.BoolAttribute{
						Required:    true,
						Description: "Auto generate the SSL certificate. ",
					},
				},
				Optional:    true,
				Description: "The SSL certificate to create. ",
			},
			"azure_sp": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"auto_generate": schema.BoolAttribute{
						Required: true,
					},
				},
				Optional:    true,
				Description: "The Azure Service Principal to create. ",
			},
			"options": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of options for the definition. ",
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
				Optional:    true,
				Description: "The Opaque definition to create. ",
			},
			"secret_indexes": schema.StringAttribute{
				Computed:    true,
				Description: "The list of secret indexes for the definition. ",
			},
			"symmetric_key": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"auto_generate": schema.BoolAttribute{
						Required:    true,
						Description: "Auto generate the symmetric key. ",
					},
					"expire_period": schema.StringAttribute{
						Required:    true,
						Description: "The expire period for the symmetric key. ",
					},
					"refresh_period": schema.StringAttribute{
						Required:    true,
						Description: "The refresh period for the symmetric key. ",
					},
					"key_size_bytes": schema.Int64Attribute{
						Optional:    true,
						Description: "The key size in bytes for the symmetric key. ",
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
	if plan.SSLCert != nil {
		tflog.Info(ctx, "SSl cert is not nil")

		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.SSLCert)
	}
	if plan.SymetricKey != nil {
		tflog.Info(ctx, "Symetric key is not nil")
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.SymetricKey)
	}
	if plan.AzureSP != nil {
		tflog.Info(ctx, "Azure SP is not nil")

		if !regexp.MustCompile(`^[a-z0-9_\-]+$`).MatchString(plan.CollectionName.ValueString()) {
			resp.Diagnostics.AddError(
				"Validation error",
				"Collection name should only contain lower case for Azure SP  ",
			)
			return
		}
		if !regexp.MustCompile(`^[a-z0-9_\-]+$`).MatchString(plan.DefinitionName.ValueString()) {
			resp.Diagnostics.AddError(
				"Validation error",
				"DefinitionName should only contain lower case for Azure SP  ",
			)
			return
		}
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.AzureSP)
	}

	if !plan.Opaque.IsNull() {
		op := Opaque{}
		tflog.Info(ctx, "Opaque is not nil")
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), op)
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
	options := []DefinitionOption{}

	for _, optionfromKmi := range definitionDetails.Option {

		options = append(options, DefinitionOption{
			Name:  types.StringValue(optionfromKmi.Name),
			Value: types.StringValue(optionfromKmi.Text),
		})
	}

	plan.Options = keySliceToList(ctx, options, &resp.Diagnostics)
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
	options := []DefinitionOption{}

	for _, optionfromKmi := range definitionDetails.Option {

		options = append(options, DefinitionOption{
			Name:  types.StringValue(optionfromKmi.Name),
			Value: types.StringValue(optionfromKmi.Text),
		})
	}
	state.Options = keySliceToList(ctx, options, &resp.Diagnostics)
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
	if plan.SSLCert != nil {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.SSLCert)
	}
	if plan.SymetricKey != nil {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.SymetricKey)
	}
	if plan.AzureSP != nil {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), plan.AzureSP)
	}

	if !plan.Opaque.IsNull() {
		err = r.client.CreateDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), Opaque{})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Definition",
				"Could not create Definition, unexpected error: "+err.Error(),
			)
			return
		}
		err = r.client.CreateOpaueSecret(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), kmi.OpaqueSecret{
			Block: struct {
				Text string "xml:\",chardata\""
				Name string "xml:\"name,attr\""
			}{
				Name: "opaque",
				Text: plan.Opaque.ValueString(),
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Opaque Secret",
				"Could not create Opaue Secret, unexpected error: "+err.Error(),
			)
			return
		}
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
	options := []DefinitionOption{}

	for _, optionfromKmi := range definitionDetails.Option {

		options = append(options, DefinitionOption{
			Name:  types.StringValue(optionfromKmi.Name),
			Value: types.StringValue(optionfromKmi.Text),
		})
	}

	plan.Options = keySliceToList(ctx, options, &resp.Diagnostics)
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

// I feel like throwing up writing this function.
func boolStr(s bool) string {
	if s {
		return "True"
	} else {
		return "False"
	}
}

type definitionResourceModel struct {
	DefinitionName types.String `tfsdk:"name"`
	CollectionName types.String `tfsdk:"collection_name"`
	LastUpdated    types.String `tfsdk:"last_updated"`
	SSLCert        *SSLCert     `tfsdk:"ssl_cert"`
	AzureSP        *AzureSP     `tfsdk:"azure_sp"`
	Opaque         types.String `tfsdk:"opaque"`
	SymetricKey    *SymetricKey `tfsdk:"symmetric_key"`
	Options        types.List   `tfsdk:"options"`
	SecretIndexes  types.String `tfsdk:"secret_indexes"`
}

type DefinitionOption struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (o DefinitionOption) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}
}

func keySliceToList(ctx context.Context, keysSliceIn []DefinitionOption, diags *diag.Diagnostics) types.List {
	keys, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: DefinitionOption{}.attrTypes()}, keysSliceIn)
	diags.Append(d...)
	return keys
}

type Opaque struct {
}

func (op Opaque) RequestPayload() ([]byte, error) {
	defn := kmi.KMIDefinition{
		Type: "opaque",
	}
	return xml.MarshalIndent(defn, "", "")
	//return xml.MarshalIndent(defn, " ", "  ")
}

type SSLCert struct {
	AutoGenerate types.Bool `tfsdk:"auto_generate"`
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

type SymetricKey struct {
	AutoGenerate  types.Bool   `tfsdk:"auto_generate"`
	ExpiryPeriod  types.String `tfsdk:"expire_period"`
	RefreshPeriod types.String `tfsdk:"refresh_period"`
	KeySizeBytes  types.Int64  `tfsdk:"key_size_bytes"`
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
