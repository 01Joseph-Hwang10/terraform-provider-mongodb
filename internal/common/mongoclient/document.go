// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongoclient

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type Documents []Document

type Document map[string]interface{}

// Converts extended json to bson
func (d *Document) ToBson() (bson.D, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	var bsonDoc bson.D
	if err := bson.UnmarshalExtJSON(b, false, &bsonDoc); err != nil {
		return nil, err
	}

	return bsonDoc, nil
}

func (d *Document) ToEJson() (string, error) {
	encoded, err := bson.MarshalExtJSON(d, false, false)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}
