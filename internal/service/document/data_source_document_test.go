// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package document_test

import (
	"fmt"
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/provider"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDocumentDataSource(t *testing.T) {
	t.Parallel()
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		var oid string
		mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
			if err != nil {
				logger.Sugar().Fatalf("failed to create a client: %v", err)
			}

			logger.Info("creating a document to test document data source")

			oid, err = client.Database("test-database").Collection("test-collection").InsertOne(mongoclient.Document{"key": "value"})
			if err != nil {
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
