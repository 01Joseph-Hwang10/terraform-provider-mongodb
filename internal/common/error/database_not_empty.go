// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewDatabaseNotEmpty(name string) *DatabaseNotEmpty {
	return &DatabaseNotEmpty{
		name: name,
	}
}

type DatabaseNotEmpty struct {
	name string
}

func (e *DatabaseNotEmpty) Error() string {
	return fmt.Sprintf(
		"Database %s contains collections, set force_destroy to true to destroy the database",
		e.name,
	)
}

func (e *DatabaseNotEmpty) Name() string {
	return "Database Not Empty"
}

func (e *DatabaseNotEmpty) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
