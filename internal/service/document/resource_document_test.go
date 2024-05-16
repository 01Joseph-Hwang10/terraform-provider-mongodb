// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package document_test

import (
	"fmt"
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/replace"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/provider"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDocumentResource_Lifecycle(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
			if err != nil {
				logger.Sugar().Fatalf("failed to create a client: %v", err)
			}

			logger.Info("creating a database and a collection to test the collection resource")
			database := client.Database("test-database")
			if err := database.Collection("test-collection").EnsureExistance(); err != nil {
				logger.Sugar().Fatalf("failed to create a collection: %v", err)
			}
		})

		tfFormat := replace.NewChain(
			replace.NewReplacement("\n", ""),
			replace.NewReplacement("\t", ""),
		)
		compFormat := tfFormat.Copy().Extend(
			replace.NewReplacement("\\\"", "\""),
		)

		firstDocument := `
			{
				\"name\":\"test-document\",
				\"with\":{
					\"some\":\"nested\",
					\"fields\":\"and\",
					\"arrays\":	[
						1,
						2,
						{
							\"three\":3
						}
					],
					\"date\":{
						\"$date\":\"2021-01-01T00:00:00Z\"
					}
				}
			}
		`
		updatedDocument := `
			{
				\"name\":\"test-document\",
				\"with\":\"some-changed-value\"
			}
		`

		logger.Info("running the test...")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactoriesWithProviderConfig(&provider.Config{Logger: logger}),
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: acc.WithProviderConfig(documentResource(
						"test-database",
						"test-collection",
						tfFormat.Apply(firstDocument),
					), server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_document.test", "collection", "test-collection"),
						resource.TestCheckResourceAttr("mongodb_database_document.test", "database", "test-database"),
						resource.TestCheckResourceAttr("mongodb_database_document.test", "document", compFormat.Apply(firstDocument)),
					),
				},
				// ImportState testing
				{
					ResourceName: "mongodb_database_document.test",
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						// Load resource data from state as JSON
						resources, err := acc.LoadResources(s.RootModule().Resources)
						if err != nil {
							return "", err
						}

						// Select the document ID
						document_id := resources["mongodb_database_document.test"].(map[string]interface{})["primary"].(map[string]interface{})["attributes"].(map[string]interface{})["document_id"].(string) //nolint:forcetypeassert

						return "databases/test-database/collections/test-collection/documents/" + document_id, nil
					},
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"document"},
				},
				// Update and Read testing
				{
					Config: acc.WithProviderConfig(documentResource(
						"test-database",
						"test-collection",
						tfFormat.Apply(updatedDocument),
					), server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("mongodb_database_document.test", "document", compFormat.Apply(updatedDocument)),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})
}

func documentResource(database string, collection string, document string) string {
	return fmt.Sprintf(`
		resource "mongodb_database_document" "test" {
			database = "%s"
			collection = "%s"
			document = "%s"
		}
	`, database, collection, document)
}
