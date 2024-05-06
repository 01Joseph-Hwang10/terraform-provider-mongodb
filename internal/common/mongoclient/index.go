// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongoclient

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Index struct {
	name       string
	field      string
	direction  int
	unique     bool
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

func (c *Collection) IndexFromField(field string, direction int, unique bool) *Index {
	return &Index{
		name:       "",
		field:      field,
		direction:  direction,
		unique:     unique,
		client:     c.client,
		database:   c.database,
		collection: c.collection,
		ctx:        c.ctx,
	}
}

func (c *Collection) Index(name string) *Index {
	return &Index{
		name:       name,
		field:      "",
		direction:  0,
		unique:     false,
		client:     c.client,
		database:   c.database,
		collection: c.collection,
		ctx:        c.ctx,
	}
}

func (i *Index) Name() string {
	return i.name
}

func (i *Index) Field() string {
	return i.field
}

func (i *Index) Direction() int {
	return i.direction
}

func (i *Index) Unique() bool {
	return i.unique
}

func (i *Index) Client() *MongoClient {
	return i.Collection().Client()
}

func (i *Index) Database() *Database {
	return i.Collection().Database()
}

func (i *Index) Collection() *Collection {
	return &Collection{
		name:       i.collection.Name(),
		client:     i.client,
		database:   i.database,
		collection: i.collection,
		ctx:        i.ctx,
	}
}

func (i *Index) WithContext(ctx context.Context) *Index {
	i.ctx = ctx
	return i
}

func (i *Index) Exists() (bool, error) {
	// Check if the index exists
	specs, err := i.collection.Indexes().ListSpecifications(i.ctx)
	if err != nil {
		return false, err
	}

	var spec *mongo.IndexSpecification
	if i.name == "" {
		spec = findIndexByField(i.field, i.direction, i.unique, specs)

		// Update the index name
		i.name = spec.Name
	} else {
		spec = findIndexByName(i.name, specs)

		// Update the index properties
		elements, err := spec.KeysDocument.Elements()
		if err != nil {
			return false, err
		}
		i.field = elements[0].Key()
		i.direction = int(elements[0].Value().Int32())
		i.unique = *spec.Unique
	}
	return spec != nil, nil
}

func findIndexByName(name string, specs []*mongo.IndexSpecification) *mongo.IndexSpecification {
	for _, spec := range specs {
		if spec.Name == name {
			return spec
		}
	}
	return nil
}

func findIndexByField(field string, direction int, unique bool, specs []*mongo.IndexSpecification) *mongo.IndexSpecification {
	for _, spec := range specs {
		if *spec.Unique == unique {
			elements, err := spec.KeysDocument.Elements()
			if err != nil {
				return nil
			}
			for _, element := range elements {
				if element.Key() == field && element.Value().Int32() == int32(direction) {
					return spec
				}
			}
		}
	}
	return nil
}

func (i *Index) EnsureExistance() error {
	// Check if the index exists
	exists, err := i.Exists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Create the index
	if i.field == "" || i.direction == 0 {
		return errors.New("unexpected error: field and direction must be set")
	}
	name, err := i.collection.Indexes().CreateOne(i.ctx, mongo.IndexModel{
		Keys:    bson.M{i.field: i.direction},
		Options: options.Index().SetUnique(i.unique),
	})
	if err != nil {
		return err
	}

	// Store the index name
	i.name = name
	return nil
}

func (i *Index) Drop() error {
	_, err := i.collection.Indexes().DropOne(i.ctx, i.name)
	return err
}
