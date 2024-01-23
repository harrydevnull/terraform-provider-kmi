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
	_ resource.Resource              = &templateResource{}
	_ resource.ResourceWithConfigure = &templateResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewTemplateResource() resource.Resource {
	return &templateResource{}
}

// templateResource is the resource implementation.
type templateResource struct {
	client *kmi.KMIRestClient
}

// Metadata returns the resource type name.
func (r *templateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

// Schema defines the schema for the resource.
func (r *templateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ca_collection": schema.StringAttribute{
				Required:    true,
				Description: "CA collection name to be created on KMI ",
			},
			"ca_definition": schema.StringAttribute{
				Required:    true,
				Description: "CA definition name to be created on KMI ",
			},
			"template_name": schema.StringAttribute{
				Required:    true,
				Description: " Certificate Signing Request template name to be created on KMI ",
			},
			"client_collection": schema.StringAttribute{
				Required:    true,
				Description: " Client collection name to be created on KMI ",
			},

			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: " The last time the group was updated. ",
			},
			"options": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"min_ttl": schema.StringAttribute{
							Optional:    true,
							Description: " The minimum period of time that the signed secret can be valid for Default is 7 day ",
						},
						"max_ttl": schema.StringAttribute{
							Optional:    true,
							Description: " The maximum period of time that the signed secret can be valid for Default is 90 days",
						},
						"leaf_exceeds_ca_ttl": schema.StringAttribute{
							Optional:    true,
							Description: " Boolean flag as to whether or not the signed secret's notAfter date can exceed that of the CA certificate ",
						},
						"allow_ca": schema.StringAttribute{
							Optional:    true,
							Description: " Whether the signed secret can have the CA option set in the BasicConstraints extension.",
						},
						"common_name": schema.StringAttribute{
							Optional:    true,
							Description: " Common name of the certificate. Can be \"*\" to allow all values, or a string with '*' as a glob character",
						},
						"dns_san": schema.StringAttribute{
							Optional:    true,
							Description: " Comma delimited list of acceptable domain names for the Subject Alternative Name extension. Names use '*' as a glob character Default no values allowed	",
						},
						"uri_san": schema.StringAttribute{
							Optional:    true,
							Description: " Comma delimited list of acceptable URIs for the Subject Alternative Name extension. Names use '*' as a glob character. Default no values allowed",
						},
						"ip_san": schema.StringAttribute{
							Optional:    true,
							Description: " omma delimited list of acceptable IPs for the Subject Alternative Name extension. Can be IP addresses or CIDRs. Default no values allowed",
						},
						"key_type": schema.StringAttribute{
							Optional:    true,
							Description: "Comma delimited list of acceptable key types for the signed certificate. Can be '*' to allow all key types.Default rsa:2048,rsa:4096,ec:secp256r1",
						},
						"hash_type": schema.StringAttribute{
							Optional:    true,
							Description: "Comma delimited list of acceptable hash_types for the signed certificate. Can be '*' to allow all key types. This constraint is ignored for key_types that don't use hashing as part of the signature (ed25519)",
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *templateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan templateResourceModel
	tflog.SetField(ctx, "Template", plan)
	tflog.Info(ctx, "Create Template Request fist")
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	elements := plan.Options // Get the list of elements from the plan
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Create Template Request")

	//This is totally wrong, but I don't know how to do it better
	constraintTypes := []kmi.ConstraintType{}
	for _, v := range elements {
		if !v.CommonName.IsNull() {
			constraintTypes = append(constraintTypes, kmi.ConstraintType{
				Type: "common_name",
				Text: v.CommonName.ValueString(),
			})
		}

	}
	kmitemplate := kmi.Template{
		Constraint: constraintTypes,
	}
	tflog.SetField(ctx, "Template", kmitemplate)
	tflog.Debug(ctx, "CreateTemplateOrSign Template")
	err := r.client.CreateTemplateOrSign(plan.CACollectionName.ValueString(), plan.CADefinitionName.ValueString(), plan.TemplateName.ValueString(), kmitemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Template",
			"Could not create template, unexpected error: "+err.Error(),
		)
		return
	}

	kmiSigner := kmi.Template{
		Collectionacl: &kmi.ApproveClientCollection{
			Target: plan.ClientCollectionName.ValueString(),
		},
	}
	tflog.Debug(ctx, "CreateTemplateOrSign CSR signer")
	err = r.client.CreateTemplateOrSign(plan.CACollectionName.ValueString(), plan.CADefinitionName.ValueString(), plan.TemplateName.ValueString(), kmiSigner)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error signing the request",
			"Could not signing the request, unexpected error: "+err.Error(),
		)
		return
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *templateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state templateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateDetails, err := r.client.GetTemplate(state.CACollectionName.ValueString(), state.CADefinitionName.ValueString(), state.TemplateName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading HashiCups Order",
			"Could not read Template "+state.TemplateName.ValueString()+": "+err.Error(),
		)
		return
	}

	for _, v := range templateDetails.Constraint {
		// state.Options[v.Text] = types.String(v.Text)
		if (v.Type == "common_name") && (v.Text != "*") {
			state.Options = append(state.Options, templateResourceModelOptions{
				CommonName: types.StringValue(v.Text),
			})
		}
	}
	optionsfromKmi := []templateResourceModelOptions{}
	state.Options = optionsfromKmi
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *templateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *templateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *templateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type templateResourceModel struct {
	CACollectionName     types.String `tfsdk:"ca_collection"`
	CADefinitionName     types.String `tfsdk:"ca_definition"`
	TemplateName         types.String `tfsdk:"template_name"`
	ClientCollectionName types.String `tfsdk:"client_collection"`

	Options     []templateResourceModelOptions `tfsdk:"options"`
	LastUpdated types.String                   `tfsdk:"last_updated"`
}
type templateResourceModelOptions struct {
	Min_ttl             types.String `tfsdk:"min_ttl"`
	Max_ttl             types.String `tfsdk:"max_ttl"`
	Leaf_exceeds_ca_ttl types.String `tfsdk:"leaf_exceeds_ca_ttl"`
	Allow_ca            types.String `tfsdk:"allow_ca"`
	CommonName          types.String `tfsdk:"common_name"`
	Dns_san             types.String `tfsdk:"dns_san"`
	Uri_san             types.String `tfsdk:"uri_san"`
	Ip_san              types.String `tfsdk:"ip_san"`
	Key_type            types.String `tfsdk:"key_type"`
	Hash_type           types.String `tfsdk:"hash_type"`
}
