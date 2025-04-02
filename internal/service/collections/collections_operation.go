// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package collections

import (
	"fmt"
	"regexp"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"go.mongodb.org/mongo-driver/bson"
)

func dataSourceRead(client *mongoclient.MongoClient, data *CollectionsDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get the list of collections
	ctx := client.Context()
	database := data.Database.ValueString()
	names, err := client.Client().Database(database).ListCollectionNames(ctx, bson.D{})

	if err != nil {
		diags.Append(errs.NewMongoClientError(err).ToDiagnostic())
		return diags
	}

	// Filter the collections by name
	pattern, err := regexp.Compile(data.Name.ValueString())
	if err != nil {
		diags.Append(errs.NewUnexpectedError(err).ToDiagnostic())
		return diags
	}

	var matched []string
	hasFilter := !(data.Name.IsNull() || data.Name.ValueString() == "")
	if hasFilter {
		for _, name := range names {
			if pattern.MatchString(name) {
				matched = append(matched, name)
			}
		}
	} else {
		matched = names
	}

	var filtered []string
	for _, name := range matched {
		if name != mongoclient.PlaceholderCollectionName {
			filtered = append(filtered, name)
		}
	}

	// Map the collections to the output format
	var collections []attr.Value
	for _, name := range filtered {
		database, errs := basetypes.NewObjectValue(
			CollectionElementType.AttrTypes,
			map[string]attr.Value{
				"id":       basetypes.NewStringValue(fmt.Sprintf("databases/%s/collections/%s", database, name)),
				"database": basetypes.NewStringValue(database),
				"name":     basetypes.NewStringValue(name),
			},
		)
		if errs != nil {
			diags.Append(errs...)
			return diags
		}
		collections = append(collections, database)
	}

	// Set the collections attribute
	v, errs := basetypes.NewListValue(CollectionElementType, collections)
	if errs != nil {
		diags.Append(errs...)
		return diags
	}
	data.Collections = v

	return diags
}
