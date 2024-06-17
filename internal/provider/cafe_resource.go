package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/inpyu/hashicups-client-go"
)

var (
	_ resource.Resource              = &cafeResource{}
	_ resource.ResourceWithConfigure = &cafeResource{}
)

func NewCafeResource() resource.Resource {
	return &cafeResource{}
}

type cafeResource struct {
	client *hashicups.Client
}

type cafeResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Address     types.String `tfsdk:"address"`
	Description types.String `tfsdk:"description"`
	Image       types.String `tfsdk:"image"`
}

func (r *cafeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cafe"
}

func (r *cafeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"address": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"image": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *cafeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cafeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cafe := hashicups.Cafe{
		Name:        plan.Name.ValueString(),
		Address:     plan.Address.ValueString(),
		Description: plan.Description.ValueString(),
		Image:       plan.Image.ValueString(),
	}

	createdCafe, err := r.client.CreateCafe([]hashicups.Cafe{cafe})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cafe",
			"Could not create cafe, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(createdCafe.ID))
	plan.Name = types.StringValue(createdCafe.Name)
	plan.Address = types.StringValue(createdCafe.Address)
	plan.Description = types.StringValue(createdCafe.Description)
	plan.Image = types.StringValue(createdCafe.Image)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *cafeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cafeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cafeID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Cafe ID",
			"Could not convert cafe ID to integer: "+err.Error(),
		)
		return
	}

	// Assume GetCafe now returns a list of cafes
	cafes, err := r.client.GetCafe(strconv.Itoa(cafeID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Cafe",
			err.Error(),
		)
		return
	}

	if len(cafes) == 0 {
		resp.Diagnostics.AddError(
			"Cafe Not Found",
			"No cafe found with the given ID",
		)
		return
	}

	cafe := cafes[0]

	// Map response body to model
	state.ID = types.StringValue(strconv.Itoa(cafe.ID))
	state.Address = types.StringValue(cafe.Address)
	state.Image = types.StringValue(cafe.Image)
	state.Name = types.StringValue(cafe.Name)
	state.Description = types.StringValue(cafe.Description)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *cafeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cafeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the ID from string to int
	cafeID, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Cafe ID",
			"Could not convert cafe ID to integer: "+err.Error(),
		)
		return
	}

	// Create a cafe object
	cafe := hashicups.Cafe{
		ID:          cafeID, // ID is an int
		Name:        plan.Name.ValueString(),
		Address:     plan.Address.ValueString(),
		Description: plan.Description.ValueString(),
		Image:       plan.Image.ValueString(),
	}

	// Update the existing cafe
	updatedCafe, err := r.client.UpdateCafe(plan.ID.ValueString(), []hashicups.Cafe{cafe})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating HashiCups Cafe",
			"Could not update cafe, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated items
	plan.ID = types.StringValue(strconv.Itoa(updatedCafe.ID))
	plan.Name = types.StringValue(updatedCafe.Name)
	plan.Address = types.StringValue(updatedCafe.Address)
	plan.Description = types.StringValue(updatedCafe.Description)
	plan.Image = types.StringValue(updatedCafe.Image)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *cafeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cafeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cafeID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Cafe",
			"Could not convert cafe ID to integer: "+err.Error(),
		)
		return
	}

	err = r.client.DeleteCafe(strconv.Itoa(cafeID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Cafe",
			"Could not delete cafe, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *cafeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hashicups.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
