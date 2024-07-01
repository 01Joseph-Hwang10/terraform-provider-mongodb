// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"fmt"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceid "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/id"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func CheckExistance(client *mongoclient.MongoClient, name string, diags *diag.Diagnostics) *mongoclient.Database {
	database := client.Database(name)
	exists, err := database.Exists()
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return nil
	}
	if !exists {
		diags.Append(
			errs.NewDatabaseNotFound(database.Name()).ToDiagnostic(),
		)
		return nil
	}
	return database
}

func CreateResourceId(name basetypes.StringValue) (basetypes.StringValue, error) {
	id, err := resourceid.New(fmt.Sprintf("databases/%s", name.ValueString()))
	if err != nil {
		return basetypes.NewStringNull(), err
	}
	return id.TerraformString(), nil
}

func dataSourceRead(client *mongoclient.MongoClient, data *DatabaseDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check if the database exists
	CheckExistance(client, data.Name.ValueString(), &diags)
	if diags.HasError() {
		return diags
	}

	// Set resource Id
	resourceId, err := CreateResourceId(data.Name)
	if err != nil {
		diags.Append(
			errs.NewInvalidResourceConfiguration(err.Error()).ToDiagnostic(),
		)
		return diags
	}
	data.Id = resourceId

	return diags
}

func resourceRead(client *mongoclient.MongoClient, data *DatabaseResourceModel) diag.Diagnostics {
	// Type cast the resource data to data source data
	d := &DatabaseDataSourceModel{
		Name: data.Name,
		Id:   data.Id,
	}

	// Read the data source
	diags := dataSourceRead(client, d)

	// Convert back to resource data
	data.Id = d.Id
	data.Name = d.Name

	return diags
}

func resourceCreate(client *mongoclient.MongoClient, data *DatabaseResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create the database
	database := client.Database(data.Name.ValueString())
	if err := database.EnsureExistance(); err != nil {
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

func resourceDelete(client *mongoclient.MongoClient, data *DatabaseResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	database := client.Database(data.Name.ValueString())

	// Check if the database is empty
	isEmpty, err := database.IsEmpty()
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}

	if !isEmpty && !data.ForceDestroy.ValueBool() {
		diags.Append(
			errs.NewDatabaseNotEmpty(database.Name()).ToDiagnostic(),
		)
		return diags
	}

	if err := database.Drop(); err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}

	return diags
}
