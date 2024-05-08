// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package acc

import (
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/versions"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// Provider name for single configuration testing.
	ProviderName = "mongodb"
)

// TestAccProtoV6ProviderFactories is used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	ProviderName: providerserver.NewProtocol6WithError(provider.New(versions.Test)()),
}

// TestAccProtoV6ProviderFactoriesWithProviderConfig is used
// to instantiate a provider during acceptance testing with provider config.
// The factory function will be invoked for every Terraform CLI command
// executed to create a provider server to which the CLI can reattach.
func TestAccProtoV6ProviderFactoriesWithProviderConfig(config *provider.Config) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		ProviderName: providerserver.NewProtocol6WithError(provider.WithConfig(versions.Test, config)()),
	}
}

func ProviderConfig(uri string) string {
	return fmt.Sprintf(`
		terraform {
			required_providers {
				mongodb = {
					source = "01Joseph-Hwang10/mongodb"
				}
			}
		}

		provider "mongodb" {
			uri = "%s"
		}	
	`, uri)
}

func WithProviderConfig(config string, uri string) string {
	return ProviderConfig(uri) + "\n\n" + config
}
