// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package collection

import (
	"fmt"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceid "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/id"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func CheckExistance(database *mongoclient.Database, name string, diags *diag.Diagnostics) *mongoclient.Collection {
	collection := database.Collection(name)
	exists, err := database.Exists()
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return nil
	}
	if !exists {
		diags.Append(
			errs.NewCollectionNotFound(collection.Name()).ToDiagnostic(),
		)
		return nil
	}
	return collection
}

func CreateResourceId(database basetypes.StringValue, name basetypes.StringValue) (basetypes.StringValue, error) {
	id, err := resourceid.New(fmt.Sprintf("databases/%s/collections/%s", database.ValueString(), name.ValueString()))
	if err != nil {
		return basetypes.NewStringNull(), err
	}
	return id.TerraformString(), nil
}

func dataSourceRead(client *mongoclient.MongoClient, data *CollectionDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check if the database exists
	database := database.CheckExistance(client, data.Database.ValueString(), &diags)
	if diags.HasError() {
		return diags
	}

	// Check if the collection exists
	CheckExistance(database, data.Name.ValueString(), &diags)
	if diags.HasError() {
		return diags
	}

	// Set resource Id
	resourceId, err := CreateResourceId(data.Database, data.Name)
	if err != nil {
		diags.Append(
			errs.NewInvalidResourceConfiguration(err.Error()).ToDiagnostic(),
		)
		return diags
	}
	data.Id = resourceId

	return diags
}

func resourceRead(client *mongoclient.MongoClient, data *CollectionResourceModel) diag.Diagnostics {
	// Type cast the resource data to data source data
	d := &CollectionDataSourceModel{
		Name:     data.Name,
		Id:       data.Id,
		Database: data.Database,
	}

	// Read the data source
	diags := dataSourceRead(client, d)

	// Convert back to resource data
	data.Id = d.Id
	data.Name = d.Name
	data.Database = d.Database

	return diags
}

func resourceCreate(client *mongoclient.MongoClient, data *CollectionResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check if the database exists
	database := database.CheckExistance(client, data.Database.ValueString(), &diags)
	if diags.HasError() {
		return diags
	}

	// Create the collection
	collection := database.Collection(data.Name.ValueString())
	if err := collection.EnsureExistance(); err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}

	// Perform the read operation
	diags.Append(resourceRead(client, data)...)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceDelete(client *mongoclient.MongoClient, data *CollectionResourceModel) diag.Diagnostics {
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
	collection := database.Collection(data.Name.ValueString())
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

	// Check if the collection is empty
	isEmpty, err := collection.IsEmpty()
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}
	if !isEmpty && !data.ForceDestroy.ValueBool() {
		diags.Append(
			errs.NewCollectionNotEmpty(collection.Name()).ToDiagnostic(),
		)
		return diags
	}

	// Delete the collection
	if err := collection.Drop(); err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}

	return diags
}
