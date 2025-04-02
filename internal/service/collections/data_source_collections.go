// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package collections

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
var _ datasource.DataSource = &CollectionsDataSource{}

func NewCollectionsDataSource() datasource.DataSource {
	return &CollectionsDataSource{}
}

// CollectionsDataSource defines the data source implementation.
type CollectionsDataSource struct {
	config *resourceconfig.ResourceConfig
}

// DocumentDataSourceModel describes the data source data model.
type CollectionsDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Database    types.String `tfsdk:"database"`
	Collections types.List   `tfsdk:"collections"`
}

var CollectionElementType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":       types.StringType,
		"database": types.StringType,
		"name":     types.StringType,
	},
}

func (d *CollectionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_collections"
}

func (d *CollectionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: mdutils.FormatResourceDescription(`
			This resource reads a list of collections
			in a specified database on a cluster.
		`),

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Name of collections to read. Regex is supported.",
			},
			"database": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the database to read collections from.",
			},
			"collections": schema.ListAttribute{
				ElementType:         CollectionElementType,
				MarkdownDescription: "List of collection resources. Refer to collection data source documentation for more details.",
				Computed:            true,
			},
		},
	}
}

func (d *CollectionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	config, diags := resourceconfig.FromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.config = config
}

func (d *CollectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	client := mongoclient.New(ctx, d.config.ClientConfig).WithLogger(d.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.Append(
				errs.NewMongoClientError(err).ToDiagnostic(),
			)
			return
		}

		var data CollectionsDataSourceModel

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
