// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package resourceconfig

import (
	"fmt"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
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
		diags.Append(
			errs.NewUnexpectedResourceConfigurationType(
				"*resourceconfig.ResourceConfig",
				fmt.Sprintf("%T", data),
			).ToDiagnostic(),
		)
		return nil, diags
	}

	if providerData.ClientConfig == nil {
		diags.Append(
			errs.NewUnexpectedResourceConfigurationType(
				"*mongoclient.Config",
				fmt.Sprintf("%T", nil),
			).ToDiagnostic(),
		)
		return nil, diags
	}

	if providerData.Logger == nil {
		diags.Append(
			errs.NewUnexpectedResourceConfigurationType(
				"*zap.Logger",
				fmt.Sprintf("%T", nil),
			).ToDiagnostic(),
		)
		return nil, diags
	}

	return providerData, diags
}
