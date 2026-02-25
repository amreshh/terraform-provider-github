package github

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGithubEnterpriseAuditLogStreamKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGithubEnterpriseAuditLogStreamKeyRead,
		Schema: map[string]*schema.Schema{
			"enterprise": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the enterprise.",
			},
			"key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the public key.",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public key.",
			},
		},
	}
}

func dataSourceGithubEnterpriseAuditLogStreamKeyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client
	enterprise := d.Get("enterprise").(string)

	key, _, err := client.Enterprise.GetAuditLogStreamKey(ctx, enterprise)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(enterprise)
	if err := d.Set("key_id", key.GetKeyID()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("key", key.GetKey()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
