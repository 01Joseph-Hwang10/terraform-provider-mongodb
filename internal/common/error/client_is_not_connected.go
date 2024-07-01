// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewClientIsNotConnected() *ClientIsNotConnected {
	return &ClientIsNotConnected{}
}

type ClientIsNotConnected struct{}

func (e *ClientIsNotConnected) Error() string {
	return "Mongo Client is not connected"
}

func (e *ClientIsNotConnected) Name() string {
	return "Client Is Not Connected"
}

func (e *ClientIsNotConnected) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
