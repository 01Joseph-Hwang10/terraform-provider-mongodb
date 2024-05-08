// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"context"
	"encoding/json"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
	errorutils "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/utils"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type StringIsJSONValidator struct {
	description         string
	markdownDescription string
	optional            bool
}

func (v StringIsJSONValidator) Description(_ context.Context) string {
	return v.description
}

func (v StringIsJSONValidator) MarkdownDescription(_ context.Context) string {
	return v.markdownDescription
}

func (v StringIsJSONValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() && !v.optional {
		resp.Diagnostics.AddError(errornames.InvalidJSONInput, "JSON value is required.")
		return
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &parsed); err != nil {
		resp.Diagnostics.AddError(
			errornames.InvalidJSONInput,
			errorutils.NewInvalidJSONInputError(
				err,
				req.ConfigValue.ValueString(),
			),
		)
	}
}

func StringIsJSON() validator.String {
	return StringIsJSONValidator{
		description:         "The value must be a valid JSON string.",
		markdownDescription: "The value must be a valid JSON string.",
		optional:            false,
	}
}

func StringIsOptionalJSON() validator.String {
	return StringIsJSONValidator{
		description:         "The value must be a valid JSON string.",
		markdownDescription: "The value must be a valid JSON string.",
		optional:            true,
	}
}
