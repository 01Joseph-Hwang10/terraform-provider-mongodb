// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewUnexpectedResourceConfigurationType(expected string, received string) *UnexpectedResourceConfigurationType {
	return &UnexpectedResourceConfigurationType{
		expected: expected,
		received: received,
	}
}

type UnexpectedResourceConfigurationType struct {
	expected string
	received string
}

func (e *UnexpectedResourceConfigurationType) Error() string {
	return fmt.Sprintf(
		"Expected %s, got: %s. Please report this issue to the provider developers.",
		e.expected,
		e.received,
	)
}

func (e *UnexpectedResourceConfigurationType) Name() string {
	return "Unexpected Resource Configuration Type"
}

func (e *UnexpectedResourceConfigurationType) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
