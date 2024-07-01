// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewCollectionNotEmpty(name string) *CollectionNotEmpty {
	return &CollectionNotEmpty{
		name: name,
	}
}

type CollectionNotEmpty struct {
	name string
}

func (e *CollectionNotEmpty) Error() string {
	return fmt.Sprintf(
		"Collection %s contains data. Set force_destroy to true to delete the collection.",
		e.name,
	)
}

func (e *CollectionNotEmpty) Name() string {
	return "Collection Not Empty"
}

func (e *CollectionNotEmpty) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
