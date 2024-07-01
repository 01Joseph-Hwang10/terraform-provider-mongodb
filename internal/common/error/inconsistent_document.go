// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewInconsistentDocument(data string) *InconsistentDocument {
	return &InconsistentDocument{
		data: data,
	}
}

type InconsistentDocument struct {
	data string
}

func (e *InconsistentDocument) Error() string {
	return fmt.Sprintf("Document is inconsistent with the data source: %s", e.data)
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
