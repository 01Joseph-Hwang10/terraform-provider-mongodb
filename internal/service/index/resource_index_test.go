// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index_test

import (
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIndexResource(t *testing.T) {
	mongolocal.WithMongoLocal(t, func(server *mongolocal.MongoLocal) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
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
						index_name := s.RootModule().Resources["mongodb_database_index.test"].Primary.RawState.GetAttr("index_name").AsString()
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
