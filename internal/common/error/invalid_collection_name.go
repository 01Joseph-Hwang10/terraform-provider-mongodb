// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewInvalidCollectionName(name string) *InvalidCollectionName {
	return &InvalidCollectionName{
		name: name,
	}
}

type InvalidCollectionName struct {
	name string
}

func (e *InvalidCollectionName) Error() string {
	return fmt.Sprintf(
		"Cannot use database name %s.",
		e.name,
	)
}

func (e *InvalidCollectionName) Name() string {
	return "Invalid Database Name"
}

func (e *InvalidCollectionName) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
