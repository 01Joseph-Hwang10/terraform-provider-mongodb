// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errornames

const (
	// MongoDB Errors.
	MongoClientError   = "Mongo Client Error"
	EJsonParseError    = "EJson Parse Error"
	DatabaseNotFound   = "Database Not Found"
	CollectionNotFound = "Collection Not Found"
	IndexNotFound      = "Index Not Found"
	DocumentNotFound   = "Document Not Found"

	// Terraform Configuration Errors.
	InvalidImportID                     = "Invalid Import ID"
	InvalidResourceConfiguration        = "Invalid Resource Configuration"
	UnexpectedResourceConfigurationType = "Unexpected Resource Configuration Type"

	// Terraform Resource Configuration Errors.
	DatabaseNotEmpty       = "Database Not Empty"
	CollectionNotEmpty     = "Collection Not Empty"
	IndexDeletionForbidden = "Index Deletion Forbidden"
	InvalidInputValue      = "Invalid Input Value"
	InvalidJSONInput       = "Invalid JSON Input"
	InconsistentDocument   = "Inconsistent Document"

	// Internal Errors.
	UnexpectedError = "Unexpected Error"
)
