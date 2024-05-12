// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongoclient

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Config struct {
	URI string
}

type MongoClient struct {
	config *Config
	client *mongo.Client
	logger *zap.Logger
	ctx    context.Context
}

func New(ctx context.Context, config *Config) *MongoClient {
	return &MongoClient{
		config: config,
		ctx:    ctx,
	}
}

func FromURI(uri string) *MongoClient {
	return New(context.Background(), &Config{
		URI: uri,
	})
}

func (c *MongoClient) Client() *mongo.Client {
	return c.client
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
	// We don't handle the error here because we are disconnecting
	// the client anyway.
	_ = c.client.Disconnect(c.ctx)

	// Set the client to nil
	c.client = nil
}

func (c *MongoClient) IsConnected() bool {
	return c.client != nil
}

func (c *MongoClient) WithContext(ctx context.Context) *MongoClient {
	c.ctx = ctx
	return c
}

func (c *MongoClient) WithLogger(logger *zap.Logger) *MongoClient {
	c.logger = logger
	return c
}

func (c *MongoClient) Run(callback func(client *MongoClient, err error)) {
	err := c.Connect()
	callback(c, err)
	c.Disconnect()
}
