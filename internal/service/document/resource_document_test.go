// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package document_test

import (
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDocumentResource(t *testing.T) {
	mongolocal.WithMongoLocal(t, func(server *mongolocal.MongoLocal) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_document" "test" {
							database = "test-database"
							collection = "test-collection"
							document = "{ \"name\": \"test-document\" }"
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_document.test", "collection", "test-collection"),
						resource.TestCheckResourceAttr("mongodb_database_document.test", "database", "test-database"),
						resource.TestCheckResourceAttr("mongodb_database_document.test", "document", "{ \"name\": \"test-document\" }"),
					),
				},
				// ImportState testing
				{
					ResourceName: "mongodb_database_document.test",
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						document_id := s.RootModule().Resources["mongodb_database_document.test"].Primary.RawState.GetAttr("document_id").AsString()
						return "databases/test-database/collections/test-collection/documents/" + document_id, nil
					},
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: acc.WithProviderConfig(`
						resource "mongodb_database_collection" "test" {
							database = "test-database"
							collection = "test-collection"
							document = "{ \"name\": \"test-document\", \"age\": 20}"
						}
					`, server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_document.test", "document", "{ \"name\": \"test-document\", \"age\": 20}"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})
}
