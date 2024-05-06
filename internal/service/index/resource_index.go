// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index

import (
	"context"
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resourceutils"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/collection"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IndexResource{}
var _ resource.ResourceWithImportState = &IndexResource{}

func NewIndexResource() resource.Resource {
	return &IndexResource{}
}

// IndexResource defines the resource implementation.
type IndexResource struct {
	client *mongoclient.MongoClient
}

// IndexResourceModel describes the resource data model.
type IndexResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Database     types.String `tfsdk:"database"`
	Collection   types.String `tfsdk:"collection"`
	Field        types.String `tfsdk:"key"`
	Direction    types.Int64  `tfsdk:"direction"`
	IndexName    types.String `tfsdk:"index_name"`
	Unique       types.Bool   `tfsdk:"unique"`
	ForceDestroy types.Bool   `tfsdk:"force_destroy"`
}

func (r *IndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_index"
}

func (r *IndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource creates an index for single field in a collection in a database on the MongoDB server.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier. Has a value with a format of databases/<database_name>/collections/<collection_name>/indexes/<index_name>.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"database": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the database to create the collection in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"collection": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the collection to create the index in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"field": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the field to create the index on.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"direction": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Direction of the index. 1 for ascending, -1 for descending.",
				Default:             int64default.StaticInt64(1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"unique": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "If true, creates a unique index.",
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"index_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the index.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"force_destroy": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "By default, the provider will fail on destroying the index. Set this to true to force destroy the index.",
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *IndexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IndexResourceModel

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

	// Check if the collection exists
	collection := collection.CheckExistance(database, data.Collection.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the index
	index := collection.IndexFromField(data.Field.ValueString(), int(data.Direction.ValueInt64()), data.Unique.ValueBool())
	if err := index.EnsureExistance(); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	// Set index name
	data.IndexName = basetypes.NewStringValue(index.Name())

	// Perform read operation
	resp.Diagnostics.Append(resourceRead(r.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IndexResourceModel
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

func (r *IndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IndexResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IndexResourceModel

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
	collection := database.Collection(data.Collection.ValueString())
	exists, err = collection.Exists()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if !exists {
		// Collection doesn't exist, nothing to delete
		return
	}

	// If force destroy is not set, fail the deletion
	if !data.ForceDestroy.ValueBool() {
		resp.Diagnostics.AddError("Deletion Forbidden", "Index deletion is not allowed by default. Set force_destroy to true to delete the index.")
		return
	}

	// Delete the index
	index := collection.Index(data.IndexName.String())
	if err := index.Drop(); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}

func (r *IndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
	if id.Index() == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Index ID is required")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), id.Database())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection"), id.Collection())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("index_name"), id.Index())...)
}
