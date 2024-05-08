// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package document

import (
	"context"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceconfig "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/config"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DocumentDataSource{}

func NewDocumentDataSource() datasource.DataSource {
	return &DocumentDataSource{}
}

// DocumentDataSource defines the data source implementation.
type DocumentDataSource struct {
	config *resourceconfig.ResourceConfig
}

// DocumentDataSourceModel describes the data source data model.
type DocumentDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Database   types.String `tfsdk:"database"`
	Collection types.String `tfsdk:"collection"`
	DocumentId types.String `tfsdk:"document_id"`
	Document   types.String `tfsdk:"document"`
}

func (d *DocumentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_document"
}

func (d *DocumentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source reads a document in a collection in a database on the MongoDB server.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier. Has a value with a format of databases/<database_name>/collections/<collection_name>/documents/<document_id>.",
			},
			"database": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the database to read the collection in.",
			},
			"collection": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the collection to read the document in.",
			},
			"document_id": schema.StringAttribute{
				MarkdownDescription: "Document ID of the document.",
				Required:            true,
			},
			"document": schema.StringAttribute{
				MarkdownDescription: "Document to insert into the collection.",
				Computed:            true,
			},
		},
	}
}

func (d *DocumentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	config, diags := resourceconfig.FromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.config = config
}

func (d *DocumentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	client := d.config.Client.WithContext(ctx).WithLogger(d.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.AddError(errornames.MongoClientError, err.Error())
			return
		}

		var data DocumentDataSourceModel

		// Read Terraform prior state data into the model
		resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform the read operation
		resp.Diagnostics.Append(dataSourceRead(client, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save updated data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	})
}
