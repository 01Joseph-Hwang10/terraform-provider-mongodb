// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package documents_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/provider"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDocumentsDataSource(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
			if err != nil {
				logger.Sugar().Fatalf("failed to create a client: %v", err)
			}

			logger.Info("creating a document to test document data source")

			collection := client.Database("test-database").Collection("test-collection")

			if _, err = collection.InsertOne(mongoclient.Document{"key": "value-1"}); err != nil {
				logger.Sugar().Fatalf("failed to insert a document: %v", err)
			}
			if _, err = collection.InsertOne(mongoclient.Document{"key": "value-2"}); err != nil {
				logger.Sugar().Fatalf("failed to insert a document: %v", err)
			}
			if _, err = collection.InsertOne(mongoclient.Document{"key": "value-3"}); err != nil {
				logger.Sugar().Fatalf("failed to insert a document: %v", err)
			}
		})

		logger.Info("running the test...")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactoriesWithProviderConfig(&provider.Config{Logger: logger}),
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: acc.WithProviderConfig(`
						data "mongodb_database_documents" "test" {
							database = "test-database"
							collection = "test-collection"
							filter = jsonencode({
								key = { "$eq" = "value-2" }
							})
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrWith("data.mongodb_database_documents.test", "documents", func(value string) error {
							var decoded []map[string]interface{}
							if err := json.Unmarshal([]byte(value), &decoded); err != nil {
								return err
							}

							if len(decoded) != 1 {
								return fmt.Errorf("expected 1 document, got %d", len(decoded))
							}

							if decoded[0]["key"] != "value-2" {
								return fmt.Errorf("expected key to be 'value-2', got %s", decoded[0]["key"])
							}

							return nil
						}),
					),
				},
			},
		})
	})
}
