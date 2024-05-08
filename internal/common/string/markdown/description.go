// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mdutils

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type DescriptionFormatter struct {
	parser   *parser.Parser
	renderer *html.Renderer
}

func NewDescriptionFormatter() *DescriptionFormatter {
	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	// Create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Create a new DescriptionFormatter
	return &DescriptionFormatter{
		parser:   p,
		renderer: renderer,
	}
}

func (f *DescriptionFormatter) Format(description string) string {
	// Parse markdown to HTML
	ast := f.parser.Parse([]byte(description))
	return string(markdown.Render(ast, f.renderer))
}
