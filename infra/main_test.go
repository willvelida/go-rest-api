package infra

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureSample(t *testing.T) {
	t.Parallel()

	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	// Configure Terraform options
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../infra",
	}

	// at the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run terraform init and apply. Fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	appName := terraform.Output(t, terraformOptions, "linux_web_app_name")
	appServicePlanId := terraform.Output(t, terraformOptions, "linux_web_app_asp")

	// Look up resources and ensure it matches the output
	actualResourceGroup := azure.ResourceGroupExists(t, resourceGroupName, subscriptionId)
	app := azure.GetAppService(t, appName, resourceGroupName, subscriptionId)

	assert.True(t, actualResourceGroup, "Resource group does not exist")
	assert.True(t, azure.AppExists(t, appName, resourceGroupName, subscriptionId))

	assert.Equal(t, appName, app.Name, "App name does not match")
	assert.Equal(t, appServicePlanId, app.ServerFarmID, "App service plan id does not match")

}
