// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index

import (
	"context"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const description = "direction must be 1 or -1"

type isDirection struct {
	validator.Int64
}

func IsDirection() validator.Int64 {
	return &isDirection{}
}

func (v *isDirection) Description(context.Context) string {
	return description
}

func (v *isDirection) MarkdownDescription(context.Context) string {
	return description
}

func (v *isDirection) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueInt64()
	if value == 1 || value == -1 {
		return
	}

	resp.Diagnostics.AddError(errornames.InvalidInputValue, description)
}
