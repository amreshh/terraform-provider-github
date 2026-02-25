package github

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGithubEnterpriseAuditLogStreamKey(t *testing.T) {
	enterpriseSlug := testAccConf.enterpriseSlug
	if enterpriseSlug == "" {
		t.Skip("Skipping because GITHUB_ENTERPRISE_SLUG is not set")
	}

	resourceName := "data.github_enterprise_audit_log_stream_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { skipUnlessMode(t, enterprise) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubEnterpriseAuditLogStreamKeyConfig(enterpriseSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enterprise", enterpriseSlug),
					resource.TestCheckResourceAttrSet(resourceName, "key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "key"),
				),
			},
		},
	})
}

func testAccGithubEnterpriseAuditLogStreamKeyConfig(enterpriseSlug string) string {
	return fmt.Sprintf(`
		data "github_enterprise_audit_log_stream_key" "test" {
			enterprise = "%s"
		}
	`, enterpriseSlug)
}
