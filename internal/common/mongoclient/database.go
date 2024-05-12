// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongoclient

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	// Placeholder collection name for explicit database creation.
	PlaceholderCollectionName = "__terraform_provider_mongodb"
)

type Database struct {
	name     string
	client   *mongo.Client
	database *mongo.Database
	ctx      context.Context
	logger   *zap.Logger
}

func (c *MongoClient) Database(name string) *Database {
	database := c.client.Database(name)
	return &Database{
		name:     name,
		client:   c.client,
		database: database,
		ctx:      c.ctx,
		logger:   c.logger,
	}
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) Client() *MongoClient {
	return &MongoClient{
		config: nil,
		client: d.client,
		ctx:    d.ctx,
		logger: d.logger,
	}
}

func (d *Database) WithContext(ctx context.Context) *Database {
	d.ctx = ctx
	return d
}

func (d *Database) WithLogger(logger *zap.Logger) *Database {
	d.logger = logger
	return d
}

func (d *Database) EnsureExistance() error {
	// Check if the database exists
	exists, err := d.Exists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Create a dummy collection to ensure the database is created
	if err := d.database.CreateCollection(d.ctx, PlaceholderCollectionName); err != nil {
		return err
	}

	return nil
}

func (d *Database) Exists() (bool, error) {
	names, err := d.client.ListDatabaseNames(d.ctx, bson.D{})
	if err != nil {
		return false, err
	}

	for _, name := range names {
		if name == d.name {
			return true, nil
		}
	}

	return false, nil
}

func (d *Database) Drop() error {
	return d.database.Drop(d.ctx)
}

func (d *Database) IsEmpty() (bool, error) {
	collections, err := d.database.ListCollectionNames(d.ctx, bson.D{})
	if err != nil {
		return false, err
	}
	isEmpty := (len(collections) == 0 ||
		(len(collections) == 1 && collections[0] == PlaceholderCollectionName))

	return isEmpty, nil
}
