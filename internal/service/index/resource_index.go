// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index

import (
	"context"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceconfig "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/config"
	resourceid "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/id"

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
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IndexResource{}
var _ resource.ResourceWithImportState = &IndexResource{}

func NewIndexResource() resource.Resource {
	return &IndexResource{}
}

// IndexResource defines the resource implementation.
type IndexResource struct {
	config *resourceconfig.ResourceConfig
}

// IndexResourceModel describes the resource data model.
type IndexResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Database     types.String `tfsdk:"database"`
	Collection   types.String `tfsdk:"collection"`
	IndexName    types.String `tfsdk:"index_name"`
	Field        types.String `tfsdk:"field"`
	Direction    types.Int64  `tfsdk:"direction"`
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
			"index_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the index.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
	config, diags := resourceconfig.FromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.config = config
}

func (r *IndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client := r.config.Client.WithContext(ctx)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		var data IndexResourceModel

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

func (r *IndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client := r.config.Client.WithContext(ctx)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		var data IndexResourceModel
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
	client := r.config.Client.WithContext(ctx)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		var data IndexResourceModel

		// Read Terraform prior state data into the model
		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform delete operation
		resp.Diagnostics.Append(resourceDelete(client, &data)...)
	})
}

func (r *IndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
	if id.Index() == "" {
		resp.Diagnostics.AddError(errornames.InvalidImportID, "Index ID is required")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), id.Database())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection"), id.Collection())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("index_name"), id.Index())...)
}
