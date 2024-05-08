// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package acc

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func LoadResources(resources map[string]*terraform.ResourceState) (jsonified map[string]interface{}, err error) {
	// Load resource data from state as JSON
	stringified, err := json.Marshal(resources)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(stringified, &jsonified); err != nil {
		return nil, err
	}

	return jsonified, nil
}
