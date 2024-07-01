// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewIndexNotFound(name string) *IndexNotFound {
	return &IndexNotFound{
		name: name,
	}
}

type IndexNotFound struct {
	name string
}

func (e *IndexNotFound) Error() string {
	return fmt.Sprintf("Index %s not found", e.name)
}

func (e *IndexNotFound) Name() string {
	return "Index Not Found"
}

func (e *IndexNotFound) ToDiagnostic() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		e.Name(),
		e.Error(),
	)
}
