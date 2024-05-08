// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package collection_test

import (
	"regexp"
	"testing"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCollectionResource_Lifecycle(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
			if err != nil {
				logger.Sugar().Fatalf("failed to create a client: %v", err)
			}

			logger.Info("creating a database to test the collection resource")
			database := client.Database("test-database")
			if err := database.Collection(mongoclient.PlaceholderCollectionName).EnsureExistance(); err != nil {
				logger.Sugar().Fatalf("failed to create a collection: %v", err)
			}
		})

		logger.Info("running the test...")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_collection" "test" {
							database = "test-database"
							name = "test-collection"
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_collection.test", "force_destroy", "false"),
						resource.TestCheckResourceAttr("mongodb_database_collection.test", "name", "test-collection"),
						resource.TestCheckResourceAttr("mongodb_database_collection.test", "database", "test-database"),
						resource.TestCheckResourceAttr("mongodb_database_collection.test", "id", "databases/test-database/collections/test-collection"),
					),
				},
				// ImportState testing
				{
					ResourceName:            "mongodb_database_collection.test",
					ImportStateId:           "databases/test-database/collections/test-collection",
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"force_destroy"},
				},
				// Update and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_collection" "test" {
							database = "test-database"
							name = "test-collection"
							force_destroy = true
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_collection.test", "force_destroy", "true"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})
}

func TestAccCollectionResource_ForceDestroy(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
			if err != nil {
				logger.Sugar().Fatalf("failed to create a client: %v", err)
			}

			logger.Info("creating a document to test force_destroy option...")

			if _, err := client.Database("test-database").Collection("test-collection").InsertOne(mongoclient.Document{"key": "value"}); err != nil {
				logger.Sugar().Fatalf("failed to insert a document: %v", err)
			}
		})

		logger.Info("running the test...")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Import the resource
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_collection" "test" {
							database = "test-database"
							name = "test-collection"
						}
					`, server.URI()),
					ResourceName:       "mongodb_database_collection.test",
					ImportStateId:      "databases/test-database/collections/test-collection",
					ImportState:        true,
					ImportStatePersist: true,
				},
				// Try to destroy the resource
				{
					Config:      acc.ProviderConfig(server.URI()),
					Destroy:     true,
					ExpectError: regexp.MustCompile(errornames.CollectionNotEmpty),
				},
				// Update the resource to force destroy
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_collection" "test" {
							database = "test-database"
							name = "test-collection"
							force_destroy = true
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_collection.test", "id", "databases/test-database/collections/test-collection"),
						resource.TestCheckResourceAttr("mongodb_database_collection.test", "name", "test-collection"),
						resource.TestCheckResourceAttr("mongodb_database_collection.test", "force_destroy", "true"),
					),
				},
				// Destroy testing automatically occurs in TestCase
				//
				// Note that `terraform-plugin-testing` expects the resource to be successfully destroyed
				// after the last step. See the issue below for more details:
				//     https://github.com/hashicorp/terraform-plugin-sdk/issues/609
			},
		})
	})
}
