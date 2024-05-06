// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resourceutils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func CheckExistance(client *mongoclient.MongoClient, name string, diags *diag.Diagnostics) *mongoclient.Database {
	database := client.Database(name)
	exists, err := database.Exists()
	if err != nil {
		diags.AddError("Client Error", err.Error())
		return nil
	}
	if !exists {
		diags.AddError("Database Not Found", fmt.Sprintf("Database %s not found", name))
		return nil
	}
	return database
}

func CreateResourceId(name basetypes.StringValue) (basetypes.StringValue, error) {
	id, err := resourceutils.NewId(fmt.Sprintf("databases/%s", name.ValueString()))
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
		diags.AddError("Invalid configuration", err.Error())
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
