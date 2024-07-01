// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewCollectionNotFound(name string) *CollectionNotFound {
	return &CollectionNotFound{
		name: name,
	}
}

type CollectionNotFound struct {
	name string
}

func (e *CollectionNotFound) Error() string {
	return fmt.Sprintf("Collection %s not found", e.name)
}

func (e *CollectionNotFound) Name() string {
	return "Collection Not Found"
}

func (e *CollectionNotFound) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
