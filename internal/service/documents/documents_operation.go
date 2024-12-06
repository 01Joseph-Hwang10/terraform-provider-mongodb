// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package documents

import (
	"encoding/json"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/collection"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func dataSourceRead(client *mongoclient.MongoClient, data *DocumentsDataSourceModel) diag.Diagnostics {
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

	// Parse the filter
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(data.Filter.ValueString()), &filter); err != nil {
		diags.Append(
			errs.NewInvalidResourceConfiguration(err.Error()).ToDiagnostic(),
		)
		return diags
	}

	// Read documents from the collection with the filter
	documents, err := collection.Find(filter)
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}
	encoded, err := json.Marshal(documents)
	if err != nil {
		diags.Append(
			errs.NewUnexpectedError(err).ToDiagnostic(),
		)
		return diags
	}
	data.Documents = basetypes.NewStringValue(string(encoded))

	return diags
}
