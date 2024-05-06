// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package document

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resourceutils"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/collection"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DocumentResource{}
var _ resource.ResourceWithImportState = &DocumentResource{}

func NewDocumentResource() resource.Resource {
	return &DocumentResource{}
}

// DocumentResource defines the resource implementation.
type DocumentResource struct {
	client *mongoclient.MongoClient
}

// DocumentResourceModel describes the resource data model.
type DocumentResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Database   types.String `tfsdk:"database"`
	Collection types.String `tfsdk:"collection"`
	DocumentId types.String `tfsdk:"document_id"`
	Document   types.String `tfsdk:"document"`
}

func (r *DocumentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_document"
}

func (r *DocumentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource creates a single document in a collection in a database on the MongoDB server.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier. Has a value with a format of databases/<database_name>/collections/<collection_name>/documents/<document_id>.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"collection": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the collection to create the document in.",
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
			"document_id": schema.StringAttribute{
				MarkdownDescription: "Document ID of the document.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"document": schema.StringAttribute{
				MarkdownDescription: "Document to insert into the collection.",
				Required:            true,
			},
		},
	}
}

func (r *DocumentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DocumentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DocumentResourceModel

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

	// Create the document
	var document interface{}
	if err := json.Unmarshal([]byte(data.Document.ValueString()), &document); err != nil {
		resp.Diagnostics.AddError("Invalid JSON Input", err.Error())
		return
	}
	oid, err := collection.InsertOne(document)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	data.DocumentId = basetypes.NewStringValue(oid)

	// Perform a read operation to get the document
	resp.Diagnostics.Append(resourceRead(r.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DocumentResourceModel

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

func (r *DocumentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DocumentResourceModel

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

	// Update the document
	var document interface{}
	if err := json.Unmarshal([]byte(data.Document.ValueString()), &document); err != nil {
		resp.Diagnostics.AddError("Invalid JSON Input", err.Error())
		return
	}
	if err := collection.UpdateByID(data.DocumentId.ValueString(), document); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	// Perform a read operation to get the document
	resp.Diagnostics.Append(resourceRead(r.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DocumentResourceModel

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

	// Delete the document
	if err := collection.DeleteByID(data.DocumentId.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}

func (r *DocumentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
	if id.Document() == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Document ID is required")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), id.Database())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection"), id.Collection())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("document_id"), id.Document())...)
}
