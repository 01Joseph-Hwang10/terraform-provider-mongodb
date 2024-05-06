// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package database_test

import (
	"context"
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAccDatabaseDataSource(t *testing.T) {
	mongolocal.WithMongoLocal(t, func(server *mongolocal.MongoLocal) {
		// Create database for testing
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(server.URI()))
		if err != nil {
			t.Fatalf("failed to connect to MongoDB: %s", err)
			return
		}
		defer client.Disconnect(context.Background())

		// Create dummy collection to ensure database exists
		database := client.Database("test-database")
		if err := database.CreateCollection(context.Background(), "test-collection"); err != nil {
			t.Fatalf("failed to create collection: %s", err)
			return
		}

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
