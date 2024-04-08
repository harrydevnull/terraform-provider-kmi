package provider

import (
	"bytes"
	"context"
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
			"adders": schema.StringAttribute{
				Optional:    true,
				Description: "The group name of the admins who will manage the definition permissions. This can be set to the KMI account admin group. ",
			},
			"modifiers": schema.StringAttribute{
				Optional:    true,
				Description: "The group name of the admins who will manage the definition permissions. This can be set to the KMI account admin group. ",
			},
			"readers": schema.StringAttribute{
				Optional:    true,
				Description: "The group name of the admins who will read the definition  ",
			},
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
					"expire_period": schema.StringAttribute{
						Optional:    true,
						Description: "The expire period for the symmetric key. ",
					},
					"refresh_period": schema.StringAttribute{
						Optional:    true,
						Description: "The refresh period for the symmetric key. ",
					},
					"issuer": schema.StringAttribute{
						Optional:    true,
						Description: "The issuer for the SSL certificate. ",
					},
					// Error from KMI side: is_ca must not exist or have an integral value
					// is_ca has to be 1 if enabled, cannot be a true/false boolean
					"is_ca": schema.Int64Attribute{
						Optional:    true,
						Description: "Is the SSL certificate a CA. ",
					},
					"cn": schema.StringAttribute{
						Optional:    true,
						Description: "Common Name of the SSL certificate. ",
					},
					"subj_alt_names": schema.StringAttribute{
						Optional:    true,
						Description: "Subject Alternative Names of the SSL certificate. ",
					},
					"ca_name": schema.StringAttribute{
						Optional:    true,
						Description: "KMI path to the template used to sign the certificate by the CA.",
					},
					"signacl": schema.StringAttribute{
						Optional:    true,
						Description: "Collection that is eligible to sign the certificate. Can be used for CA definition setup.",
					},
					"signacldomain": schema.StringAttribute{
						Optional:    true,
						Description: "Much like a signacl rule, it restricts signing to the named collection. However, it has the additional restriction of only applying to a particular domain name or wildcarded domain (denoted by a domain starting with '*.' ). Can be used for CA definition setup.",
					},
					"signaclgroup": schema.StringAttribute{
						Optional:    true,
						Description: "Group that is eligible to sign the certificate. Can be used for CA definition setup.",
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
			"b64encoded": schema.BoolAttribute{
				Optional:    true,
				Description: "Should the secret be Base64-encoded? If it's not set, then is \"false\"",
			},
			"transparent": schema.StringAttribute{
				Optional:    true,
				Description: "The Transparent definition to create. ",
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

	definition := kmi.KMIDefinition{
		Adders:    plan.Adders.ValueString(),
		Modifiers: plan.Modifiers.ValueString(),
		Readers:   plan.Readers.ValueString(),
	}

	var err error
	if plan.SSLCert != nil {
		tflog.Info(ctx, "SSl cert is not nil")
		if !plan.SSLCert.Issuer.IsNull() && !plan.SSLCert.IsCA.IsNull() {
			resp.Diagnostics.AddError(
				"Validation error",
				"IsCA should not be set if Issuer is set ",
			)
			return

		}

		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, plan.SSLCert)
	}
	if plan.SymetricKey != nil {
		tflog.Info(ctx, "Symetric key is not nil")
		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, plan.SymetricKey)
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
		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, plan.AzureSP)
	}

	if !plan.Opaque.IsNull() {
		tflog.Info(ctx, "Opaque is not nil")
		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, Opaque{})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Definition",
				"Could not create Definition, unexpected error: "+err.Error(),
			)
			return
		}
		err = r.client.CreateBlockSecret(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), kmi.BlockSecret{
			Block: struct {
				Text       string "xml:\",chardata\""
				Name       string "xml:\"name,attr\""
				B64Encoded string `xml:"b64encoded,attr"`
			}{
				Name:       "opaque",
				Text:       plan.Opaque.ValueString(),
				B64Encoded: boolStr(plan.B64Encoded.ValueBool()),
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
	if !plan.Transparent.IsNull() {
		tflog.Info(ctx, "Transparent is not nil")
		transparent := Transparent{}

		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, transparent)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Definition",
				"Could not create Definition, unexpected error: "+err.Error(),
			)
			return
		}
		err = r.client.CreateBlockSecret(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), kmi.BlockSecret{
			Block: struct {
				Text       string "xml:\",chardata\""
				Name       string "xml:\"name,attr\""
				B64Encoded string `xml:"b64encoded,attr"`
			}{
				Name:       "transparent",
				Text:       plan.Transparent.ValueString(),
				B64Encoded: boolStr(plan.B64Encoded.ValueBool()),
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Transparent Block",
				"Could not create Transparent Block, unexpected error: "+err.Error(),
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

	definition := kmi.KMIDefinition{
		Adders:    plan.Adders.ValueString(),
		Modifiers: plan.Modifiers.ValueString(),
		Readers:   plan.Readers.ValueString(),
	}

	var err error
	if plan.SSLCert != nil {
		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, plan.SSLCert)
	}
	if plan.SymetricKey != nil {
		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, plan.SymetricKey)
	}
	if plan.AzureSP != nil {
		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, plan.AzureSP)
	}

	if !plan.Opaque.IsNull() {
		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, Opaque{})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Definition",
				"Could not create Definition, unexpected error: "+err.Error(),
			)
			return
		}
		err = r.client.CreateBlockSecret(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), kmi.BlockSecret{
			Block: struct {
				Text       string "xml:\",chardata\""
				Name       string "xml:\"name,attr\""
				B64Encoded string `xml:"b64encoded,attr"`
			}{
				Name:       "opaque",
				Text:       plan.Opaque.ValueString(),
				B64Encoded: boolStr(plan.B64Encoded.ValueBool()),
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

	if !plan.Transparent.IsNull() {
		tflog.Info(ctx, "Transparent is not nil")
		transparent := Transparent{}

		err = r.createDefinition(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), definition, transparent)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Definition",
				"Could not create Definition, unexpected error: "+err.Error(),
			)
			return
		}
		err = r.client.CreateBlockSecret(plan.CollectionName.ValueString(), plan.DefinitionName.ValueString(), kmi.BlockSecret{
			Block: struct {
				Text       string "xml:\",chardata\""
				Name       string "xml:\"name,attr\""
				B64Encoded string `xml:"b64encoded,attr"`
			}{
				Name:       "transparent",
				Text:       plan.Transparent.ValueString(),
				B64Encoded: boolStr(plan.B64Encoded.ValueBool()),
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Transparent Block",
				"Could not create Transparent Block, unexpected error: "+err.Error(),
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
	Adders         types.String `tfsdk:"adders"`
	Modifiers      types.String `tfsdk:"modifiers"`
	Readers        types.String `tfsdk:"readers"`
	DefinitionName types.String `tfsdk:"name"`
	CollectionName types.String `tfsdk:"collection_name"`
	LastUpdated    types.String `tfsdk:"last_updated"`
	SSLCert        *SSLCert     `tfsdk:"ssl_cert"`
	AzureSP        *AzureSP     `tfsdk:"azure_sp"`
	Opaque         types.String `tfsdk:"opaque"`
	B64Encoded     types.Bool   `tfsdk:"b64encoded"`
	Transparent    types.String `tfsdk:"transparent"`
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

type kmigenerator interface {
	RequestPayload(kmi.KMIDefinition) (kmi.KMIDefinition, error)
}

func (r *definitionsResource) createDefinition(collectionName string, definitionName string, definition kmi.KMIDefinition, kmigenerator kmigenerator) error {
	out, err := kmigenerator.RequestPayload(definition)
	fmt.Printf("CreateDefinition payload %v\n", out)
	if err != nil {
		return err
	}
	return r.client.CreateDefinition(collectionName, definitionName, out)
}

type Opaque struct {
}

func (op Opaque) RequestPayload(definition kmi.KMIDefinition) (kmi.KMIDefinition, error) {
	definition.Type = "opaque"
	return definition, nil
}

type Transparent struct {
}

func (op Transparent) RequestPayload(definition kmi.KMIDefinition) (kmi.KMIDefinition, error) {
	definition.Type = "transparent"
	return definition, nil
}

type SSLCert struct {
	AutoGenerate  types.Bool   `tfsdk:"auto_generate"`
	ExpiryPeriod  types.String `tfsdk:"expire_period"`
	RefreshPeriod types.String `tfsdk:"refresh_period"`
	Issuer        types.String `tfsdk:"issuer"`
	IsCA          types.Int64  `tfsdk:"is_ca"`
	Cn            types.String `tfsdk:"cn"`
	Sans          types.String `tfsdk:"subj_alt_names"`
	CAName        types.String `tfsdk:"ca_name"`
	SignACL       types.String `tfsdk:"signacl"`
	SignACLDomain types.String `tfsdk:"signacldomain"`
	SignACLGroup  types.String `tfsdk:"signaclgroup"`
}

func (s SSLCert) RequestPayload(definition kmi.KMIDefinition) (kmi.KMIDefinition, error) {

	if !s.Issuer.IsNull() && !s.IsCA.IsNull() {
		return kmi.KMIDefinition{}, fmt.Errorf("IsCA should not be set if Issuer is set ")
	}

	var options []*kmi.KMIOption
	if !s.IsCA.IsNull() {
		option := &kmi.KMIOption{
			Name: "is_ca",
			Text: fmt.Sprintf("%d", s.IsCA.ValueInt64()),
		}
		options = append(options, option)
	}
	if !s.Issuer.IsNull() {
		option := &kmi.KMIOption{
			Name: "issuer",
			Text: s.Issuer.ValueString(),
		}
		options = append(options, option)
	}
	if !s.Cn.IsNull() {
		option := &kmi.KMIOption{
			Name: "cn",
			Text: s.Cn.ValueString(),
		}
		options = append(options, option)
	}
	if !s.Sans.IsNull() {
		option := &kmi.KMIOption{
			Name: "subj_alt_names",
			Text: s.Sans.ValueString(),
		}
		options = append(options, option)
	}
	if !s.SignACL.IsNull() {
		option := &kmi.KMIOption{
			Name: "signacl:" + s.SignACL.ValueString(),
			Text: "true",
		}
		options = append(options, option)
	}
	if !s.SignACLDomain.IsNull() {
		option := &kmi.KMIOption{
			Name: "signacldomain:" + s.SignACLDomain.ValueString(),
			Text: "true",
		}
		options = append(options, option)
	}
	if !s.SignACLGroup.IsNull() {
		option := &kmi.KMIOption{
			Name: "signaclgroup:" + s.SignACLGroup.ValueString(),
			Text: "true",
		}
		options = append(options, option)
	}
	if !s.CAName.IsNull() {
		option := &kmi.KMIOption{
			Name: "ca_name",
			Text: s.CAName.ValueString(),
		}
		options = append(options, option)
	}

	definition.AutoGenerate = boolStr(s.AutoGenerate.ValueBool())
	definition.Type = "ssl_cert"
	definition.ExpirePeriod = s.ExpiryPeriod.ValueString()
	definition.RefreshPeriod = s.RefreshPeriod.ValueString()
	definition.Options = options
	return definition, nil
}

type AzureSP struct {
	AutoGenerate types.Bool `tfsdk:"auto_generate"`
}

func (sp AzureSP) RequestPayload(definition kmi.KMIDefinition) (kmi.KMIDefinition, error) {
	definition.AutoGenerate = boolStr(sp.AutoGenerate.ValueBool())
	definition.Type = "azure_sp"
	return definition, nil
}

type SymetricKey struct {
	AutoGenerate  types.Bool   `tfsdk:"auto_generate"`
	ExpiryPeriod  types.String `tfsdk:"expire_period"`
	RefreshPeriod types.String `tfsdk:"refresh_period"`
	KeySizeBytes  types.Int64  `tfsdk:"key_size_bytes"`
}

func (sk SymetricKey) RequestPayload(definition kmi.KMIDefinition) (kmi.KMIDefinition, error) {
	var options []*kmi.KMIOption
	if !sk.KeySizeBytes.IsNull() {
		option := &kmi.KMIOption{
			Name: "key_size_bytes",
			Text: fmt.Sprintf("%d", sk.KeySizeBytes.ValueInt64()),
		}
		options = append(options, option)
	}

	definition.AutoGenerate = boolStr(sk.AutoGenerate.ValueBool())
	definition.Type = "symmetric_key"
	definition.ExpirePeriod = sk.ExpiryPeriod.ValueString()
	definition.RefreshPeriod = sk.RefreshPeriod.ValueString()
	definition.Options = options
	return definition, nil
}
