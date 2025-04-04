// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongoclient

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Collection struct {
	name       string
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
	logger     *zap.Logger
}

func (d *Database) Collection(name string) *Collection {
	collection := d.database.Collection(name)
	return &Collection{
		name:       name,
		client:     d.client,
		database:   d.database,
		collection: collection,
		ctx:        d.ctx,
		logger:     d.logger,
	}
}

func (c *Collection) Name() string {
	return c.name
}

func (c *Collection) Client() *MongoClient {
	return c.Database().Client()
}

func (c *Collection) Database() *Database {
	return &Database{
		name:     c.database.Name(),
		client:   c.client,
		database: c.database,
		ctx:      c.ctx,
		logger:   c.logger,
	}
}

func (c *Collection) WithContext(ctx context.Context) *Collection {
	c.ctx = ctx
	return c
}

func (c *Collection) WithLogger(logger *zap.Logger) *Collection {
	c.logger = logger
	return c
}

func (c *Collection) Exists() (bool, error) {
	names, err := c.database.ListCollectionNames(
		c.ctx,
		bson.D{{Key: "name", Value: c.name}},
	)
	if err != nil {
		return false, err
	}

	return len(names) > 0, nil
}

func (c *Collection) EnsureExistance() error {
	// Check if the collection exists
	exists, err := c.Exists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Create the collection
	if err := c.database.CreateCollection(c.ctx, c.name); err != nil {
		return err
	}

	return nil
}

func (c *Collection) Drop() error {
	return c.collection.Drop(c.ctx)
}

func (c *Collection) IsEmpty() (bool, error) {
	count, err := c.collection.CountDocuments(c.ctx, bson.D{})
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (c *Collection) Find(filter Document) (Documents, error) {
	bsonFilter, err := filter.ToBson()
	if err != nil {
		return nil, err
	}

	cursor, err := c.collection.Find(c.ctx, bsonFilter)
	if err != nil {
		return nil, err
	}

	var documents Documents

	for cursor.Next(c.ctx) {
		var document Document
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

type FindByIdOptions struct {
	IncludeId bool
}

func (c *Collection) FindById(id string, opts *FindByIdOptions) (Document, error) {
	// Convert the id to an ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Retrieve the document
	var document Document
	filter := bson.D{{Key: "_id", Value: oid}}
	if err := c.collection.FindOne(c.ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	// If opts.IncludeId is not true, exclude the id from the document
	includeId := opts != nil && opts.IncludeId
	if !includeId {
		delete(document, "_id")
	}

	return document, nil
}

func (c *Collection) InsertOne(document Document) (string, error) {
	bsonDoc, err := document.ToBson()
	if err != nil {
		return "", err
	}

	res, err := c.collection.InsertOne(c.ctx, bsonDoc)
	if err != nil {
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to convert InsertedID to ObjectID")
	}
	return oid.Hex(), nil
}

func (c *Collection) UpdateByID(id string, update Document) error {
	bsonDoc, err := update.ToBson()
	if err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	_, err = c.collection.UpdateOne(c.ctx, filter, bson.D{{Key: "$set", Value: bsonDoc}})
	return err
}

func (c *Collection) DeleteByID(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	_, err = c.collection.DeleteOne(c.ctx, filter)
	return err
}
