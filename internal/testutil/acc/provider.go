// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package acc

import (
	"fmt"

	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// Provider name for single configuration testing
	ProviderName = "mongodb"
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	ProviderName: providerserver.NewProtocol6WithError(provider.New("test")()),
}

func providerConfig(uri string) string {
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
	return providerConfig(uri) + "\n\n" + config
}
