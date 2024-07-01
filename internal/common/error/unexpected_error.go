// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import "github.com/hashicorp/terraform-plugin-framework/diag"

func NewUnexpectedError(err error) *UnexpectedError {
	return &UnexpectedError{
		err: err,
	}
}

type UnexpectedError struct {
	err error
}

func (e *UnexpectedError) Error() string {
	return e.err.Error()
}

func (e *UnexpectedError) Name() string {
	return "Unexpected Error"
}

func (e *UnexpectedError) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
