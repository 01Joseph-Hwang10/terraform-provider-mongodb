// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import "github.com/hashicorp/terraform-plugin-framework/diag"

func NewInvalidInputValue(description string) *InvalidInputValue {
	return &InvalidInputValue{
		description: description,
	}
}

type InvalidInputValue struct {
	description string
}

func (e *InvalidInputValue) Error() string {
	return e.description
}

func (e *InvalidInputValue) Name() string {
	return "Invalid Input Value"
}

func (e *InvalidInputValue) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
