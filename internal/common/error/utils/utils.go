// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package errorutils

import "fmt"

func NewInvalidJSONInputError(err error, json string) string {
	return fmt.Sprintf(
		"%s (received: %s)",
		err.Error(),
		json,
	)
}
