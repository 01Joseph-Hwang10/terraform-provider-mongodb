// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mdutils

import (
	"fmt"
	"html"
	"strings"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/indent"
)

func FormatResourceDescription(description string, args ...any) string {
	description = fmt.Sprintf(description, args...)
	description = strings.Trim(description, "\n")
	description = indent.Sanitize(description, indent.ProjectTabSize)
	return description
}

func FormatSchemaDescription(description string, args ...any) string {
	formatter := NewDescriptionFormatter()

	description = fmt.Sprintf(description, args...)
	description = strings.Trim(description, "\n")
	description = indent.Sanitize(description, indent.ProjectTabSize)
	description = formatter.Format(description)
	description = strings.ReplaceAll(description, "\n", " ")
	return description
}

func CodeBlock(language string, code string) string {
	code = indent.Sanitize(code, indent.ProjectTabSize)
	code = strings.Trim(code, "\n")
	code = html.EscapeString(code)
	code = strings.ReplaceAll(code, "\n", "<br />")

	langClass := ""
	if language != "" {
		langClass = fmt.Sprintf("language-%s", language)
	}

	return fmt.Sprintf(
		"<pre><code class=\"%s\">%s</code></pre>",
		langClass,
		code,
	)
}

func InlineCodeBlock(code string) string {
	return fmt.Sprintf(
		"<code>%s</code>",
		html.EscapeString(code),
	)
}
