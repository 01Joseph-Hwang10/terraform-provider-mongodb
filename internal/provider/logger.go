// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/versions"
	"go.uber.org/zap"
)

func configureLogger(p *MongoProvider) (*zap.Logger, error) {
	if p.config.Logger != nil {
		return p.config.Logger, nil
	}

	if p.version == versions.Dev {
		return zap.NewDevelopment()
	}

	return zap.NewNop(), nil
}
