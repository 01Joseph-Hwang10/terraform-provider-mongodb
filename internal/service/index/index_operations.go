// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index

import (
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resourceutils"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/collection"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func CheckExistance(index *mongoclient.Index, diags *diag.Diagnostics) *mongoclient.Index {
	exists, err := index.Exists()
	if err != nil {
		diags.AddError("Client Error", err.Error())
		return nil
	}
	if !exists {
		diags.AddError("Index not found", fmt.Sprintf("Index on field %s with direction %d not found", index.Field(), index.Direction()))
		return nil
	}
	return index
}

func CreateResourceId(database basetypes.StringValue, collection basetypes.StringValue, index basetypes.StringValue) (basetypes.StringValue, error) {
	id, err := resourceutils.NewId(fmt.Sprintf("databases/%s/collections/%s/indexes/%s", database.ValueString(), collection.ValueString(), index.ValueString()))
	if err != nil {
		return basetypes.NewStringNull(), err
	}
	return id.TerraformString(), nil
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

func dataSourceRead(client *mongoclient.MongoClient, data *IndexDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if !client.IsConnected() {
		diags.AddError("Client Error", "Client is not connected")
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
	index := collection.Index(data.IndexName.String())
	CheckExistance(index, &diags)
	if diags.HasError() {
		return diags
	}

	// Set resource Id
	resourceId, err := CreateResourceId(data.Database, data.Collection, data.IndexName)
	if err != nil {
		diags.AddError("Invalid configuration", err.Error())
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
