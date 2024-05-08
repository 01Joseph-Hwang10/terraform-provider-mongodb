// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"context"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceconfig "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/config"
	mdutils "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/markdown"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DatabaseDataSource{}

func NewDatabaseDataSource() datasource.DataSource {
	return &DatabaseDataSource{}
}

// DatabaseDataSource defines the data source implementation.
type DatabaseDataSource struct {
	config *resourceconfig.ResourceConfig
}

// DatabaseDataSourceModel describes the data source data model.
type DatabaseDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *DatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *DatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: mdutils.FormatResourceDescription(`
			This resource reads a visible database on the MongoDB server.

			The meaning of **visible** is that you can get the database information
			with [listDatabases](https://www.mongodb.com/docs/manual/reference/command/listDatabases/)
			command or [show dbs](https://www.mongodb.com/docs/mongodb-shell/reference/access-mdb-shell-help/#show-available-databases) command.

			By default, mongodb automatically creates a database 
			when you first store data in that database. 
			(See [this MongoDB documentation](https://www.mongodb.com/docs/manual/core/databases-and-collections/#create-a-database) for more information)

			So even if you once created a database before,
			if the database is empty, this data source will fail to read the database.
		`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: mdutils.FormatSchemaDescription(
					`
						Resource identifier.
						
						ID has a value with a format of the following:

						%s
					`,
					mdutils.CodeBlock("", "databases/<database>"),
				),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the database",
				Required:            true,
			},
		},
	}
}

func (d *DatabaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	config, diags := resourceconfig.FromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.config = config
}

func (d *DatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	client := d.config.Client.WithContext(ctx).WithLogger(d.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.AddError(errornames.MongoClientError, err.Error())
			return
		}

		var data DatabaseDataSourceModel

		// Read Terraform prior state data into the model
		resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform data source read operation
		resp.Diagnostics.Append(dataSourceRead(client, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save updated data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	})
}
