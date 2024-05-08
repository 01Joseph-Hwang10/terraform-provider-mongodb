// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index_test

import (
	"fmt"
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/provider"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/acc"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexDataSource(t *testing.T) {
	mongolocal.RunWithServer(t, func(server *mongolocal.MongoLocal) {
		logger := server.Logger()

		resp := acc.PreTestAccIndexDataSource(server, logger)

		logger.Info("running the test...")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acc.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactoriesWithProviderConfig(&provider.Config{Logger: logger}),
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: acc.WithProviderConfig(fmt.Sprintf(`
						data "mongodb_database_index" "test" {
							database = "test-database"
							collection = "test-collection"
							index_name = "%s"
						}
					`, resp.IndexName), server.URI()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "id", fmt.Sprintf("databases/test-database/collections/test-collection/indexes/%s", resp.IndexName)),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "index_name", resp.IndexName),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "database", "test-database"),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "collection", "test-collection"),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "field", "test-field"),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "direction", "1"),
						resource.TestCheckResourceAttr("data.mongodb_database_index.test", "unique", "false"),
					),
				},
			},
		})
	})
}
