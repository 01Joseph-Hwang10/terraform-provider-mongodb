// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewDocumentNotFound(id string) *DocumentNotFound {
	return &DocumentNotFound{
		id: id,
	}
}

type DocumentNotFound struct {
	id string
}

func (e *DocumentNotFound) Error() string {
	return fmt.Sprintf("Document with ID %s not found", e.id)
}

func (e *DocumentNotFound) Name() string {
	return "Document Not Found"
}

func (e *DocumentNotFound) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
