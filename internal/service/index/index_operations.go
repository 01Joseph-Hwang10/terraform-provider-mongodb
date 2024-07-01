// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index

import (
	"fmt"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceid "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/id"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/collection"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func CreateResourceId(database basetypes.StringValue, collection basetypes.StringValue, index basetypes.StringValue) (basetypes.StringValue, error) {
	id, err := resourceid.New(fmt.Sprintf("databases/%s/collections/%s/indexes/%s", database.ValueString(), collection.ValueString(), index.ValueString()))
	if err != nil {
		return basetypes.NewStringNull(), err
	}
	return id.TerraformString(), nil
}

func dataSourceRead(client *mongoclient.MongoClient, data *IndexDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if !client.IsConnected() {
		diags.Append(
			errs.NewClientIsNotConnected().ToDiagnostic(),
		)
		return diags
	}

	// Check if the database exists
	database := database.CheckExistance(client, data.Database.ValueString(), &diags)
	if diags.HasError() {
		return diags
	}

	// Check if the collection exists
	collection := collection.CheckExistance(database, data.Collection.ValueString(), &diags)
	if diags.HasError() {
		return diags
	}

	// Check if the index exists
	index := collection.Index(data.IndexName.ValueString())
	spec, err := index.GetSpec()
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}
	if spec == nil {
		diags.Append(
			errs.NewIndexNotFound(data.IndexName.ValueString()).ToDiagnostic(),
		)
		return diags
	}

	// Hydrate index information
	index.Hydrate(spec)

	// Set resource Id
	resourceId, err := CreateResourceId(data.Database, data.Collection, data.IndexName)
	if err != nil {
		diags.Append(
			errs.NewInvalidResourceConfiguration(err.Error()).ToDiagnostic(),
		)
		return diags
	}

	data.Id = resourceId
	data.Collection = basetypes.NewStringValue(index.Collection().Name())
	data.Database = basetypes.NewStringValue(index.Database().Name())
	data.IndexName = basetypes.NewStringValue(index.Name())
	data.Field = basetypes.NewStringValue(index.Field())
	data.Direction = basetypes.NewInt64Value(int64(index.Direction()))
	data.Unique = basetypes.NewBoolValue(index.Unique())

	return diags
}

func resourceRead(client *mongoclient.MongoClient, data *IndexResourceModel) diag.Diagnostics {
	// Type cast the resource data to data source data
	d := &IndexDataSourceModel{
		Id:         data.Id,
		Database:   data.Database,
		Collection: data.Collection,
		IndexName:  data.IndexName,
		Field:      data.Field,
		Direction:  data.Direction,
		Unique:     data.Unique,
	}

	// Read the data source
	diags := dataSourceRead(client, d)

	// Convert back to resource data
	data.Id = d.Id
	data.Database = d.Database
	data.Collection = d.Collection
	data.IndexName = d.IndexName
	data.Field = d.Field
	data.Direction = d.Direction
	data.Unique = d.Unique

	return diags
}

func resourceCreate(client *mongoclient.MongoClient, data *IndexResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check if the database exists
	database := database.CheckExistance(client, data.Database.ValueString(), &diags)
	if diags.HasError() {
		return diags
	}

	// Check if the collection exists
	collection := collection.CheckExistance(database, data.Collection.ValueString(), &diags)
	if diags.HasError() {
		return diags
	}

	// Create the index
	field := data.Field.ValueString()
	direction := int(data.Direction.ValueInt64())
	unique := data.Unique.ValueBool()
	index := collection.IndexFromField(field, direction, unique)
	if err := index.EnsureExistance(); err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}

	// Get the index name
	name := index.Name()

	// If the index name is not set,
	// infer it from the field name and direction
	if name == "" {
		name = fmt.Sprintf("%s_%d", field, direction)
	}

	// Set index name
	data.IndexName = basetypes.NewStringValue(name)

	// Perform read operation
	diags.Append(resourceRead(client, data)...)

	return diags
}

func resourceDelete(client *mongoclient.MongoClient, data *IndexResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check if the database exists
	database := client.Database(data.Database.ValueString())
	exists, err := database.Exists()
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}
	if !exists {
		// We don't need to check if the collection exists,
		// as the database doesn't exist
		return diags
	}

	// Check if the collection exists
	collection := database.Collection(data.Collection.ValueString())
	exists, err = collection.Exists()
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}
	if !exists {
		// Collection doesn't exist, nothing to delete
		return diags
	}

	// If force destroy is not set, fail the deletion
	if !data.ForceDestroy.ValueBool() {
		diags.Append(
			errs.NewIndexDeletionForbidden().ToDiagnostic(),
		)
		return diags
	}

	// Delete the index
	index := collection.Index(data.IndexName.ValueString())
	if err := index.Drop(); err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}

	return diags
}
