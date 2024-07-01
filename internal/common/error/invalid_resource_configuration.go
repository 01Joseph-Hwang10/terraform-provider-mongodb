// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import "github.com/hashicorp/terraform-plugin-framework/diag"

func NewInvalidResourceConfiguration(message string) *InvalidResourceConfiguration {
	return &InvalidResourceConfiguration{
		message: message,
	}
}

type InvalidResourceConfiguration struct {
	message string
}

func (e *InvalidResourceConfiguration) Error() string {
	return e.message
}

func (e *InvalidResourceConfiguration) Name() string {
	return "Invalid Resource Configuration"
}

func (e *InvalidResourceConfiguration) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
