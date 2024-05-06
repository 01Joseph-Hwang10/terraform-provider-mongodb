// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongoclient

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI string
}

type MongoClient struct {
	config *Config
	client *mongo.Client
	ctx    context.Context
}

func New(ctx context.Context, config *Config) *MongoClient {
	return &MongoClient{
		config: config,
	}
}

func (c *MongoClient) Connect() error {
	// Create a new client
	client, err := mongo.Connect(c.ctx, options.Client().ApplyURI(c.config.URI))
	if err != nil {
		return err
	}
	// Check the connection
	if err := client.Ping(c.ctx, nil); err != nil {
		return err
	}
	c.client = client
	return nil
}

func (c *MongoClient) Disconnect() {
	if err := c.client.Disconnect(c.ctx); err != nil {
		panic(err)
	}
	c.client = nil
}

func (c *MongoClient) IsConnected() bool {
	return c.client != nil
}

func (c *MongoClient) WithContext(ctx context.Context) *MongoClient {
	c.ctx = ctx
	return c
}
