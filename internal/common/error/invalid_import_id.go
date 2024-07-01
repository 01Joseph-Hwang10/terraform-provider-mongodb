// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import "github.com/hashicorp/terraform-plugin-framework/diag"

func NewInvalidImportID(message string) *InvalidImportID {
	return &InvalidImportID{
		message: message,
	}
}

type InvalidImportID struct {
	message string
}

func (e *InvalidImportID) Error() string {
	return e.message
}

func (e *InvalidImportID) Name() string {
	return "Invalid Import ID"
}

func (e *InvalidImportID) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
