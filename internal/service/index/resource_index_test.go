// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index_test

import (
	"fmt"
	"regexp"
	"testing"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
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

		acc.PreTestAccIndexResource(server, logger)

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
						resource.TestCheckResourceAttr("mongodb_database_index.test", "force_destroy", "false"),
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

						index_name := resources["mongodb_database_index.test"].(map[string]interface{})["primary"].(map[string]interface{})["attributes"].(map[string]interface{})["index_name"].(string) //nolint:forcetypeassert
						return "databases/test-database/collections/test-collection/indexes/" + index_name, nil
					},
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"force_destroy"},
				},
				// Update and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_index" "test" {
							database = "test-database"
							collection = "test-collection"
							field = "test-field"
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

func TestAccIndexResource_Variant(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		acc.PreTestAccIndexResource(server, logger)

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
							direction = -1
							unique = true
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_index.test", "database", "test-database"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "collection", "test-collection"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "field", "test-field"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "direction", "-1"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "unique", "true"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "force_destroy", "false"),
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

						index_name := resources["mongodb_database_index.test"].(map[string]interface{})["primary"].(map[string]interface{})["attributes"].(map[string]interface{})["index_name"].(string) //nolint:forcetypeassert
						return "databases/test-database/collections/test-collection/indexes/" + index_name, nil
					},
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"force_destroy"},
				},
				// Update and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_index" "test" {
							database = "test-database"
							collection = "test-collection"
							field = "test-field"
							direction = -1
							unique = true
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

func TestAccIndexResource_ForceDestroy(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		resp := acc.PreTestAccIndexDataSource(server, logger)

		logger.Info("running the test...")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactoriesWithProviderConfig(&provider.Config{Logger: logger}),
			Steps: []resource.TestStep{
				// Import the resource
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_index" "test" {
							database = "test-database"
							collection = "test-collection"
							field = "test-field"
						}
					`, server.URI()),
					ResourceName:       "mongodb_database_index.test",
					ImportStateId:      fmt.Sprintf("databases/test-database/collections/test-collection/indexes/%s", resp.IndexName),
					ImportState:        true,
					ImportStatePersist: true,
				},
				// Try to destroy the resource
				{
					Config:      acc.ProviderConfig(server.URI()),
					Destroy:     true,
					ExpectError: regexp.MustCompile(errornames.IndexDeletionForbidden),
				},
				// Update the resource to force destroy
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_index" "test" {
							database = "test-database"
							collection = "test-collection"
							field = "test-field"
							force_destroy = true
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_index.test", "database", "test-database"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "collection", "test-collection"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "field", "test-field"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "direction", "1"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "unique", "false"),
						resource.TestCheckResourceAttr("mongodb_database_index.test", "force_destroy", "true"),
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
