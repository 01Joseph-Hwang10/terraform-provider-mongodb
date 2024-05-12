// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package resourceconfig

import (
	"fmt"

	errornames "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error/names"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.uber.org/zap"
)

type ResourceConfig struct {
	ClientConfig *mongoclient.Config
	Logger       *zap.Logger
}

func FromProviderData(data any) (config *ResourceConfig, diags diag.Diagnostics) {
	// Prevent panic if the provider has not been configured.
	if data == nil {
		return nil, diags
	}

	providerData, ok := data.(*ResourceConfig)
	if !ok {
		diags.AddError(
			errornames.UnexpectedResourceConfigurationType,
			fmt.Sprintf("Expected *resourceconfig.ResourceConfig, got: %T. Please report this issue to the provider developers.", data),
		)
		return nil, diags
	}

	if providerData.ClientConfig == nil {
		diags.AddError(
			errornames.UnexpectedResourceConfigurationType,
			"Expected *mongoclient.Config, got nil. Please report this issue to the provider developers.",
		)
		return nil, diags
	}

	if providerData.Logger == nil {
		diags.AddError(
			errornames.UnexpectedResourceConfigurationType,
			"Expected *zap.Logger, got nil. Please report this issue to the provider developers.",
		)
		return nil, diags
	}

	return providerData, diags
}
