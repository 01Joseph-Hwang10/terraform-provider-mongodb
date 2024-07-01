// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import "github.com/hashicorp/terraform-plugin-framework/diag"

func NewEJsonParseError(err error) *EJsonParseError {
	return &EJsonParseError{
		err: err,
	}
}

type EJsonParseError struct {
	err error
}

func (e *EJsonParseError) Error() string {
	return e.err.Error()
}

func (e *EJsonParseError) Name() string {
	return "EJson Parse Error"
}

func (e *EJsonParseError) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
