// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package database_test

import (
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseResource(t *testing.T) {
	mongolocal.WithMongoLocal(t, func(server *mongolocal.MongoLocal) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database" "test" {
							name = "test-database"
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database.test", "force_destroy", "false"),
						resource.TestCheckResourceAttr("mongodb_database.test", "name", "test-database"),
						resource.TestCheckResourceAttr("mongodb_database.test", "id", "databases/test-database"),
					),
				},
				// ImportState testing
				{
					ResourceName:      "mongodb_database.test",
					ImportStateId:     "databases/test-database",
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database" "test" {
							name = "test-database"
							force_destroy = true
						}
						`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database.test", "force_destroy", "true"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})
}
