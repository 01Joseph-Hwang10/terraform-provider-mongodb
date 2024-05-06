package resourceutils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ResourceId struct {
	database   string
	collection string
	document   string
	index      string
}

func NewId(id string) (*ResourceId, error) {
	isUnderDatabase := strings.Contains(id, "databases/")
	isUnderCollection := strings.Contains(id, "collections/")
	isUnderDocument := strings.Contains(id, "documents/")
	isUnderIndex := strings.Contains(id, "indexes/")

	id = strings.Replace(id, "databases/", "", -1)
	id = strings.Replace(id, "collections/", "", -1)
	id = strings.Replace(id, "documents/", "", -1)
	id = strings.Replace(id, "indexes/", "", -1)

	parts := strings.Split(id, "/")
	if len(parts) == 0 || len(parts) > 3 {
		return nil, errors.New("import ID must be in the format databases/<database_name>/collections/<collection_name>/... ")
	}

	database := ""
	collection := ""
	document := ""
	index := ""
	if isUnderDatabase {
		database = parts[0]
	}
	if isUnderCollection {
		collection = parts[1]
	}
	if isUnderDocument {
		document = parts[2]
	}
	if isUnderIndex {
		index = parts[3]
	}

	return &ResourceId{
		database:   database,
		collection: collection,
		document:   document,
		index:      index,
	}, nil
}

func (r *ResourceId) Collection() string {
	return r.collection
}

func (r *ResourceId) Database() string {
	return r.database
}

func (r *ResourceId) Document() string {
	return r.document
}

func (r *ResourceId) Index() string {
	return r.index
}

func (r *ResourceId) String() string {
	id := ""
	if r.database != "" {
		id += fmt.Sprintf("databases/%s", r.database)
	}
	if r.collection != "" {
		id += fmt.Sprintf("/collections/%s", r.collection)
	}
	if r.document != "" {
		id += fmt.Sprintf("/documents/%s", r.document)
	}
	if r.index != "" {
		id += fmt.Sprintf("/indexes/%s", r.index)
	}
	return id
}

func (r *ResourceId) TerraformString() basetypes.StringValue {
	return basetypes.NewStringValue(r.String())
}
