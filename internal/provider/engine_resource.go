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
	_ resource.Resource              = &engineResource{}
	_ resource.ResourceWithConfigure = &engineResource{}
)

// NewEngineResource is a helper function to simplify the provider implementation.
func NewEngineResource() resource.Resource {

	return &engineResource{}
}

// engineResource is the resource implementation.
type engineResource struct {
	client *kmi.KMIRestClient
}

// Metadata returns the resource type name.
func (r *engineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engine"
}

// Schema defines the schema for the resource.
func (r *engineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"engine": schema.StringAttribute{
				Required: true,
			},
			"account_name": schema.StringAttribute{
				Required: true,
			},
			"api_endpoint": schema.StringAttribute{
				Required: true,
			},
			"cas_base64": schema.StringAttribute{
				Required: true,
			},
			"source": schema.StringAttribute{
				Optional: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"workloads": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"serviceaccount": schema.StringAttribute{
							Required: true,
						},
						"namespace": schema.StringAttribute{
							Required: true,
						},
						"region": schema.StringAttribute{
							Required: true,
						},
						"kubernetes_service_account": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *engineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EngineResourceModel
	diags := req.Plan.Get(ctx, &plan)
	tflog.SetField(ctx, "Plan Engine", plan)
	tflog.Debug(ctx, "Getting Identity engine")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	options := []kmi.KMIOption{{
		Text: plan.CertificateDataAuthority.ValueString(),
		Name: "cas_base64",
	}, {
		Text: plan.ApiEndpoint.ValueString(),
		Name: "endpoint_url",
	}}
	workloads := []kmi.KMIWorkload{}

	for _, projection := range plan.Workloads {
		kubernetes_service_account := fmt.Sprintf("system:serviceaccount:%s:%s", projection.Namespace.ValueString(), projection.ServiceAccount.ValueString())
		wkload := kmi.KMIWorkload{
			Projection:               projection.Name.ValueString(),
			KubernetesServiceAccount: kubernetes_service_account,
			Region:                   projection.Region.ValueString(),
		}
		workloads = append(workloads, wkload)
	}
	engine := kmi.KMIEngine{
		Cloud:     "linode",
		Type:      "kubernetes",
		Option:    options,
		Workloads: workloads,
	}
	err := r.client.SaveIdentityEngine(plan.AccountName.ValueString(), plan.Engine.ValueString(), engine)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Identity Engine",
			"Could not create Identity, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "After Saving Identity engine")

	identityEngine, err := r.client.GetIdentityEngine(plan.AccountName.ValueString(), plan.Engine.ValueString())
	tflog.SetField(ctx, "Identity Engine", identityEngine)
	tflog.Debug(ctx, "Getting Identity engine")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Identity Engine",
			"Could not create Identity, unexpected error: "+err.Error(),
		)
		return
	}

	for wrkIndex, projectionafter := range identityEngine.Workload {
		kmiprojection, err := r.client.GetWorkloadDetails(plan.AccountName.ValueString(), plan.Engine.ValueString(), projectionafter.Projection)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Identity Engine",
				"Could not create Identity, unexpected error: "+err.Error(),
			)
			return
		}

		kmiserviceAcc := kmiprojection.KubernetesServiceAccount.Text
		k8ServiceAccount := strings.Split(kmiserviceAcc, ":")[3]
		k8Namepace := strings.Split(kmiserviceAcc, ":")[2]
		wrkmodel := WorkloadResourceModel{
			Name:                     types.StringValue(kmiprojection.Projection),
			KubernetesServiceAccount: types.StringValue(kmiserviceAcc),
			ServiceAccount:           types.StringValue(k8ServiceAccount),
			Namespace:                types.StringValue(k8Namepace),
			Region:                   types.StringValue(kmiprojection.Region.Text),
		}
		plan.Workloads[wrkIndex] = wrkmodel
	}
	// plan.LastUpdated = identityEngine.Published
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *engineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EngineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityEngine, err := r.client.GetIdentityEngine(state.AccountName.ValueString(), state.Engine.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Identity Engine",
			"Could not create Identity, unexpected error: "+err.Error(),
		)
		return
	}

	if len(identityEngine.Workload) != 0 {
		state.Workloads = []WorkloadResourceModel{}
	}

	for _, projectionafter := range identityEngine.Workload {
		kmiprojection, err := r.client.GetWorkloadDetails(state.AccountName.ValueString(), state.Engine.ValueString(), projectionafter.Projection)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Identity Engine",
				"Could not create Identity, unexpected error: "+err.Error(),
			)
			return
		}

		kmiserviceAcc := kmiprojection.KubernetesServiceAccount.Text
		k8ServiceAccount := strings.Split(kmiserviceAcc, ":")[3]
		k8Namepace := strings.Split(kmiserviceAcc, ":")[2]
		wrkmodel := WorkloadResourceModel{
			Name:                     types.StringValue(kmiprojection.Projection),
			KubernetesServiceAccount: types.StringValue(kmiserviceAcc),
			ServiceAccount:           types.StringValue(k8ServiceAccount),
			Namespace:                types.StringValue(k8Namepace),
			Region:                   types.StringValue(kmiprojection.Region.Text),
		}
		state.Workloads = append(state.Workloads, wrkmodel)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *engineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan EngineResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	options := []kmi.KMIOption{{
		Text: plan.CertificateDataAuthority.ValueString(),
		Name: "cas_base64",
	}, {
		Text: plan.ApiEndpoint.ValueString(),
		Name: "endpoint_url",
	}}
	workloads := []kmi.KMIWorkload{}

	for _, projection := range plan.Workloads {
		kubernetes_service_account := fmt.Sprintf("system:serviceaccount:%s:%s", projection.Namespace.ValueString(), projection.ServiceAccount.ValueString())
		wkload := kmi.KMIWorkload{
			Projection:               projection.Name.ValueString(),
			KubernetesServiceAccount: kubernetes_service_account,
			Region:                   projection.Region.ValueString(),
		}
		workloads = append(workloads, wkload)
	}
	engine := kmi.KMIEngine{
		Cloud:     "linode",
		Type:      "kubernetes",
		Option:    options,
		Workloads: workloads,
	}
	err := r.client.SaveIdentityEngine(plan.AccountName.ValueString(), plan.Engine.ValueString(), engine)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Identity Engine",
			"Could not create Identity, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "After Saving Identity engine")

	identityEngine, err := r.client.GetIdentityEngine(plan.AccountName.ValueString(), plan.Engine.ValueString())
	tflog.SetField(ctx, "Identity Engine", identityEngine)
	tflog.Debug(ctx, "Getting Identity engine")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Identity Engine",
			"Could not create Identity, unexpected error: "+err.Error(),
		)
		return
	}

	for wrkIndex, projectionafter := range identityEngine.Workload {
		kmiprojection, err := r.client.GetWorkloadDetails(plan.AccountName.ValueString(), plan.Engine.ValueString(), projectionafter.Projection)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Identity Engine",
				"Could not create Identity, unexpected error: "+err.Error(),
			)
			return
		}

		kmiserviceAcc := kmiprojection.KubernetesServiceAccount.Text
		k8ServiceAccount := strings.Split(kmiserviceAcc, ":")[3]
		k8Namepace := strings.Split(kmiserviceAcc, ":")[2]
		wrkmodel := WorkloadResourceModel{
			Name:                     types.StringValue(kmiprojection.Projection),
			KubernetesServiceAccount: types.StringValue(kmiserviceAcc),
			ServiceAccount:           types.StringValue(k8ServiceAccount),
			Namespace:                types.StringValue(k8Namepace),
			Region:                   types.StringValue(kmiprojection.Region.Text),
		}
		plan.Workloads[wrkIndex] = wrkmodel
	}
	// plan.LastUpdated = identityEngine.Published
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *engineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EngineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityEngine, err := r.client.GetIdentityEngine(state.AccountName.ValueString(), state.Engine.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Identity Engine",
			"Could not create Identity, unexpected error: "+err.Error(),
		)
		return
	}

	for _, projectionafter := range identityEngine.Workload {
		err := r.client.DeleteWorkload(state.AccountName.ValueString(), state.Engine.ValueString(), projectionafter.Projection)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Identity Engine",
				"Could not create Identity, unexpected error: "+err.Error(),
			)
			return
		}

	}

	err = r.client.DeleteIdentityEngine(state.AccountName.ValueString(), state.Engine.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Identity Engine",
			"Could not create Identity, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *engineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type EngineResourceModel struct {
	Engine                   types.String            `tfsdk:"engine"`
	AccountName              types.String            `tfsdk:"account_name"`
	ApiEndpoint              types.String            `tfsdk:"api_endpoint"`
	CertificateDataAuthority types.String            `tfsdk:"cas_base64"`
	Source                   types.String            `tfsdk:"source"`
	Workloads                []WorkloadResourceModel `tfsdk:"workloads"`
	LastUpdated              types.String            `tfsdk:"last_updated"`
}

// Name            types.String `tfsdk:"name"`
type WorkloadResourceModel struct {
	Name                     types.String `tfsdk:"name"`
	KubernetesServiceAccount types.String `tfsdk:"kubernetes_service_account"`
	ServiceAccount           types.String `tfsdk:"serviceaccount"`
	Namespace                types.String `tfsdk:"namespace"`
	Region                   types.String `tfsdk:"region"`
}
