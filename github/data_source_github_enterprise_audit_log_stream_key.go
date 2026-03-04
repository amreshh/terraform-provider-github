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
			"slug": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the enterprise.",
			},
			"key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Audit Log Streaming Public Key ID.",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Audit Log Streaming Public Key.",
			},
		},
	}
}

func dataSourceGithubEnterpriseAuditLogStreamKeyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3clientV84
	slug := d.Get("slug").(string)

	key, _, err := client.Enterprise.GetAuditLogStreamKey(ctx, slug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(slug)
	if err := d.Set("key_id", key.KeyID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("key", key.Key); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
