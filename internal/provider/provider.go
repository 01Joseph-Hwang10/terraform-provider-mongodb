// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/collection"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/document"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/index"
)

// Ensure MongoProvider satisfies various provider interfaces.
var _ provider.Provider = &MongoProvider{}
var _ provider.ProviderWithFunctions = &MongoProvider{}

// MongoProvider defines the provider implementation.
type MongoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MongoProviderModel describes the provider data model.
type MongoProviderModel struct {
	URI string `tfsdk:"uri"`
}

func (p *MongoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mongodb"
	resp.Version = p.version
}

func (p *MongoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				MarkdownDescription: "URI to connect to the MongoDB server.",
				Required:            true,
			},
		},
	}
}

func (p *MongoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MongoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := mongoclient.New(ctx, &mongoclient.Config{
		URI: data.URI,
	})

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *MongoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		database.NewDatabaseResource,
		collection.NewCollectionResource,
		document.NewDocumentResource,
		index.NewIndexResource,
	}
}

func (p *MongoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		database.NewDatabaseDataSource,
		collection.NewCollectionDataSource,
		document.NewDocumentDataSource,
		index.NewIndexDataSource,
	}
}

func (p *MongoProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MongoProvider{
			version: version,
		}
	}
}
