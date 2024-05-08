// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongolocal

import "fmt"

func watcherScript(parent int, child int) string {
	return fmt.Sprintf(
		"while kill -0 %d; do "+
			"sleep 1; "+
			"done; "+
			"kill -9 %d ",
		parent, child)
}
