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
	"go.uber.org/zap"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceconfig "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/config"
	mdutils "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/markdown"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/collection"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/document"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/index"
)

// Ensure MongoProvider satisfies various provider interfaces.
var _ provider.Provider = &MongoProvider{}
var _ provider.ProviderWithFunctions = &MongoProvider{}

type Config struct {
	Logger *zap.Logger
}

// MongoProvider defines the provider implementation.
type MongoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	// config configures provider behavior.
	config *Config
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
		MarkdownDescription: mdutils.FormatResourceDescription(
			`
				%s allows you to manage 
				MongoDB databases, collections, documents, and indexes.
			`,
			mdutils.InlineCodeBlock("01Joseph-Hwang10/terraform-provider-mongodb"),
		),

		Attributes: map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				MarkdownDescription: mdutils.FormatSchemaDescription(`
					URI to connect to the MongoDB server. 

					You should include valid username and password whose roles have the necessary permissions 
					for the operations you want to perform in the connection string. 
					
					Also, you should attach the options as a query string to the connection string 
					if you want to use it
				`),
				Required: true,
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

	// Prepare mongo client
	client := mongoclient.FromURI(data.URI).WithContext(ctx)

	// Prepare the logger
	logger, err := configureLogger(p)
	if err != nil {
		resp.Diagnostics.AddError(errornames.UnexpectedError, err.Error())
		return
	}

	providerData := &resourceconfig.ResourceConfig{
		Client: client,
		Logger: logger,
	}

	resp.ResourceData = providerData
	resp.DataSourceData = providerData
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
	// Set default config values
	config := &Config{
		Logger: nil,
	}

	return WithConfig(version, config)
}

func WithConfig(version string, config *Config) func() provider.Provider {
	return func() provider.Provider {
		return &MongoProvider{
			version: version,
			config:  config,
		}
	}
}
