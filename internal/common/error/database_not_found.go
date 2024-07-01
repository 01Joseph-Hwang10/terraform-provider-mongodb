// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewDatabaseNotFound(name string) *DatabaseNotFound {
	return &DatabaseNotFound{
		name: name,
	}
}

type DatabaseNotFound struct {
	name string
}

func (e *DatabaseNotFound) Error() string {
	return fmt.Sprintf("Database %s not found", e.name)
}

func (e *DatabaseNotFound) Name() string {
	return "Database Not Found"
}

func (e *DatabaseNotFound) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
