// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import "github.com/hashicorp/terraform-plugin-framework/diag"

func NewIndexDeletionForbidden() *IndexDeletionForbidden {
	return &IndexDeletionForbidden{}
}

type IndexDeletionForbidden struct{}

func (e *IndexDeletionForbidden) Error() string {
	return "Index deletion is not allowed by default. Set force_destroy to true to delete the index."
}

func (e *IndexDeletionForbidden) Name() string {
	return "Index Deletion Forbidden"
}

func (e *IndexDeletionForbidden) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
