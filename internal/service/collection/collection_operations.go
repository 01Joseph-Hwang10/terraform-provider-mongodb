// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package collection

import (
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resourceutils"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func CheckExistance(database *mongoclient.Database, name string, diags *diag.Diagnostics) *mongoclient.Collection {
	collection := database.Collection(name)
	exists, err := database.Exists()
	if err != nil {
		diags.AddError("Client Error", err.Error())
		return nil
	}
	if !exists {
		diags.AddError("Collection Not Found", fmt.Sprintf("Collection %s not found", name))
		return nil
	}
	return collection
}

func CreateResourceId(database basetypes.StringValue, name basetypes.StringValue) (basetypes.StringValue, error) {
	id, err := resourceutils.NewId(fmt.Sprintf("databases/%s/collections/%s", database.ValueString(), name.ValueString()))
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
		diags.AddError("Invalid configuration", err.Error())
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
