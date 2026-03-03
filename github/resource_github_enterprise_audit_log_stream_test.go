package github

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGithubEnterpriseAuditLogStream_azureBlob(t *testing.T) {
	t.Parallel()

	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	container := os.Getenv("GITHUB_AZURE_BLOB_CONTAINER")
	keyID := os.Getenv("GITHUB_AZURE_BLOB_KEY_ID")
	encryptedSasURL := os.Getenv("GITHUB_AZURE_BLOB_SAS_URL")
	if container == "" || keyID == "" || encryptedSasURL == "" {
		t.Skip("Skipping because one or more Azure Blob env vars are not set " +
			"(GITHUB_AZURE_BLOB_CONTAINER, GITHUB_AZURE_BLOB_KEY_ID, GITHUB_AZURE_BLOB_SAS_URL)")
	}

	resourceName := "github_enterprise_audit_log_stream.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGithubEnterpriseAuditLogStreamDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with enabled = true
			{
				Config: testAccGithubEnterpriseAuditLogStreamAzureBlobConfig(enterpriseSlug, container, keyID, encryptedSasURL, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "azure_blob_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_blob_config.0.container", container),
					resource.TestCheckResourceAttr(resourceName, "azure_blob_config.0.key_id", keyID),
					resource.TestCheckResourceAttr(resourceName, "azure_blob_config.0.encrypted_sas_url", encryptedSasURL),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
				),
			},
			// Step 2: Import and verify
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"azure_blob_config.#",
					"azure_blob_config.0.%",
					"azure_blob_config.0.key_id",
					"azure_blob_config.0.encrypted_sas_url",
					"azure_blob_config.0.container",
				},
			},
			// Step 3: Update — disable the stream
			{
				Config: testAccGithubEnterpriseAuditLogStreamAzureBlobConfig(enterpriseSlug, container, keyID, encryptedSasURL, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "azure_blob_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_blob_config.0.container", container),
					resource.TestCheckResourceAttr(resourceName, "azure_blob_config.0.key_id", keyID),
				),
			},
			// Step 4: Update — re-enable the stream
			{
				Config: testAccGithubEnterpriseAuditLogStreamAzureBlobConfig(enterpriseSlug, container, keyID, encryptedSasURL, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccGithubEnterpriseAuditLogStream_enabledDefault(t *testing.T) {
	t.Parallel()

	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	container := os.Getenv("GITHUB_AZURE_BLOB_CONTAINER")
	keyID := os.Getenv("GITHUB_AZURE_BLOB_KEY_ID")
	encryptedSasURL := os.Getenv("GITHUB_AZURE_BLOB_SAS_URL")
	if container == "" || keyID == "" || encryptedSasURL == "" {
		t.Skip("Skipping because one or more Azure Blob env vars are not set " +
			"(GITHUB_AZURE_BLOB_CONTAINER, GITHUB_AZURE_BLOB_KEY_ID, GITHUB_AZURE_BLOB_SAS_URL)")
	}

	resourceName := "github_enterprise_audit_log_stream.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGithubEnterpriseAuditLogStreamDestroy,
		Steps: []resource.TestStep{
			// Verify that omitting enabled defaults to true
			{
				Config: testAccGithubEnterpriseAuditLogStreamAzureBlobConfigNoEnabled(enterpriseSlug, container, keyID, encryptedSasURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "azure_blob_config.#", "1"),
				),
			},
		},
	})
}

func TestAccGithubEnterpriseAuditLogStream_azureHub(t *testing.T) {
	t.Parallel()

	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	name := os.Getenv("GITHUB_AZURE_HUB_NAME")
	keyID := os.Getenv("GITHUB_AZURE_HUB_KEY_ID")
	encryptedConnstring := os.Getenv("GITHUB_AZURE_HUB_ENCRYPTED_CONNSTRING")
	if name == "" || keyID == "" || encryptedConnstring == "" {
		t.Skip("Skipping because one or more Azure Hub env vars are not set " +
			"(GITHUB_AZURE_HUB_NAME, GITHUB_AZURE_HUB_KEY_ID, GITHUB_AZURE_HUB_ENCRYPTED_CONNSTRING)")
	}

	resourceName := "github_enterprise_audit_log_stream.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGithubEnterpriseAuditLogStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubEnterpriseAuditLogStreamAzureHubConfig(enterpriseSlug, name, keyID, encryptedConnstring, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "azure_hub_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "azure_hub_config.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "azure_hub_config.0.key_id", keyID),
					resource.TestCheckResourceAttr(resourceName, "azure_hub_config.0.encrypted_connstring", encryptedConnstring),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"azure_hub_config.#",
					"azure_hub_config.0.%",
					"azure_hub_config.0.name",
					"azure_hub_config.0.key_id",
					"azure_hub_config.0.encrypted_connstring",
				},
			},
		},
	})
}

func TestAccGithubEnterpriseAuditLogStream_amazonS3OIDC(t *testing.T) {
	t.Parallel()

	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	bucket := os.Getenv("GITHUB_S3_OIDC_BUCKET")
	region := os.Getenv("GITHUB_S3_OIDC_REGION")
	keyID := os.Getenv("GITHUB_S3_OIDC_KEY_ID")
	arnRole := os.Getenv("GITHUB_S3_OIDC_ARN_ROLE")
	if bucket == "" || region == "" || keyID == "" || arnRole == "" {
		t.Skip("Skipping because one or more S3 OIDC env vars are not set " +
			"(GITHUB_S3_OIDC_BUCKET, GITHUB_S3_OIDC_REGION, GITHUB_S3_OIDC_KEY_ID, GITHUB_S3_OIDC_ARN_ROLE)")
	}

	resourceName := "github_enterprise_audit_log_stream.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGithubEnterpriseAuditLogStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubEnterpriseAuditLogStreamAmazonS3OIDCConfig(enterpriseSlug, bucket, region, keyID, arnRole, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_oidc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_oidc_config.0.bucket", bucket),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_oidc_config.0.region", region),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_oidc_config.0.key_id", keyID),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_oidc_config.0.arn_role", arnRole),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"amazon_s3_oidc_config.#",
					"amazon_s3_oidc_config.0.%",
					"amazon_s3_oidc_config.0.bucket",
					"amazon_s3_oidc_config.0.region",
					"amazon_s3_oidc_config.0.key_id",
					"amazon_s3_oidc_config.0.arn_role",
				},
			},
		},
	})
}

func TestAccGithubEnterpriseAuditLogStream_amazonS3AccessKeys(t *testing.T) {
	t.Parallel()

	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	bucket := os.Getenv("GITHUB_S3_ACCESS_KEYS_BUCKET")
	region := os.Getenv("GITHUB_S3_ACCESS_KEYS_REGION")
	keyID := os.Getenv("GITHUB_S3_ACCESS_KEYS_KEY_ID")
	encryptedSecretKey := os.Getenv("GITHUB_S3_ACCESS_KEYS_ENCRYPTED_SECRET_KEY")
	encryptedAccessKeyID := os.Getenv("GITHUB_S3_ACCESS_KEYS_ENCRYPTED_ACCESS_KEY_ID")
	if bucket == "" || region == "" || keyID == "" || encryptedSecretKey == "" || encryptedAccessKeyID == "" {
		t.Skip("Skipping because one or more S3 Access Keys env vars are not set " +
			"(GITHUB_S3_ACCESS_KEYS_BUCKET, GITHUB_S3_ACCESS_KEYS_REGION, GITHUB_S3_ACCESS_KEYS_KEY_ID, " +
			"GITHUB_S3_ACCESS_KEYS_ENCRYPTED_SECRET_KEY, GITHUB_S3_ACCESS_KEYS_ENCRYPTED_ACCESS_KEY_ID)")
	}

	resourceName := "github_enterprise_audit_log_stream.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGithubEnterpriseAuditLogStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubEnterpriseAuditLogStreamAmazonS3AccessKeysConfig(
					enterpriseSlug, bucket, region, keyID, encryptedSecretKey, encryptedAccessKeyID, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_access_keys_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_access_keys_config.0.bucket", bucket),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_access_keys_config.0.region", region),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_access_keys_config.0.key_id", keyID),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_access_keys_config.0.encrypted_secret_key", encryptedSecretKey),
					resource.TestCheckResourceAttr(resourceName, "amazon_s3_access_keys_config.0.encrypted_access_key_id", encryptedAccessKeyID),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"amazon_s3_access_keys_config.#",
					"amazon_s3_access_keys_config.0.%",
					"amazon_s3_access_keys_config.0.bucket",
					"amazon_s3_access_keys_config.0.region",
					"amazon_s3_access_keys_config.0.key_id",
					"amazon_s3_access_keys_config.0.encrypted_secret_key",
					"amazon_s3_access_keys_config.0.encrypted_access_key_id",
				},
			},
		},
	})
}

func TestAccGithubEnterpriseAuditLogStream_splunk(t *testing.T) {
	t.Parallel()

	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	domain := os.Getenv("GITHUB_SPLUNK_DOMAIN")
	port := os.Getenv("GITHUB_SPLUNK_PORT")
	keyID := os.Getenv("GITHUB_SPLUNK_KEY_ID")
	encryptedToken := os.Getenv("GITHUB_SPLUNK_ENCRYPTED_TOKEN")
	sslVerify := os.Getenv("GITHUB_SPLUNK_SSL_VERIFY")
	if domain == "" || port == "" || keyID == "" || encryptedToken == "" || sslVerify == "" {
		t.Skip("Skipping because one or more Splunk env vars are not set " +
			"(GITHUB_SPLUNK_DOMAIN, GITHUB_SPLUNK_PORT, GITHUB_SPLUNK_KEY_ID, " +
			"GITHUB_SPLUNK_ENCRYPTED_TOKEN, GITHUB_SPLUNK_SSL_VERIFY)")
	}

	resourceName := "github_enterprise_audit_log_stream.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGithubEnterpriseAuditLogStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubEnterpriseAuditLogStreamSplunkConfig(
					enterpriseSlug, domain, port, keyID, encryptedToken, sslVerify, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "splunk_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "splunk_config.0.domain", domain),
					resource.TestCheckResourceAttr(resourceName, "splunk_config.0.port", port),
					resource.TestCheckResourceAttr(resourceName, "splunk_config.0.key_id", keyID),
					resource.TestCheckResourceAttr(resourceName, "splunk_config.0.encrypted_token", encryptedToken),
					resource.TestCheckResourceAttr(resourceName, "splunk_config.0.ssl_verify", sslVerify),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"splunk_config.#",
					"splunk_config.0.%",
					"splunk_config.0.domain",
					"splunk_config.0.port",
					"splunk_config.0.key_id",
					"splunk_config.0.encrypted_token",
					"splunk_config.0.ssl_verify",
				},
			},
		},
	})
}

func TestAccGithubEnterpriseAuditLogStream_hec(t *testing.T) {
	t.Parallel()

	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	domain := os.Getenv("GITHUB_HEC_DOMAIN")
	port := os.Getenv("GITHUB_HEC_PORT")
	keyID := os.Getenv("GITHUB_HEC_KEY_ID")
	encryptedToken := os.Getenv("GITHUB_HEC_ENCRYPTED_TOKEN")
	path := os.Getenv("GITHUB_HEC_PATH")
	sslVerify := os.Getenv("GITHUB_HEC_SSL_VERIFY")
	if domain == "" || port == "" || keyID == "" || encryptedToken == "" || path == "" || sslVerify == "" {
		t.Skip("Skipping because one or more HEC env vars are not set " +
			"(GITHUB_HEC_DOMAIN, GITHUB_HEC_PORT, GITHUB_HEC_KEY_ID, " +
			"GITHUB_HEC_ENCRYPTED_TOKEN, GITHUB_HEC_PATH, GITHUB_HEC_SSL_VERIFY)")
	}

	resourceName := "github_enterprise_audit_log_stream.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGithubEnterpriseAuditLogStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubEnterpriseAuditLogStreamHecConfig(
					enterpriseSlug, domain, port, keyID, encryptedToken, path, sslVerify, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "hec_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "hec_config.0.domain", domain),
					resource.TestCheckResourceAttr(resourceName, "hec_config.0.port", port),
					resource.TestCheckResourceAttr(resourceName, "hec_config.0.key_id", keyID),
					resource.TestCheckResourceAttr(resourceName, "hec_config.0.encrypted_token", encryptedToken),
					resource.TestCheckResourceAttr(resourceName, "hec_config.0.path", path),
					resource.TestCheckResourceAttr(resourceName, "hec_config.0.ssl_verify", sslVerify),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"hec_config.#",
					"hec_config.0.%",
					"hec_config.0.domain",
					"hec_config.0.port",
					"hec_config.0.key_id",
					"hec_config.0.encrypted_token",
					"hec_config.0.path",
					"hec_config.0.ssl_verify",
				},
			},
		},
	})
}

func TestAccGithubEnterpriseAuditLogStream_googleCloud(t *testing.T) {
	t.Parallel()

	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	bucket := os.Getenv("GITHUB_GOOGLE_CLOUD_BUCKET")
	keyID := os.Getenv("GITHUB_GOOGLE_CLOUD_KEY_ID")
	encryptedJSONCreds := os.Getenv("GITHUB_GOOGLE_CLOUD_ENCRYPTED_JSON_CREDENTIALS")
	if bucket == "" || keyID == "" || encryptedJSONCreds == "" {
		t.Skip("Skipping because one or more Google Cloud env vars are not set " +
			"(GITHUB_GOOGLE_CLOUD_BUCKET, GITHUB_GOOGLE_CLOUD_KEY_ID, GITHUB_GOOGLE_CLOUD_ENCRYPTED_JSON_CREDENTIALS)")
	}

	resourceName := "github_enterprise_audit_log_stream.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGithubEnterpriseAuditLogStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubEnterpriseAuditLogStreamGoogleCloudConfig(
					enterpriseSlug, bucket, keyID, encryptedJSONCreds, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_config.0.bucket", bucket),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_config.0.key_id", keyID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_config.0.encrypted_json_credentials", encryptedJSONCreds),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"google_cloud_config.#",
					"google_cloud_config.0.%",
					"google_cloud_config.0.bucket",
					"google_cloud_config.0.key_id",
					"google_cloud_config.0.encrypted_json_credentials",
				},
			},
		},
	})
}

func TestAccGithubEnterpriseAuditLogStream_datadog(t *testing.T) {
	t.Parallel()

	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	encryptedToken := os.Getenv("GITHUB_DATADOG_ENCRYPTED_TOKEN")
	site := os.Getenv("GITHUB_DATADOG_SITE")
	keyID := os.Getenv("GITHUB_DATADOG_KEY_ID")
	if encryptedToken == "" || site == "" || keyID == "" {
		t.Skip("Skipping because one or more Datadog env vars are not set " +
			"(GITHUB_DATADOG_ENCRYPTED_TOKEN, GITHUB_DATADOG_SITE, GITHUB_DATADOG_KEY_ID)")
	}

	resourceName := "github_enterprise_audit_log_stream.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGithubEnterpriseAuditLogStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubEnterpriseAuditLogStreamDatadogConfig(
					enterpriseSlug, encryptedToken, site, keyID, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "datadog_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "datadog_config.0.encrypted_token", encryptedToken),
					resource.TestCheckResourceAttr(resourceName, "datadog_config.0.site", site),
					resource.TestCheckResourceAttr(resourceName, "datadog_config.0.key_id", keyID),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"datadog_config.#",
					"datadog_config.0.%",
					"datadog_config.0.encrypted_token",
					"datadog_config.0.site",
					"datadog_config.0.key_id",
				},
			},
		},
	})
}

func testAccCheckGithubEnterpriseAuditLogStreamDestroy(s *terraform.State) error {
	meta, err := getTestMeta()
	if err != nil {
		return err
	}
	conn := meta.v3clientV84

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "github_enterprise_audit_log_stream" {
			continue
		}

		enterprise, streamID, err := parseAuditLogStreamID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, _, err = conn.Enterprise.GetAuditLogStream(context.Background(), enterprise, streamID)
		if err == nil {
			return fmt.Errorf("audit log stream %d still exists in enterprise %s", streamID, enterprise)
		}
	}

	return nil
}

// --- Config helpers ---

func testAccGithubEnterpriseAuditLogStreamAzureBlobConfig(enterpriseSlug, container, keyID, encryptedSasURL string, enabled bool) string {
	return fmt.Sprintf(`
resource "github_enterprise_audit_log_stream" "test" {
  enterprise = "%s"
  enabled    = %s

  azure_blob_config {
    container         = "%s"
    key_id            = "%s"
    encrypted_sas_url = "%s"
  }
}
`, enterpriseSlug, strconv.FormatBool(enabled), container, keyID, encryptedSasURL)
}

func testAccGithubEnterpriseAuditLogStreamAzureBlobConfigNoEnabled(enterpriseSlug, container, keyID, encryptedSasURL string) string {
	return fmt.Sprintf(`
resource "github_enterprise_audit_log_stream" "test" {
  enterprise = "%s"

  azure_blob_config {
    container         = "%s"
    key_id            = "%s"
    encrypted_sas_url = "%s"
  }
}
`, enterpriseSlug, container, keyID, encryptedSasURL)
}

func testAccGithubEnterpriseAuditLogStreamAzureHubConfig(enterpriseSlug, name, keyID, encryptedConnstring string, enabled bool) string {
	return fmt.Sprintf(`
resource "github_enterprise_audit_log_stream" "test" {
  enterprise = "%s"
  enabled    = %s

  azure_hub_config {
    name                 = "%s"
    key_id               = "%s"
    encrypted_connstring = "%s"
  }
}
`, enterpriseSlug, strconv.FormatBool(enabled), name, keyID, encryptedConnstring)
}

func testAccGithubEnterpriseAuditLogStreamAmazonS3OIDCConfig(enterpriseSlug, bucket, region, keyID, arnRole string, enabled bool) string {
	return fmt.Sprintf(`
resource "github_enterprise_audit_log_stream" "test" {
  enterprise = "%s"
  enabled    = %s

  amazon_s3_oidc_config {
    bucket   = "%s"
    region   = "%s"
    key_id   = "%s"
    arn_role = "%s"
  }
}
`, enterpriseSlug, strconv.FormatBool(enabled), bucket, region, keyID, arnRole)
}

func testAccGithubEnterpriseAuditLogStreamAmazonS3AccessKeysConfig(enterpriseSlug, bucket, region, keyID, encryptedSecretKey, encryptedAccessKeyID string, enabled bool) string {
	return fmt.Sprintf(`
resource "github_enterprise_audit_log_stream" "test" {
  enterprise = "%s"
  enabled    = %s

  amazon_s3_access_keys_config {
    bucket                  = "%s"
    region                  = "%s"
    key_id                  = "%s"
    encrypted_secret_key    = "%s"
    encrypted_access_key_id = "%s"
  }
}
`, enterpriseSlug, strconv.FormatBool(enabled), bucket, region, keyID, encryptedSecretKey, encryptedAccessKeyID)
}

func testAccGithubEnterpriseAuditLogStreamSplunkConfig(enterpriseSlug, domain, port, keyID, encryptedToken, sslVerify string, enabled bool) string {
	return fmt.Sprintf(`
resource "github_enterprise_audit_log_stream" "test" {
  enterprise = "%s"
  enabled    = %s

  splunk_config {
    domain          = "%s"
    port            = %s
    key_id          = "%s"
    encrypted_token = "%s"
    ssl_verify      = %s
  }
}
`, enterpriseSlug, strconv.FormatBool(enabled), domain, port, keyID, encryptedToken, sslVerify)
}

func testAccGithubEnterpriseAuditLogStreamHecConfig(enterpriseSlug, domain, port, keyID, encryptedToken, path, sslVerify string, enabled bool) string {
	return fmt.Sprintf(`
resource "github_enterprise_audit_log_stream" "test" {
  enterprise = "%s"
  enabled    = %s

  hec_config {
    domain          = "%s"
    port            = %s
    key_id          = "%s"
    encrypted_token = "%s"
    path            = "%s"
    ssl_verify      = %s
  }
}
`, enterpriseSlug, strconv.FormatBool(enabled), domain, port, keyID, encryptedToken, path, sslVerify)
}

func testAccGithubEnterpriseAuditLogStreamGoogleCloudConfig(enterpriseSlug, bucket, keyID, encryptedJSONCredentials string, enabled bool) string {
	return fmt.Sprintf(`
resource "github_enterprise_audit_log_stream" "test" {
  enterprise = "%s"
  enabled    = %s

  google_cloud_config {
    bucket                     = "%s"
    key_id                     = "%s"
    encrypted_json_credentials = "%s"
  }
}
`, enterpriseSlug, strconv.FormatBool(enabled), bucket, keyID, encryptedJSONCredentials)
}

func testAccGithubEnterpriseAuditLogStreamDatadogConfig(enterpriseSlug, encryptedToken, site, keyID string, enabled bool) string {
	return fmt.Sprintf(`
resource "github_enterprise_audit_log_stream" "test" {
  enterprise = "%s"
  enabled    = %s

  datadog_config {
    encrypted_token = "%s"
    site            = "%s"
    key_id          = "%s"
  }
}
`, enterpriseSlug, strconv.FormatBool(enabled), encryptedToken, site, keyID)
}
