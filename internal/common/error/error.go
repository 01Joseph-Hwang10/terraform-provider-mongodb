// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errs

import "github.com/hashicorp/terraform-plugin-framework/diag"

type DiagnosableError interface {
	error
	Name() string
	ToDiagnostic() diag.Diagnostic
}
