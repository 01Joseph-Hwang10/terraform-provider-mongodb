// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package database_test

import (
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseDataSource(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
			if err != nil {
				logger.Sugar().Fatalf("failed to create a client: %v", err)
			}

			logger.Info("creating a database to test the data source")
			database := client.Database("test-database")
			if err := database.Collection(mongoclient.PlaceholderCollectionName).EnsureExistance(); err != nil {
				logger.Sugar().Fatalf("failed to create a collection: %v", err)
			}
		})

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: acc.WithProviderConfig(`
						data "mongodb_database" "test" {
							name = "test-database"
						}
						`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.mongodb_database.test", "id", "databases/test-database"),
					),
				},
			},
		})
	})
}
