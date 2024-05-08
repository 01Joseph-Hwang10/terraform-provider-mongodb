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

func NewReplacement(old string, new string) replacement {
	return replacement{
		old: old,
		new: new,
	}
}

type replaceChain struct {
	replacements []replacement
}

func NewChain(replacements ...replacement) *replaceChain {
	return &replaceChain{
		replacements: replacements,
	}
}

func (c *replaceChain) Apply(s string) string {
	for _, r := range c.replacements {
		s = strings.Replace(s, r.old, r.new, -1)
	}
	return s
}

func (c *replaceChain) Copy() *replaceChain {
	return &replaceChain{
		replacements: append([]replacement{}, c.replacements...),
	}
}

func (c *replaceChain) Extend(replacements ...replacement) *replaceChain {
	c.replacements = append(c.replacements, replacements...)
	return c
}
