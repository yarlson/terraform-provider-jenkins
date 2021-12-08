package jenkins

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceJenkinsNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJenkinsNodeCreate,
		ReadContext:   resourceJenkinsNodeRead,
		UpdateContext: resourceJenkinsNodeUpdate,
		DeleteContext: resourceJenkinsNodeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The unique name of the JenkinsCI node.",
				Required:         true,
				ForceNew:         true,
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
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip": {
				Type:             schema.TypeString,
				Description:      "The IP of the JenkinsCI node to connect to",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateNodeIP,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "The port of the JenkinsCI node to connect to",
				Required:    true,
				ForceNew:    true,
			},
			"credential_name": {
				Type:        schema.TypeString,
				Description: "The name of the credential, used to connect to the node",
				Required:    true,
				ForceNew:    true,
			},
			"remote_root_directory": {
				Type:        schema.TypeString,
				Description: "The remote root directory of the JenkinsCI node",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceJenkinsNodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels").([]interface{})
	ip := d.Get("ip").(string)
	port := d.Get("port").(int)
	credentials := d.Get("credential_name").(string)
	remoteFS := d.Get("remote_root_directory").(string)

	var ciLabels []string
	for _, label := range labels {
		ciLabels = append(ciLabels, label.(string))
	}
	_, err := client.RegisterNode(ctx, name, 1, description, remoteFS, credentials, ciLabels, ip, port)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	return resourceJenkinsNodeRead(ctx, d, meta)
}

func resourceJenkinsNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name := d.Id()

	log.Printf("[DEBUG] jenkins::read - Looking for node %q", name)
	_, err := client.GetNode(ctx, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::read - Could not find node %q: %w", name, err))
	}

	log.Printf("[DEBUG] jenkins::read - Node %q exists", name)

	return nil
}

func resourceJenkinsNodeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJenkinsNodeRead(ctx, d, meta)
}

func resourceJenkinsNodeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name := d.Id()

	log.Printf("[DEBUG] jenkins::delete - Removing %q", name)

	ok, err := client.GetNode(ctx, name)
	if err != nil {
		if err.Error() != "No node found" {
			return diag.FromErr(fmt.Errorf("jenkins::delete - Could not find node %q: %w", name, err))
		}
	}

	_, err = client.DeleteNode(ctx, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::delete - Could not delete node %q: %w", name, err))
	}

	log.Printf("[DEBUG] jenkins::delete - %q removed: %v", name, ok)
	return nil
}
