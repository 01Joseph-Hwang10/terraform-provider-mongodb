// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package collection_test

import (
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCollectionDataSource(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
			if err != nil {
				logger.Sugar().Fatalf("failed to create a client: %v", err)
			}

			logger.Info("creating a database and a collection to test the collection resource")
			database := client.Database("test-database")
			if err := database.Collection("test-collection").EnsureExistance(); err != nil {
				logger.Sugar().Fatalf("failed to create a collection: %v", err)
			}
		})

		logger.Info("running the test...")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: acc.WithProviderConfig(`
						data "mongodb_database_collection" "test" {
							database = "test-database"
							name = "test-collection"
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.mongodb_database_collection.test", "id", "databases/test-database/collections/test-collection"),
					),
				},
			},
		})
	})
}
