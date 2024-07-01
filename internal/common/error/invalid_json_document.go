// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewInvalidJSONDocument(description string, json string) *InvalidJSONDocument {
	return &InvalidJSONDocument{
		description: description,
		json:        json,
	}
}

type InvalidJSONDocument struct {
	description string
	json        string
}

func (e *InvalidJSONDocument) Error() string {
	return fmt.Sprintf("%s (received: %s)", e.description, e.json)
}

func (e *InvalidJSONDocument) Name() string {
	return "Invalid JSON Document"
}

func (e *InvalidJSONDocument) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
