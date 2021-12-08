package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJenkinsNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJenkinsNodeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The unique name of the JenkinsCI node.",
				Required:         true,
				ValidateDiagFunc: validateNodeName,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the JenkinsCI node.",
				Optional:    true,
			},
			"labels": {
				Type:        schema.TypeList,
				Description: "The labels the job.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip": {
				Type:             schema.TypeString,
				Description:      "The IP of the JenkinsCI node to connect to",
				Required:         true,
				ValidateDiagFunc: validateNodeIP,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "The port of the JenkinsCI node to connect to",
				Required:    true,
			},
			"credential_name": {
				Type:        schema.TypeString,
				Description: "The name of the credential, used to connect to the node",
				Required:    true,
			},
			"remote_root_directory": {
				Type:        schema.TypeString,
				Description: "The remote root directory of the JenkinsCI node",
				Required:    true,
			},
		},
	}
}

func dataSourceJenkinsNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	d.SetId(name)

	return resourceJenkinsJobRead(ctx, d, meta)
}
