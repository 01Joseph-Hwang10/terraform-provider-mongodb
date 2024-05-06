// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index

import (
	"context"
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexDataSource{}

func NewIndexDataSource() datasource.DataSource {
	return &IndexDataSource{}
}

// IndexDataSource defines the data source implementation.
type IndexDataSource struct {
	client *mongoclient.MongoClient
}

// IndexDataSourceModel describes the data source data model.
type IndexDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Database   types.String `tfsdk:"database"`
	Collection types.String `tfsdk:"collection"`
	Field      types.String `tfsdk:"key"`
	Direction  types.Int64  `tfsdk:"direction"`
	IndexName  types.String `tfsdk:"index_name"`
	Unique     types.Bool   `tfsdk:"unique"`
}

func (d *IndexDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_index"
}

func (d *IndexDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source reads an index for single field in a collection in a database on the MongoDB server.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier. Has a value with a format of databases/<database_name>/collections/<collection_name>/indexes/<index_name>.",
			},
			"database": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the database to read the collection in.",
			},
			"collection": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the collection to read the index in.",
			},
			"field": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the field to create the index on.",
			},
			"direction": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Direction of the index. 1 for ascending, -1 for descending.",
			},
			"unique": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, this index has a unique constraint.",
			},
			"index_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the index.",
			},
		},
	}
}

func (d *IndexDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *IndexDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IndexDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Connect to the MongoDB server
	if err := d.client.WithContext(ctx).Connect(); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	defer d.client.Disconnect()

	// Perform read operation
	resp.Diagnostics.Append(dataSourceRead(d.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
