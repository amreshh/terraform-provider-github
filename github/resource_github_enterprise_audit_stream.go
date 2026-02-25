package github

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/google/go-github/v83/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGithubEnterpriseAuditLogStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGithubEnterpriseAuditLogStreamCreate,
		ReadContext:   resourceGithubEnterpriseAuditLogStreamRead,
		UpdateContext: resourceGithubEnterpriseAuditLogStreamUpdate,
		DeleteContext: resourceGithubEnterpriseAuditLogStreamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"enterprise": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The slug of the enterprise.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the audit log stream is enabled.",
			},
			"stream_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the audit log stream.",
			},
			"azure_blob_config": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				AtLeastOneOf: []string{"azure_blob_config"},
				Description:  "The configuration for an Azure Blob Storage audit log stream.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the public key used to encrypt the SAS URL.",
						},
						"encrypted_sas_url": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The encrypted SAS URL for the Azure Blob Storage container.",
						},
						"container": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the Azure Blob Storage container.",
						},
					},
				},
			},
		},
	}
}

// parseAuditLogStreamID extracts the enterprise slug and stream ID from the
// composite resource ID (enterprise:stream_id).
func parseAuditLogStreamID(id string) (enterprise string, streamID int64, err error) {
	enterprise, streamIDStr, err := parseTwoPartID(id, "enterprise", "stream_id")
	if err != nil {
		return "", 0, err
	}
	streamID, err = strconv.ParseInt(streamIDStr, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid stream_id %q: %w", streamIDStr, err)
	}
	return enterprise, streamID, nil
}

// expandAzureBlobConfig reads the azure_blob_config block from ResourceData and
// returns an AuditLogStreamConfig ready to send to the API. Returns nil if no
// azure_blob_config block is present.
func expandAzureBlobConfig(d *schema.ResourceData, enabled bool) *github.AuditLogStreamConfig {
	v, ok := d.GetOk("azure_blob_config")
	if !ok {
		return nil
	}
	azureBlob := v.([]any)[0].(map[string]any)
	azureConfig := &github.AzureBlobConfig{
		KeyID:           github.Ptr(azureBlob["key_id"].(string)),
		EncryptedSasURL: github.Ptr(azureBlob["encrypted_sas_url"].(string)),
		Container:       github.Ptr(azureBlob["container"].(string)),
	}
	return github.NewAzureBlobStreamConfig(enabled, azureConfig)
}

func resourceGithubEnterpriseAuditLogStreamCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client
	enterprise := d.Get("enterprise").(string)
	enabled := d.Get("enabled").(bool)

	config := expandAzureBlobConfig(d, enabled)
	if config == nil {
		return diag.Errorf("one of azure_blob_config must be specified")
	}

	stream, _, err := client.Enterprise.CreateAuditLogStream(ctx, enterprise, config)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(enterprise, strconv.FormatInt(stream.GetID(), 10)))

	return resourceGithubEnterpriseAuditLogStreamRead(ctx, d, meta)
}

func resourceGithubEnterpriseAuditLogStreamRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client

	enterprise, streamID, err := parseAuditLogStreamID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	stream, _, err := client.Enterprise.GetAuditLogStream(ctx, enterprise, streamID)
	if err != nil {
		var ghErr *github.ErrorResponse
		if errors.As(err, &ghErr) && ghErr.Response.StatusCode == 404 {
			log.Printf("[INFO] Removing audit log stream %d from state because it no longer exists in GitHub", streamID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("enterprise", enterprise); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", stream.GetEnabled()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("stream_id", stream.GetID()); err != nil {
		return diag.FromErr(err)
	}

	// The GitHub API does not return vendor-specific config details (encrypted
	// fields, key_id, container, etc.) in the GetAuditLogStream response — it
	// only returns a summary string in StreamDetails. The azure_blob_config
	// block is therefore preserved from prior state automatically by Terraform.
	return nil
}

func resourceGithubEnterpriseAuditLogStreamUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client
	enterprise, streamID, err := parseAuditLogStreamID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	enabled := d.Get("enabled").(bool)

	config := expandAzureBlobConfig(d, enabled)
	if config == nil {
		return diag.Errorf("one of azure_blob_config must be specified")
	}

	_, _, err = client.Enterprise.UpdateAuditLogStream(ctx, enterprise, streamID, config)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGithubEnterpriseAuditLogStreamRead(ctx, d, meta)
}

func resourceGithubEnterpriseAuditLogStreamDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client
	enterprise, streamID, err := parseAuditLogStreamID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Enterprise.DeleteAuditLogStream(ctx, enterprise, streamID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
