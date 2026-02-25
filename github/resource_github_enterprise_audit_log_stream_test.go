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

func TestAccGithubEnterpriseAuditLogStream(t *testing.T) {
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
				Config: testAccGithubEnterpriseAuditLogStreamConfig(enterpriseSlug, container, keyID, encryptedSasURL, true),
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
			// The API does not return vendor-specific config, so we must ignore
			// the entire azure_blob_config block on import.
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
				Config: testAccGithubEnterpriseAuditLogStreamConfig(enterpriseSlug, container, keyID, encryptedSasURL, false),
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
				Config: testAccGithubEnterpriseAuditLogStreamConfig(enterpriseSlug, container, keyID, encryptedSasURL, true),
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
				Config: testAccGithubEnterpriseAuditLogStreamConfigNoEnabled(enterpriseSlug, container, keyID, encryptedSasURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "azure_blob_config.#", "1"),
				),
			},
		},
	})
}

func testAccCheckGithubEnterpriseAuditLogStreamDestroy(s *terraform.State) error {
	meta, err := getTestMeta()
	if err != nil {
		return err
	}
	conn := meta.v3client

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

func testAccGithubEnterpriseAuditLogStreamConfig(enterpriseSlug, container, keyID, encryptedSasURL string, enabled bool) string {
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

func testAccGithubEnterpriseAuditLogStreamConfigNoEnabled(enterpriseSlug, container, keyID, encryptedSasURL string) string {
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
