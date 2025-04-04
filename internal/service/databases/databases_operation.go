// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package databases

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

func dataSourceRead(client *mongoclient.MongoClient, data *DatabasesDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get the list of databases
	ctx := client.Context()
	names, err := client.Client().ListDatabaseNames(ctx, bson.D{})

	if err != nil {
		diags.Append(errs.NewMongoClientError(err).ToDiagnostic())
		return diags
	}

	// Filter the databases by name
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

	// Map the databases to the output format
	var databases []attr.Value
	for _, name := range matched {
		database, errs := basetypes.NewObjectValue(
			DatabaseElementType.AttrTypes,
			map[string]attr.Value{
				"id":   basetypes.NewStringValue(fmt.Sprintf("databases/%s", name)),
				"name": basetypes.NewStringValue(name),
			},
		)
		if errs != nil {
			diags.Append(errs...)
			return diags
		}
		databases = append(databases, database)
	}

	// Set the databases attribute
	v, errs := basetypes.NewListValue(DatabaseElementType, databases)
	if errs != nil {
		diags.Append(errs...)
		return diags
	}
	data.Databases = v

	return diags
}
