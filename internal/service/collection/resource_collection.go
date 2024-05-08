// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package collection

import (
	"context"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceconfig "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/config"
	resourceid "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/id"
	mdutils "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/markdown"

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
	config *resourceconfig.ResourceConfig
}

// CollectionResourceModel describes the resource data model.
type CollectionResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Database     types.String `tfsdk:"database"`
	Name         types.String `tfsdk:"name"`
	ForceDestroy types.Bool   `tfsdk:"force_destroy"`
}

func (r *CollectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_collection"
}

func (r *CollectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: mdutils.FormatResourceDescription(`
			This resource creates a collection in a database on the MongoDB server.
		`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: mdutils.FormatSchemaDescription(
					`
						Resource identifier.
						
						ID has a value with a format of the following:

						%s

						Note that this format is used for importing the resource into Terraform state.
						Import the resource using the following command:

						%s
					`,
					mdutils.CodeBlock("", "databases/<database>/collections/<name>"),
					mdutils.CodeBlock("bash", "terraform import mongodb_database_collection.<resource_name> databases/<database>/collections/<name>"),
				),
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
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the collection",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"force_destroy": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: mdutils.FormatSchemaDescription(`
					Whether to force destroy the collection.
					
					By default, the provider will not destroy the collection if it contains any data. 
					The provider decides whether the collection contains data based on the collection's document count. If the collection contains any documents, the provider will not destroy the collection.

					Set this to true to force destroy the collection even if it contains data.
				`),
			},
		},
	}
}

func (r *CollectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config, diags := resourceconfig.FromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.config = config
}

func (r *CollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client := r.config.Client.WithContext(ctx).WithLogger(r.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.AddError(errornames.MongoClientError, err.Error())
			return
		}

		var data CollectionResourceModel

		// Read Terraform plan data into the model
		resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform create operation
		resp.Diagnostics.Append(resourceCreate(client, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	})
}

func (r *CollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client := r.config.Client.WithContext(ctx).WithLogger(r.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.AddError(errornames.MongoClientError, err.Error())
			return
		}

		var data CollectionResourceModel

		// Read Terraform prior state data into the model
		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform read operation
		resp.Diagnostics.Append(resourceRead(client, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save updated data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	})
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
	client := r.config.Client.WithContext(ctx).WithLogger(r.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.AddError(errornames.MongoClientError, err.Error())
			return
		}

		var data CollectionResourceModel
		// Read Terraform prior state data into the model
		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform delete operation
		resp.Diagnostics.Append(resourceDelete(client, &data)...)
	})
}

func (r *CollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := resourceid.New(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(errornames.InvalidImportID, err.Error())
		return
	}
	if id.Database() == "" {
		resp.Diagnostics.AddError(errornames.InvalidImportID, "Database name is required")
		return
	}
	if id.Collection() == "" {
		resp.Diagnostics.AddError(errornames.InvalidImportID, "Collection name is required")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), id.Database())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), id.Collection())...)
}
