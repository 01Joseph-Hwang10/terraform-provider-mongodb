// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package testenv

import "os"

func IsDebug() bool {
	return os.Getenv("DEBUG") != ""
}

func ExecRoot() string {
	return os.Getenv("EXEC_ROOT")
}
