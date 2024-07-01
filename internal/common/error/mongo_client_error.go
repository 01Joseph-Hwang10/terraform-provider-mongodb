// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import "github.com/hashicorp/terraform-plugin-framework/diag"

func NewMongoClientError(err error) *MongoClientError {
	return &MongoClientError{
		err: err,
	}
}

type MongoClientError struct {
	err error
}

func (e *MongoClientError) Error() string {
	return e.err.Error()
}

func (e *MongoClientError) Name() string {
	return "Mongo Client Error"
}

func (e *MongoClientError) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
