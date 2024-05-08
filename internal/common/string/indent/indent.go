// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package indent

import (
	"strings"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/replace"
)

const (
	ProjectTabSize = 2

	indentUndefined = -1
)

// Removes leading indents in respect to the shortest indent.
func Sanitize(s string, tabSize int) string {
	// Replace every tab with spaces.
	s = strings.ReplaceAll(s, "\t", strings.Repeat(" ", tabSize))

	// Split the string into lines.
	lines := strings.Split(s, "\n")

	// Find the shortest indent.
	minIndent := indentUndefined
	for _, line := range lines {
		// Skip empty lines.
		if isEmpty(line) {
			continue
		}

		// Get the indent size and update the minIndent.
		indent := getIndentSize(line)
		if minIndent == indentUndefined || indent < minIndent {
			minIndent = indent
		}
	}

	// Remove the shortest indent.
	var sanitized string
	for _, line := range lines {
		// Skip empty lines.
		if isEmpty(line) {
			sanitized += "\n"
			continue
		}

		// Remove the indent.
		sanitized += line[minIndent:] + "\n"
	}

	return sanitized
}

func getIndentSize(s string) int {
	for i, r := range s {
		if r != ' ' {
			return i
		}
	}
	return len(s)
}

func isEmpty(s string) bool {
	return replace.NewChain(
		replace.NewReplacement(" ", ""),
		replace.NewReplacement("\t", ""),
	).Apply(s) == ""
}
