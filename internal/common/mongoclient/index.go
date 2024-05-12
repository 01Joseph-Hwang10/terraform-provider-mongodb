// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongoclient

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type SanitizedIndexSpec struct {
	Name      string
	Field     string
	Direction int
	Unique    bool
}

type Index struct {
	name       string
	field      string
	direction  int
	unique     bool
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
	logger     *zap.Logger
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
		logger:     c.logger,
	}
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
		logger:     c.logger,
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
		logger:     i.logger,
	}
}

func (i *Index) WithContext(ctx context.Context) *Index {
	i.ctx = ctx
	return i
}

func (i *Index) GetSpec() (*SanitizedIndexSpec, error) {
	// Check if the index exists
	specs, err := i.collection.Indexes().ListSpecifications(i.ctx)
	if err != nil {
		return nil, err
	}

	var spec *SanitizedIndexSpec
	if i.name == "" {
		// If the index name is not set, find the index by field,
		// direction, and unique constraint
		spec, err = i.findIndexByField(i.field, i.direction, i.unique, specs)
		if err != nil {
			return nil, err
		}
		if spec == nil {
			return nil, nil
		}
	} else {
		// If the index name is set, find the index by name
		spec, err = i.findIndexByName(i.name, specs)
		if err != nil {
			return nil, err
		}
		if spec == nil {
			return nil, nil
		}
	}
	return spec, nil
}

func (i *Index) findIndexByName(name string, specs []*mongo.IndexSpecification) (*SanitizedIndexSpec, error) {
	for _, spec := range specs {
		if spec.Name == name {
			// Get field name and direction
			elements, err := spec.KeysDocument.Elements()
			if err != nil {
				return nil, err
			}
			var field string
			var direction int
			for _, element := range elements {
				field = element.Key()
				direction = int(element.Value().Int32())
				break
			}

			// Get unique constraint
			unique := false
			if spec.Unique != nil && *spec.Unique {
				unique = true
			}

			// Return the index spec
			return &SanitizedIndexSpec{
				Name:      spec.Name,
				Field:     field,
				Direction: direction,
				Unique:    unique,
			}, nil
		}
	}
	return nil, nil
}

func (i *Index) findIndexByField(field string, direction int, unique bool, specs []*mongo.IndexSpecification) (*SanitizedIndexSpec, error) {
	for _, spec := range specs {
		uniqueMatches := (!unique && spec.Unique == nil) || (spec.Unique != nil && unique == *spec.Unique)
		if uniqueMatches {
			elements, err := spec.KeysDocument.Elements()
			if err != nil {
				return nil, err
			}
			for _, element := range elements {
				if element.Key() == field && element.Value().Int32() == int32(direction) {
					return &SanitizedIndexSpec{
						Name:      spec.Name,
						Field:     field,
						Direction: direction,
						Unique:    unique,
					}, nil
				}
			}
		}
	}
	return nil, nil
}

func (i *Index) Hydrate(spec *SanitizedIndexSpec) *Index {
	i.name = spec.Name
	i.field = spec.Field
	i.direction = spec.Direction
	i.unique = spec.Unique
	return i
}

func (i *Index) Exists() (bool, error) {
	spec, err := i.GetSpec()
	if err != nil {
		return false, err
	}
	return spec != nil, nil
}

func (i *Index) EnsureExistance() error {
	// Check if the index exists
	spec, err := i.GetSpec()
	if err != nil {
		return err
	}
	if spec != nil {
		return nil
	}

	// Create the index
	if i.field == "" || i.direction == 0 {
		return errors.New("unexpected error: field and direction must be set")
	}
	filter := bson.D{{Key: i.field, Value: i.direction}}
	name, err := i.collection.Indexes().CreateOne(i.ctx, mongo.IndexModel{
		Keys:    filter,
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
