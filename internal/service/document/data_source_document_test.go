// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package document_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAccDocumentDataSource(t *testing.T) {
	mongolocal.WithMongoLocal(t, func(server *mongolocal.MongoLocal) {
		// Create document for testing
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(server.URI()))
		if err != nil {
			t.Fatalf("failed to connect to MongoDB: %s", err)
			return
		}
		defer client.Disconnect(context.Background())

		collection := client.Database("test-database").Collection("test-collection")
		res, err := collection.InsertOne(context.Background(), map[string]interface{}{
			"name": "test-document",
		})
		if err != nil {
			t.Fatalf("failed to insert document: %s", err)
			return
		}
		oid := res.InsertedID.(primitive.ObjectID).Hex()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: acc.WithProviderConfig(fmt.Sprintf(`
						data "mongodb_database_document" "test" {
							database = "test-database"
							collection = "test-collection"
							document_id = "%s"
						}
					`, oid), server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.mongodb_database_document.test", "id", fmt.Sprintf("databases/test-database/collections/test-collection/documents/%s", oid)),
					),
				},
			},
		})
	})
}
