package github

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/google/go-github/v84/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var vendorConfigKeys = []string{
	"azure_blob_config",
	"azure_hub_config",
	"amazon_s3_oidc_config",
	"amazon_s3_access_keys_config",
	"splunk_config",
	"hec_config",
	"google_cloud_config",
	"datadog_config",
}

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
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  vendorConfigKeys,
				Description:   "The configuration for an Azure Blob Storage audit log stream.",
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
			"azure_hub_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  vendorConfigKeys,
				Description:   "The configuration for an Azure Event Hubs audit log stream.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the Azure Event Hub.",
						},
						"encrypted_connstring": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The encrypted connection string for the Azure Event Hub.",
						},
						"key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the encryption key.",
						},
					},
				},
			},
			"amazon_s3_oidc_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  vendorConfigKeys,
				Description:   "The configuration for an Amazon S3 (OIDC) audit log stream.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the Amazon S3 bucket.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The AWS region of the S3 bucket.",
						},
						"key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the encryption key.",
						},
						"arn_role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ARN of the IAM role to assume for OIDC authentication.",
						},
					},
				},
			},
			"amazon_s3_access_keys_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  vendorConfigKeys,
				Description:   "The configuration for an Amazon S3 (access keys) audit log stream.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the Amazon S3 bucket.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The AWS region of the S3 bucket.",
						},
						"key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the encryption key.",
						},
						"encrypted_secret_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The encrypted AWS secret key.",
						},
						"encrypted_access_key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The encrypted AWS access key ID.",
						},
					},
				},
			},
			"splunk_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  vendorConfigKeys,
				Description:   "The configuration for a Splunk audit log stream.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The domain of the Splunk instance.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The port of the Splunk instance.",
						},
						"key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the encryption key.",
						},
						"encrypted_token": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The encrypted HEC token for the Splunk instance.",
						},
						"ssl_verify": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether SSL verification is enabled.",
						},
					},
				},
			},
			"hec_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  vendorConfigKeys,
				Description:   "The configuration for an HTTPS Event Collector audit log stream.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The domain of the HEC endpoint.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The port of the HEC endpoint.",
						},
						"key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the encryption key.",
						},
						"encrypted_token": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The encrypted HEC token.",
						},
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The path for the HEC endpoint.",
						},
						"ssl_verify": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether SSL verification is enabled.",
						},
					},
				},
			},
			"google_cloud_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  vendorConfigKeys,
				Description:   "The configuration for a Google Cloud Storage audit log stream.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the Google Cloud Storage bucket.",
						},
						"key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the encryption key.",
						},
						"encrypted_json_credentials": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The encrypted JSON credentials for Google Cloud.",
						},
					},
				},
			},
			"datadog_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  vendorConfigKeys,
				Description:   "The configuration for a Datadog audit log stream.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encrypted_token": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The encrypted Datadog API token.",
						},
						"site": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Datadog site (e.g. US, US3, US5, EU1, US1-FED, AP1).",
						},
						"key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the encryption key.",
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
	m := v.([]any)[0].(map[string]any)
	cfg := &github.AzureBlobConfig{
		KeyID:           m["key_id"].(string),
		EncryptedSASURL: m["encrypted_sas_url"].(string),
		Container:       m["container"].(string),
	}
	return github.NewAzureBlobStreamConfig(enabled, cfg)
}

func expandAzureHubConfig(d *schema.ResourceData, enabled bool) *github.AuditLogStreamConfig {
	v, ok := d.GetOk("azure_hub_config")
	if !ok {
		return nil
	}
	m := v.([]any)[0].(map[string]any)
	cfg := &github.AzureHubConfig{
		Name:                m["name"].(string),
		EncryptedConnstring: m["encrypted_connstring"].(string),
		KeyID:               m["key_id"].(string),
	}
	return github.NewAzureHubStreamConfig(enabled, cfg)
}

func expandAmazonS3OIDCConfig(d *schema.ResourceData, enabled bool) *github.AuditLogStreamConfig {
	v, ok := d.GetOk("amazon_s3_oidc_config")
	if !ok {
		return nil
	}
	m := v.([]any)[0].(map[string]any)
	cfg := &github.AmazonS3OIDCConfig{
		Bucket:  m["bucket"].(string),
		Region:  m["region"].(string),
		KeyID:   m["key_id"].(string),
		ArnRole: m["arn_role"].(string),
	}
	return github.NewAmazonS3OIDCStreamConfig(enabled, cfg)
}

func expandAmazonS3AccessKeysConfig(d *schema.ResourceData, enabled bool) *github.AuditLogStreamConfig {
	v, ok := d.GetOk("amazon_s3_access_keys_config")
	if !ok {
		return nil
	}
	m := v.([]any)[0].(map[string]any)
	cfg := &github.AmazonS3AccessKeysConfig{
		Bucket:               m["bucket"].(string),
		Region:               m["region"].(string),
		KeyID:                m["key_id"].(string),
		EncryptedSecretKey:   m["encrypted_secret_key"].(string),
		EncryptedAccessKeyID: m["encrypted_access_key_id"].(string),
	}
	return github.NewAmazonS3AccessKeysStreamConfig(enabled, cfg)
}

func expandSplunkConfig(d *schema.ResourceData, enabled bool) *github.AuditLogStreamConfig {
	v, ok := d.GetOk("splunk_config")
	if !ok {
		return nil
	}
	m := v.([]any)[0].(map[string]any)
	cfg := &github.SplunkConfig{
		Domain:         m["domain"].(string),
		Port:           uint16(m["port"].(int)),
		KeyID:          m["key_id"].(string),
		EncryptedToken: m["encrypted_token"].(string),
		SSLVerify:      m["ssl_verify"].(bool),
	}
	return github.NewSplunkStreamConfig(enabled, cfg)
}

func expandHecConfig(d *schema.ResourceData, enabled bool) *github.AuditLogStreamConfig {
	v, ok := d.GetOk("hec_config")
	if !ok {
		return nil
	}
	m := v.([]any)[0].(map[string]any)
	cfg := &github.HecConfig{
		Domain:         m["domain"].(string),
		Port:           uint16(m["port"].(int)),
		KeyID:          m["key_id"].(string),
		EncryptedToken: m["encrypted_token"].(string),
		Path:           m["path"].(string),
		SSLVerify:      m["ssl_verify"].(bool),
	}
	return github.NewHecStreamConfig(enabled, cfg)
}

func expandGoogleCloudConfig(d *schema.ResourceData, enabled bool) *github.AuditLogStreamConfig {
	v, ok := d.GetOk("google_cloud_config")
	if !ok {
		return nil
	}
	m := v.([]any)[0].(map[string]any)
	cfg := &github.GoogleCloudConfig{
		Bucket:                   m["bucket"].(string),
		KeyID:                    m["key_id"].(string),
		EncryptedJSONCredentials: m["encrypted_json_credentials"].(string),
	}
	return github.NewGoogleCloudStreamConfig(enabled, cfg)
}

func expandDatadogConfig(d *schema.ResourceData, enabled bool) *github.AuditLogStreamConfig {
	v, ok := d.GetOk("datadog_config")
	if !ok {
		return nil
	}
	m := v.([]any)[0].(map[string]any)
	cfg := &github.DatadogConfig{
		EncryptedToken: m["encrypted_token"].(string),
		Site:           m["site"].(string),
		KeyID:          m["key_id"].(string),
	}
	return github.NewDatadogStreamConfig(enabled, cfg)
}

// expandVendorConfig tries each vendor-specific expand function and returns the
// first non-nil result. Returns nil if no vendor config block is set.
func expandVendorConfig(d *schema.ResourceData, enabled bool) *github.AuditLogStreamConfig {
	if c := expandAzureBlobConfig(d, enabled); c != nil {
		return c
	}
	if c := expandAzureHubConfig(d, enabled); c != nil {
		return c
	}
	if c := expandAmazonS3OIDCConfig(d, enabled); c != nil {
		return c
	}
	if c := expandAmazonS3AccessKeysConfig(d, enabled); c != nil {
		return c
	}
	if c := expandSplunkConfig(d, enabled); c != nil {
		return c
	}
	if c := expandHecConfig(d, enabled); c != nil {
		return c
	}
	if c := expandGoogleCloudConfig(d, enabled); c != nil {
		return c
	}
	if c := expandDatadogConfig(d, enabled); c != nil {
		return c
	}
	return nil
}

func resourceGithubEnterpriseAuditLogStreamCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3clientV84
	enterprise := d.Get("enterprise").(string)
	enabled := d.Get("enabled").(bool)

	config := expandVendorConfig(d, enabled)
	if config == nil {
		return diag.Errorf("one of %v must be specified", vendorConfigKeys)
	}

	stream, _, err := client.Enterprise.CreateAuditLogStream(ctx, enterprise, *config)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(enterprise, strconv.FormatInt(stream.ID, 10)))

	return resourceGithubEnterpriseAuditLogStreamRead(ctx, d, meta)
}

func resourceGithubEnterpriseAuditLogStreamRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3clientV84

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
	if err := d.Set("enabled", stream.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("stream_id", stream.ID); err != nil {
		return diag.FromErr(err)
	}

	// The GitHub API does not return vendor-specific config details (encrypted
	// fields, key_id, container, etc.) in the GetAuditLogStream response — it
	// only returns a summary string in StreamDetails. The vendor config block
	// is therefore preserved from prior state automatically by Terraform.
	return nil
}

func resourceGithubEnterpriseAuditLogStreamUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3clientV84
	enterprise, streamID, err := parseAuditLogStreamID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	enabled := d.Get("enabled").(bool)

	config := expandVendorConfig(d, enabled)
	if config == nil {
		return diag.Errorf("one of %v must be specified", vendorConfigKeys)
	}

	_, _, err = client.Enterprise.UpdateAuditLogStream(ctx, enterprise, streamID, *config)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGithubEnterpriseAuditLogStreamRead(ctx, d, meta)
}

func resourceGithubEnterpriseAuditLogStreamDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3clientV84
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
