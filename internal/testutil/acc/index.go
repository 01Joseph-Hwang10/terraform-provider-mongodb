// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package acc

import (
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/testutil/mongolocal"
	"go.uber.org/zap"
)

func PreTestAccIndexResource(server *mongolocal.MongoLocal, logger *zap.Logger) string {
	var oid string

	mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			logger.Sugar().Fatalf("failed to create a client: %v", err)
		}

		logger.Info("creating a document for the test")

		collection := client.Database("test-database").Collection("test-collection")
		oid, err = collection.InsertOne(mongoclient.Document{"test-field": "test-value"})
		if err != nil {
			logger.Sugar().Fatalf("failed to insert a document: %v", err)
		}
	})

	return oid
}

type preTestAccIndexDataSourceResponse struct {
	DocumentId string
	IndexName  string
}

func PreTestAccIndexDataSource(server *mongolocal.MongoLocal, logger *zap.Logger) *preTestAccIndexDataSourceResponse {
	var oid string
	var indexName string

	mongoclient.FromURI(server.URI()).Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			logger.Sugar().Fatalf("failed to create a client: %v", err)
		}

		logger.Info("creating a document for the test")

		collection := client.Database("test-database").Collection("test-collection")
		oid, err = collection.InsertOne(mongoclient.Document{"test-field": "test-value"})
		if err != nil {
			logger.Sugar().Fatalf("failed to insert a document: %v", err)
		}

		logger.Info("creating an index for the test")
		index := collection.IndexFromField("test-field", 1, false)
		if err := index.EnsureExistance(); err != nil {
			logger.Sugar().Fatalf("failed to create an index: %v", err)
		}
		indexName = index.Name()
	})

	return &preTestAccIndexDataSourceResponse{
		DocumentId: oid,
		IndexName:  indexName,
	}
}
