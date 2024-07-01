// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package replace

import (
	"strings"
)

type replacement struct {
	old string
	new string
}

func NewReplacement(oldStr string, newStr string) replacement {
	return replacement{
		old: oldStr,
		new: newStr,
	}
}

type ReplaceChain struct {
	replacements []replacement
}

func NewChain(replacements ...replacement) *ReplaceChain {
	return &ReplaceChain{
		replacements: replacements,
	}
}

func (c *ReplaceChain) Apply(s string) string {
	for _, r := range c.replacements {
		s = strings.Replace(s, r.old, r.new, -1)
	}
	return s
}

func (c *ReplaceChain) Copy() *ReplaceChain {
	return &ReplaceChain{
		replacements: append([]replacement{}, c.replacements...),
	}
}

func (c *ReplaceChain) Extend(replacements ...replacement) *ReplaceChain {
	c.replacements = append(c.replacements, replacements...)
	return c
}
