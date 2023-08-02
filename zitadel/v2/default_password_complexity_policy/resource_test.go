package default_password_complexity_policy_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccDefaultPasswordComplexityPolicy(t *testing.T) {
	resourceName := "zitadel_default_password_complexity_policy"
	initialProperty := true
	updatedProperty := false
	frame, err := test_utils.NewInstanceTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		func(configProperty, _ interface{}) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  min_length    = "8"
  has_uppercase = true
  has_lowercase = true
  has_number    = true
  has_symbol    = %t
}`, resourceName, frame.UniqueResourcesID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckNothing,
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(interface{}) resource.TestCheckFunc {
	return func(expect interface{}) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetPasswordComplexityPolicy(frame, &admin.GetPasswordComplexityPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetHasSymbol()
			if actual != expect {
				return fmt.Errorf("expected %t, but got %t", expect, actual)
			}
			return nil
		}
	}
}
