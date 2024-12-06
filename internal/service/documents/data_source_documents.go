// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package documents

import (
	"context"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceconfig "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/config"
	mdutils "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/markdown"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DocumentsDataSource{}

func NewDocumentsDataSource() datasource.DataSource {
	return &DocumentsDataSource{}
}

// DocumentsDataSource defines the data source implementation.
type DocumentsDataSource struct {
	config *resourceconfig.ResourceConfig
}

// DocumentDataSourceModel describes the data source data model.
type DocumentsDataSourceModel struct {
	Database   types.String `tfsdk:"database"`
	Collection types.String `tfsdk:"collection"`
	Filter     types.String `tfsdk:"filter"`
	Documents  types.String `tfsdk:"documents"`
}

func (d *DocumentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_documents"
}

func (d *DocumentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: mdutils.FormatResourceDescription(`
			This resource reads a single document in a collection 
			in a database on the MongoDB server.
		`),

		Attributes: map[string]schema.Attribute{
			"database": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the database to read the collection in.",
			},
			"collection": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the collection to read the document in.",
			},
			"filter": schema.StringAttribute{
				MarkdownDescription: mdutils.FormatSchemaDescription(
					`
						Filter to find the document in the collection.

						The value of this attribute is a stringified JSON.
						Note that you should escape every double quote in the JSON string.

						In terraform, you can achieve this by simply using the 
						%s function:

						%s
					`,
					mdutils.InlineCodeBlock("jsonencode"),
					mdutils.CodeBlock("terraform", "document = jsonencode({ key = \"value\" })"),
				),
				Required: true,
			},
			"documents": schema.StringAttribute{
				MarkdownDescription: mdutils.FormatSchemaDescription(
					`
						Documents read from the collection.

						The value of this attribute is a stringified JSON, 
						with every double quote escaped with a backslash.
						This means that the JSON string contains backslashes before every double quote.

						In terraform, you'll be able to smoothly decode the JSON string by using the %s function.

						%s
					`,
					mdutils.InlineCodeBlock("jsondecode"),
					mdutils.CodeBlock("terraform", "decoded = jsondecode(document)"),
				),
				Computed: true,
			},
		},
	}
}

func (d *DocumentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	config, diags := resourceconfig.FromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.config = config
}

func (d *DocumentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	client := mongoclient.New(ctx, d.config.ClientConfig).WithLogger(d.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.Append(
				errs.NewMongoClientError(err).ToDiagnostic(),
			)
			return
		}

		var data DocumentsDataSourceModel

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
