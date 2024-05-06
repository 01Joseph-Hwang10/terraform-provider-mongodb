// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAccIndexDataSource(t *testing.T) {
	mongolocal.WithMongoLocal(t, func(server *mongolocal.MongoLocal) {
		// Create collection index for testing
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(server.URI()))
		if err != nil {
			t.Fatalf("failed to connect to MongoDB: %s", err)
			return
		}
		defer client.Disconnect(context.Background())

		collection := client.Database("test-database").Collection("test-collection")
		index_name, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.M{"test-field": 1},
		})
		if err != nil {
			t.Fatalf("failed to create index: %s", err)
			return
		}

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: acc.WithProviderConfig(fmt.Sprintf(`
						data "mongodb_database_index" "test" {
							database = "test-database"
							collection = "test-collection"
							index_name = "%s"
						}
					`, index_name), server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "id", fmt.Sprintf("databases/test-database/collections/test-collection/indexes/%s", index_name)),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "index_name", index_name),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "database", "test-database"),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "collection", "test-collection"),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "field", "test-field"),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "direction", "1"),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "unique", "false"),
					),
				},
			},
		})
	})
}
