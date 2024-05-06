// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package collection

import (
	"context"
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resourceutils"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CollectionResource{}
var _ resource.ResourceWithImportState = &CollectionResource{}

func NewCollectionResource() resource.Resource {
	return &CollectionResource{}
}

// CollectionResource defines the resource implementation.
type CollectionResource struct {
	client *mongoclient.MongoClient
}

// CollectionResourceModel describes the resource data model.
type CollectionResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Id           types.String `tfsdk:"id"`
	Database     types.String `tfsdk:"database"`
	ForceDestroy types.Bool   `tfsdk:"force_destroy"`
}

func (r *CollectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_collection"
}

func (r *CollectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource creates a single collection in a database on the MongoDB server.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the collection",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"database": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the database to create the collection in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier. Has a value with a format of databases/<database_name>/collections/<collection_name>.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"force_destroy": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "By default, the provider will not destroy the collection if it contains any data. Set this to true to force destroy the collection even if it contains data.",
			},
		},
	}
}

func (r *CollectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*mongoclient.MongoClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *mongoclient.MongoClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *CollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CollectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Connect to the MongoDB server
	if err := r.client.WithContext(ctx).Connect(); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	defer r.client.Disconnect()

	// Check if the database exists
	database := database.CheckExistance(r.client, data.Database.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the collection
	collection := database.Collection(data.Name.ValueString())
	if err := collection.EnsureExistance(); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	// Perform the read operation
	resp.Diagnostics.Append(resourceRead(r.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CollectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Connect to the MongoDB server
	if err := r.client.WithContext(ctx).Connect(); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	defer r.client.Disconnect()

	// Perform read operation
	resp.Diagnostics.Append(resourceRead(r.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CollectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CollectionResourceModel
	if err := r.client.WithContext(ctx).Connect(); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	defer r.client.Disconnect()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the database exists
	database := r.client.Database(data.Database.ValueString())
	exists, err := database.Exists()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if !exists {
		// We don't need to check if the collection exists,
		// as the database doesn't exist
		return
	}

	// Check if the collection exists
	collection := database.Collection(data.Name.ValueString())
	exists, err = collection.Exists()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if !exists {
		// Collection doesn't exist, nothing to delete
		return
	}

	// Check if the collection is empty
	isEmpty, err := collection.IsEmpty()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if !isEmpty && !data.ForceDestroy.ValueBool() {
		resp.Diagnostics.AddError("Collection Not Empty", "Collection contains data. Set force_destroy to true to delete the collection.")
		return
	}

	// Delete the collection
	if err := collection.Drop(); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}

func (r *CollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := resourceutils.NewId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())
		return
	}
	if id.Database() == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Database name is required")
		return
	}
	if id.Collection() == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Collection name is required")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), id.Database())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), id.Collection())...)
}
