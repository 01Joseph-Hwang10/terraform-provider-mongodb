// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package document

import (
	"encoding/json"
	"fmt"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceid "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/id"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/collection"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/service/database"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/wI2L/jsondiff"
)

func CreateResourceId(database basetypes.StringValue, collection basetypes.StringValue, documentId basetypes.StringValue) (basetypes.StringValue, error) {
	id, err := resourceid.New(
		fmt.Sprintf("databases/%s/collections/%s/documents/%s", database.ValueString(), collection.ValueString(), documentId.ValueString()),
	)
	if err != nil {
		return basetypes.NewStringNull(), err
	}
	return id.TerraformString(), nil
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
	document, err := collection.FindById(data.DocumentId.ValueString(), nil)
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}
	if document == nil {
		diags.Append(
			errs.NewDocumentNotFound(data.DocumentId.ValueString()).ToDiagnostic(),
		)
		return diags
	}

	// Set document
	encoded, err := document.ToEJson()
	if err != nil {
		diags.Append(
			errs.NewEJsonParseError(err).ToDiagnostic(),
		)
		return diags
	}
	data.Document = basetypes.NewStringValue(encoded)

	// Set resource Id
	resourceId, err := CreateResourceId(data.Database, data.Collection, data.DocumentId)
	if err != nil {
		diags.Append(
			errs.NewInvalidResourceConfiguration(err.Error()).ToDiagnostic(),
		)
		return diags
	}
	data.Id = resourceId

	return diags
}

func resourceRead(client *mongoclient.MongoClient, r *DocumentResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	d := DocumentDataSourceModel{
		Database:   r.Database,
		Collection: r.Collection,
		DocumentId: r.DocumentId,
		Document:   basetypes.NewStringNull(),
		Id:         basetypes.NewStringNull(),
	}

	diags.Append(dataSourceRead(client, &d)...)
	if diags.HasError() {
		return diags
	}

	// Assign retrieved document to the data model
	// if document in data model is not set
	if r.Document.IsNull() {
		r.Document = d.Document
	}

	// Validate document consistency
	// Document consistency validation is only performed
	// when SyncWithDatabase is enabled.
	if r.SyncWithDatabase.ValueBool() {
		diags.Append(validateDocumentConsistency(r, &d)...)
		if diags.HasError() {
			return diags
		}
	}

	// Set resource Id
	resourceId, err := CreateResourceId(r.Database, r.Collection, r.DocumentId)
	if err != nil {
		diags.Append(
			errs.NewInvalidResourceConfiguration(err.Error()).ToDiagnostic(),
		)
		return diags
	}
	r.Id = resourceId

	return diags
}

func validateDocumentConsistency(r *DocumentResourceModel, d *DocumentDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Parse document read from the data source
	var document mongoclient.Document
	rawDocument := d.Document.ValueString()
	if err := json.Unmarshal([]byte(rawDocument), &document); err != nil {
		diags.Append(
			errs.NewInvalidJSONDocument(err.Error(), rawDocument).ToDiagnostic(),
		)
		return diags
	}

	// Compare received document from the state
	// with the document in the data source
	//
	// As terraform-plugin-framework v2
	// strictly requires data consistency among the operations,
	// we need to use `data.Document` value if it exists.
	//
	// See more details at the links below:
	// - https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/data-consistency-errors
	// - https://discuss.hashicorp.com/t/is-it-possible-to-have-statefunc-like-behavior-with-the-plugin-framework/58377
	var expected mongoclient.Document
	rawExpected := r.Document.ValueString()
	if err := json.Unmarshal([]byte(rawExpected), &expected); err != nil {
		diags.Append(
			errs.NewInvalidJSONDocument(err.Error(), rawExpected).ToDiagnostic(),
		)
		return diags
	}
	patch, err := jsondiff.Compare(document, expected)
	if err != nil {
		diags.Append(
			errs.NewUnexpectedError(err).ToDiagnostic(),
		)
		return diags
	}
	if patch.String() != "" {
		diags.Append(
			errs.NewInconsistentDocument(
				rawDocument,
				rawExpected,
				patch.String(),
			).ToDiagnostic(),
		)
		return diags
	}

	return diags
}

func resourceCreate(client *mongoclient.MongoClient, data *DocumentResourceModel) diag.Diagnostics {
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

	// Create the document
	var document mongoclient.Document
	rawDocument := data.Document.ValueString()
	if err := json.Unmarshal([]byte(rawDocument), &document); err != nil {
		diags.Append(
			errs.NewInvalidJSONDocument(err.Error(), rawDocument).ToDiagnostic(),
		)
		return diags
	}
	oid, err := collection.InsertOne(document)
	if err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}

	data.DocumentId = basetypes.NewStringValue(oid)

	// Perform a read operation to get the document
	diags.Append(resourceRead(client, data)...)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceUpdate(client *mongoclient.MongoClient, data *DocumentResourceModel) diag.Diagnostics {
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

	// Update the document
	var document mongoclient.Document
	rawDocument := data.Document.ValueString()
	if err := json.Unmarshal([]byte(rawDocument), &document); err != nil {
		diags.Append(
			errs.NewInvalidJSONDocument(err.Error(), rawDocument).ToDiagnostic(),
		)
		return diags
	}
	if err := collection.UpdateByID(data.DocumentId.ValueString(), document); err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}

	// Perform a read operation to get the document
	diags.Append(resourceRead(client, data)...)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceDelete(client *mongoclient.MongoClient, data *DocumentResourceModel) diag.Diagnostics {
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

	// Delete the document
	if err := collection.DeleteByID(data.DocumentId.ValueString()); err != nil {
		diags.Append(
			errs.NewMongoClientError(err).ToDiagnostic(),
		)
		return diags
	}

	return diags
}
