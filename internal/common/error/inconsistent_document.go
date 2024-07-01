// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewInconsistentDocument(base string, compare string, diff string) *InconsistentDocument {
	return &InconsistentDocument{
		base:    base,
		compare: compare,
		diff:    diff,
	}
}

type InconsistentDocument struct {
	base    string
	compare string
	diff    string
}

func (e *InconsistentDocument) Error() string {
	return fmt.Sprintf(`
Document is inconsistent with the data source.

Base: %s

Compare: %s

Diff: %s
		`,
		e.base,
		e.compare,
		e.diff,
	)
}

func (e *InconsistentDocument) Name() string {
	return "Inconsistent Document"
}

func (e *InconsistentDocument) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
