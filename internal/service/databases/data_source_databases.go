// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package databases

import (
	"context"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceconfig "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/config"
	mdutils "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/markdown"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DatabasesDataSource{}

func NewDatabasesDataSource() datasource.DataSource {
	return &DatabasesDataSource{}
}

// DatabasesDataSource defines the data source implementation.
type DatabasesDataSource struct {
	config *resourceconfig.ResourceConfig
}

// DocumentDataSourceModel describes the data source data model.
type DatabasesDataSourceModel struct {
	Name      types.String `tfsdk:"database"`
	Databases types.List   `tfsdk:"databases"`
}

var DatabaseElementType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	},
}

func (d *DatabasesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_databases"
}

func (d *DatabasesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: mdutils.FormatResourceDescription(`
			This resource reads a list of documents 
			in a database with a given filter.
		`),

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Name of databases to read. Regex is supported.",
			},
			"databases": schema.ListAttribute{
				ElementType:         DatabaseElementType,
				MarkdownDescription: "Array of database resource.",
				Computed:            true,
			},
		},
	}
}

func (d *DatabasesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	config, diags := resourceconfig.FromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.config = config
}

func (d *DatabasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	client := mongoclient.New(ctx, d.config.ClientConfig).WithLogger(d.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.Append(
				errs.NewMongoClientError(err).ToDiagnostic(),
			)
			return
		}

		var data DatabasesDataSourceModel

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
