// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index_test

import (
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/provider"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIndexResource_Lifecycle(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
			if err != nil {
				logger.Sugar().Fatalf("failed to create a client: %v", err)
			}

			logger.Info("creating a document to test index resource")

			if _, err := client.Database("test-database").Collection("test-collection").InsertOne(mongoclient.Document{"test-field": "test-value"}); err != nil {
				logger.Sugar().Fatalf("failed to insert a document: %v", err)
			}
		})

		logger.Info("running the test...")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactoriesWithProviderConfig(&provider.Config{Logger: logger}),
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_index" "test" {
							database = "test-database"
							collection = "test-collection"
							field = "test-field"
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_index.test", "database", "test-database"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "collection", "test-collection"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "field", "test-field"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "direction", "1"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "unique", "false"),
					),
				},
				// ImportState testing
				{
					ResourceName: "mongodb_database_index.test",
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						resources, err := acc.LoadResources(s.RootModule().Resources)
						if err != nil {
							return "", err
						}

						index_name := resources["mongodb_database_index.test"].(map[string]interface{})["primary"].(map[string]interface{})["attributes"].(map[string]interface{})["index_name"].(string)
						return "databases/test-database/collections/test-collection/indexes/" + index_name, nil
					},
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_index" "test" {
							database = "test-database"
							collection = "test-collection"
							force_destroy = true
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_index.test", "force_destroy", "true"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})
}
