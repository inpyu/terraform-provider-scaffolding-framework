package provider

import (
	"context"
	"fmt"

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

type cafeDataSourceModel struct {
	Cafe []cafeResourceModel `tfsdk:"cafe"`
}

// orderResource is the resource implementation.
type cafeResource struct {
	client *hashicups.Client
}

type cafeResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Address     types.String `tfsdk:"address"`
	Description types.String `tfsdk:"description"`
	Image       types.String `tfsdk:"image"`
}

// Metadata returns the resource type name.
func (r *cafeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cafe"
}

// Schema defines the schema for the resource.
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
				Computed: true,
			},
			"address": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"image": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create a new resource.
func (r *cafeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cafeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var items []hashicups.Cafe

	// Create new cafe
	cafe, err := r.client.CreateCafe(items)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cafe",
			"Could not create cafe, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.Int64Value(int64(cafe.ID))
	plan.Name = types.StringValue(cafe.Name)
	plan.Address = types.StringValue(cafe.Address)
	plan.Description = types.StringValue(cafe.Description)
	plan.Image = types.StringValue(cafe.Image)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (d *cafeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cafeDataSourceModel

	cafes, err := d.client.GetCafes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Cafes",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, cafe := range cafes {
		cafeState := cafeResourceModel{
			ID:          types.Int64Value(int64(cafe.ID)),
			Name:        types.StringValue(cafe.Name),
			Address:     types.StringValue(cafe.Address),
			Description: types.StringValue(cafe.Description),
			Image:       types.StringValue(cafe.Image),
		}
		state.Cafe = append(state.Cafe, cafeState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *cafeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cafeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var hashicupsItems []hashicups.Cafe

	// Update existing order
	_, err := r.client.UpdateCafe(plan.ID.String(), hashicupsItems)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating HashiCups Cafe",
			"Could not update cafe, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	cafe, err := r.client.GetCafe(plan.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading HashiCups Cafe",
			"Could not read HashiCups ocaferder ID "+plan.ID.String()+": "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.Name = types.StringValue(cafe.Name)
	plan.Address = types.StringValue(cafe.Address)
	plan.Description = types.StringValue(cafe.Description)
	plan.Image = types.StringValue(cafe.Image)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *cafeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cafeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing cafe
	err := r.client.DeleteCafe(state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Cafe",
			"Could not delete cafe, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
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
