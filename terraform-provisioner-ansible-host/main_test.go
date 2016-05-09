package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/terraform"
)

func TestOutputHostFile(t *testing.T) {
	expectedOutput := `[testGroup]
example.test.com hostname=tester distro=debian`
	rawConfig, _ := config.NewRawConfig(map[string]interface{}{
		"inventory_path": "test_inventory",
		"host":           "example.test.com",
		"group":          "testGroup",
		"hostname":       "tester",
		"distro":         "debian",
	})
	config := terraform.NewResourceConfig(rawConfig)
	mockUIOutput := new(terraform.MockUIOutput)
	provisioner := new(ResourceProvisioner)

	err := provisioner.Apply(mockUIOutput, nil, config)

	if err != nil {
		t.Errorf("Exected error but got none")
	}
	b, fileReadErr := ioutil.ReadFile("test_inventory")
	if fileReadErr != nil {
		t.Fatalf("Error reading inventory: %s", fileReadErr)
	}
	if string(b) != expectedOutput {
		t.Errorf("Expected file content to be \"%s\" but was \"%s\"", expectedOutput, string(b))
	}
	os.Remove("test_inventory")
}
