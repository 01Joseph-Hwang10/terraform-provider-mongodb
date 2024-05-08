// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package indent_test

import (
	"testing"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/indent"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/replace"
)

type TestCase struct {
	name     string
	input    string
	expected string
}

func TestIndentSanitize(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name: "test-case-1",
			input: `
						hello
						    this string expects to be indented
					  this string should not be indented
					`,
			expected: "\n" +
				"hello\n" +
				"    this string expects to be indented\n" +
				"this string should not be indented\n" +
				"\n",
		},
		{
			name: "test-case-2",
			input: `
								this string expects to be indented
					  this string should not be indented

						hello
					`,
			expected: "\n" +
				"    this string expects to be indented\n" +
				"this string should not be indented\n" +
				"\n" +
				"hello\n" +
				"\n",
		},
		{
			name: "test-case-3",
			input: "" +
				`
				Document to insert into the collection.

				The value of this attribute is a stringified JSON.
				Note that you should escape every double quote in the JSON string.

				In terraform, you can achieve this by simply using the 
				%s function:

				%s
			`,
			expected: "\n" +
				"Document to insert into the collection.\n" +
				"\n" +
				"The value of this attribute is a stringified JSON.\n" +
				"Note that you should escape every double quote in the JSON string.\n" +
				"\n" +
				"In terraform, you can achieve this by simply using the \n" +
				"%s function:\n" +
				"\n" +
				"%s\n" +
				"\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			output := indent.Sanitize(testCase.input, indent.ProjectTabSize)
			if output != testCase.expected {
				t.Errorf(
					"---EXPECTED:\n%s\n\n---ACTUAL:\n%s\n\n",
					formatTestString(testCase.expected),
					formatTestString(output),
				)
			}
		})
	}
}

func formatTestString(s string) string {
	return replace.NewChain(
		replace.NewReplacement("\t", "\\t"),
		replace.NewReplacement(" ", "\\s"),
		replace.NewReplacement("\n", "<EOL\n"),
	).Apply(s)
}
