// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"

	testutils "github.com/hewlettpackard/hpegl-metal-terraform-resources/internal/test-utils"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = testutils.ProviderFunc()()
	testAccProviders = map[string]*schema.Provider{
		"hpegl": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
}

func TestProvider(t *testing.T) {
	if err := testutils.ProviderFunc()().InternalValidate(); err != nil {
		t.Fatalf("%s\n", err)
	}
	testAccPreCheck(t)
}

//nolint: funlen  // ignoring for test functions
func TestAccProvider_Auth0Error(t *testing.T) {
	filePath, err := getDefaultMetalConfigPath()
	if err != nil {
		t.Fatalf("Failed to get default qjwt file path: %v", err)
	}

	filePathOrig := filePath + ".orig"

	// ignoring error if file doesn't exist
	err = os.Rename(filePath, filePathOrig)
	if err == nil {
		defer func() {
			if err := os.Rename(filePathOrig, filePath); err != nil {
				t.Logf("igoring error: %v", err)
			}
		}()
	}

	glTokenEnvOrig := os.Getenv("HPEGL_METAL_GL_TOKEN")

	// turning on Metal Auth mode
	os.Setenv("HPEGL_METAL_GL_TOKEN", "false")
	defer os.Setenv("HPEGL_METAL_GL_TOKEN", glTokenEnvOrig)

	tCases := []struct {
		name string
		conf map[string]string
	}{
		{
			name: "metal token - rest_url missing",
			conf: map[string]string{
				"jwt":       "hastoken",
				"member_id": "portal_owner_id",
			},
		},
		{
			name: "metal token - jwt missing",
			conf: map[string]string{
				"rest_url":  "https://test",
				"member_id": "portal_owner_id",
			},
		},
		{
			name: "metal token - member_id missing",
			conf: map[string]string{
				"rest_url": "https://test",
				"jwt":      "havetoken",
			},
		},
	}

	for _, tc := range tCases {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				if err := writeToYAMLFile(tc.conf, filePath); err != nil {
					t.Fatalf("failed to create .qjwt file: %v", err)
				}
			},
			Providers: testAccProviders,

			Steps: []resource.TestStep{
				{
					Config:      testAccProviderMetalBasic(),
					ExpectError: regexp.MustCompile("configuration error"),
				},
			},
		})
	}
}

func TestAccProvider_IAMAuthError(t *testing.T) {
	filePath, err := getDefaultIAMConfigPath()
	if err != nil {
		t.Fatalf("Failed to get default .gltform file path: %v", err)
	}

	filePathOrig := filePath + ".orig"

	// ignoring error if file doesn't exist
	err = os.Rename(filePath, filePathOrig)
	if err == nil {
		defer func() {
			if err := os.Rename(filePathOrig, filePath); err != nil {
				t.Logf("igoring error: %v", err)
			}
		}()
	}

	glTokenEnvOrig := os.Getenv("HPEGL_METAL_GL_TOKEN")

	// turning on Metal Auth mode
	os.Setenv("HPEGL_METAL_GL_TOKEN", "true")
	defer os.Setenv("HPEGL_METAL_GL_TOKEN", glTokenEnvOrig)

	tCases := []struct {
		name string
		conf map[string]string
	}{
		{
			name: "glc token - rest_url missing",
			conf: map[string]string{
				"project_id": "team_project_id",
			},
		},
	}

	for _, tc := range tCases {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				if err := writeToYAMLFile(tc.conf, filePath); err != nil {
					t.Fatalf("failed to create .gltform file: %v", err)
				}
			},
			Providers: testAccProviders,

			Steps: []resource.TestStep{
				{
					Config:      testAccProviderMetalBasic(),
					ExpectError: regexp.MustCompile("configuration error"),
				},
			},
		})
	}
}

func testAccProviderMetalBasic() string {
	return `
provider "hpegl" {
	metal {
	}
}

data "hpegl_metal_available_resources" "physical" {

}`
}

func getDefaultMetalConfigPath() (string, error) {
	filename := ".qjwt"
	// In Metal auth mode, the lookup order is Home dir and then the Work dir.
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return filepath.Clean(filepath.Join(dir, filename)), nil
}

func getDefaultIAMConfigPath() (string, error) {
	filename := ".gltform"
	// In IAM auth mode, the lookup order is Work dir and then the Home dir.
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return filepath.Clean(filepath.Join(dir, filename)), nil
}

func writeToYAMLFile(d map[string]string, filePath string) error {
	// Marshal config
	b, err := yaml.Marshal(d)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	err = f.Sync()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
