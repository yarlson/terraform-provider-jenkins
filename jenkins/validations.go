package jenkins

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func validateJobName(val interface{}, path cty.Path) diag.Diagnostics {
	if strings.Contains(val.(string), "/") {
		return diag.FromErr(fmt.Errorf("provided name includes path characters. Please use the 'folder' property if specifying a job within a subfolder"))
	}

	return diag.Diagnostics{}
}

func validateFolderName(val interface{}, path cty.Path) diag.Diagnostics {
	return diag.Diagnostics{}
}

func validateCredentialScope(val interface{}, path cty.Path) diag.Diagnostics {
	var supportedCredentialScopes = []string{"SYSTEM", "GLOBAL"}
	for _, supported := range supportedCredentialScopes {
		if val == supported {
			return diag.Diagnostics{}
		}
	}
	return diag.Errorf("Invalid scope: %s. Supported scopes are: %s", val, strings.Join(supportedCredentialScopes, ", "))
}

func validateNodeName(val interface{}, path cty.Path) diag.Diagnostics {
	if strings.Contains(val.(string), "/") {
		return diag.FromErr(fmt.Errorf("provided name includes path characters"))
	}

	return diag.Diagnostics{}
}

func validateNodeIP(val interface{}, path cty.Path) diag.Diagnostics {
	if net.ParseIP(val.(string)) == nil {
		return diag.FromErr(fmt.Errorf("provided IP is not valid"))
	}

	return diag.Diagnostics{}
}
