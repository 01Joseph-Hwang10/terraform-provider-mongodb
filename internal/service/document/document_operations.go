// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package document

import (
	"encoding/json"
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resourceutils"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/collection"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func CreateResourceId(database basetypes.StringValue, collection basetypes.StringValue, documentId basetypes.StringValue) (basetypes.StringValue, error) {
	id, err := resourceutils.NewId(
		fmt.Sprintf("databases/%s/collections/%s/documents/%s", database.ValueString(), collection.ValueString(), documentId.ValueString()),
	)
	if err != nil {
		return basetypes.NewStringNull(), err
	}
	return id.TerraformString(), nil
}

func resourceRead(client *mongoclient.MongoClient, data *DocumentResourceModel) diag.Diagnostics {
	// Type cast the resource data to data source data
	d := &DocumentDataSourceModel{
		Id:         data.Id,
		Database:   data.Database,
		Collection: data.Collection,
		DocumentId: data.DocumentId,
		Document:   data.Document,
	}

	// Read the data source
	diags := dataSourceRead(client, d)

	// Convert back to resource data
	data.Id = d.Id
	data.Database = d.Database
	data.Collection = d.Collection
	data.DocumentId = d.DocumentId
	data.Document = d.Document

	return diags
}

func dataSourceRead(client *mongoclient.MongoClient, data *DocumentDataSourceModel) diag.Diagnostics {
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

	// Read the document
	document, err := collection.FindById(data.DocumentId.ValueString())
	if err != nil {
		diags.AddError("Client Error", err.Error())
		return diags
	}

	// Set document
	encoded, err := json.Marshal(document)
	if err != nil {
		diags.AddError("Invalid JSON Input", err.Error())
		return diags
	}
	data.Document = basetypes.NewStringValue(string(encoded))

	// Set resource Id
	resourceId, err := CreateResourceId(data.Database, data.Collection, data.DocumentId)
	if err != nil {
		diags.AddError("Invalid configuration", err.Error())
		return diags
	}
	data.Id = resourceId

	return diags
}
