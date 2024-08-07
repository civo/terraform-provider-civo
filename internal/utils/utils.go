package utils

// ValidateNameSize is a functo check the size of a name
// func ValidateNameSize

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/civo/civogo"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// VersionInfo stores Provider's version Info
type VersionInfo struct {
	ProviderSelections map[string]string `json:"provider_selections"`
}

// ValidateName is a function to check if the name is valid
func ValidateName(v interface{}, _ string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

// ValidateCNIName is a function to check if the cni name is valid
func ValidateCNIName(v interface{}, _ string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected CNI to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("CNI cannot contain whitespace. Got %s", value))
		return warns, errs
	}

	if value != "flannel" && value != "cilium" {
		errs = append(errs, fmt.Errorf("CNI plugin provided isn't valid/supported"))
		return warns, errs
	}

	return warns, errs
}

// ValidateNameSize is a function to check the size of a name
func ValidateNameSize(v interface{}, _ string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}

	if len(value) > 63 {
		errs = append(errs, fmt.Errorf("the len of the name has to be less than 63. Got %d", len(value)))
		return warns, errs
	}

	return warns, errs
}

// ResourceCommonParseID is a function to parse the ID of a resource
func ResourceCommonParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
	}

	return parts[0], parts[1], nil
}

// CheckAPPName is a function to check if the app name is valid
func CheckAPPName(appName string, client *civogo.Client) bool {
	allAPP, err := client.ListKubernetesMarketplaceApplications()
	if err != nil {
		return false
	}

	for _, v := range allAPP {
		if strings.Contains(appName, v.Name) {
			return true
		}
	}

	return false
}

// GetCommaSeparatedAllowedKeys is used by "tfplugindocs" CLI to generate Markdown docs
func GetCommaSeparatedAllowedKeys(allowedKeys []string) string {
	res := []string{}
	for _, ak := range allowedKeys {
		res = append(res, fmt.Sprintf("`%s`", ak))
	}
	sort.Strings(res)
	return strings.Join(res, ", ")
}

// ValidateNameOnlyContainsAlphanumericCharacters validate name only contains alphanumeric characters, hyphens, underscores and dots
func ValidateNameOnlyContainsAlphanumericCharacters(v interface{}, _ cty.Path) diag.Diagnostics {
	value := v.(string)
	var diags diag.Diagnostics

	_, ok := v.(string)
	if !ok {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "wrong value",
			Detail:   "expected name to be string",
		}
		diags = append(diags, diag)
	}

	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cannot contain whitespace",
			Detail:   fmt.Sprintf("name cannot contain whitespace. Got %s", value),
		}
		diags = append(diags, diag)
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9-_.]+$`).Match([]byte(value)) {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "alphanumeric characters",
			Detail:   fmt.Sprintf("name can only contain alphanumeric characters, hyphens, underscores and dots. Got %s", value),
		}
		diags = append(diags, diag)
	}

	return diags
}

// StringToInt converts a string to an int
func StringToInt(s string) (int, error) {
	s = strings.Replace(s, "G", "", 1)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// InPool is a utility function to check if a node pool is in a kubernetes cluster
func InPool(id string, list []civogo.KubernetesClusterPoolConfig) bool {
	for _, b := range list {
		if b.ID == id {
			return true
		}
	}
	return false
}

// ValidateClusterType Validates if the user has provided a supported cluster type.
func ValidateClusterType(v interface{}, path cty.Path) diag.Diagnostics {
	val := v.(string)
	var diags diag.Diagnostics
	if val != "k3s" && val != "talos" {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Cluster Type",
			Detail:   "The specified cluster type is invalid. Please choose either 'k3s' or 'talos'.",
		})
	}
	return diags
}

// ValidateProviderVersion function compares the current provider verson of the user with the threshold version and shows warning accordingly
func ValidateProviderVersion(v interface{}, path cty.Path) diag.Diagnostics {
	var versionInfo VersionInfo
	diags := diag.Diagnostics{}

	cmd := exec.Command("terraform", "version", "-json")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("[ERROR] error running terraform show: %v\n", err)
		return diags
	}

	err = json.Unmarshal(output, &versionInfo)
	if err != nil {
		log.Printf("[ERROR] error parsing JSON: %v\n", err)
		return diags
	}
	versionField := "registry.terraform.io/civo/civo"
	currentProviderVersion := versionInfo.ProviderSelections[versionField]
	thresholdProviderVersion := "1.0.49"

	v1, err := version.NewSemver(currentProviderVersion)
	if err != nil {
		log.Println("[ERROR] error parsing the given version")
		return diags
	}
	v2, err := version.NewVersion(thresholdProviderVersion)
	if err != nil {
		log.Println("[ERROR] error parsing the given version")
		return diags
	}

	lastStep := path[len(path)-1]
	var field string
	if step, ok := lastStep.(cty.GetAttrStep); ok {
		field = step.Name
	}

	if v1.LessThanOrEqual(v2) {
		if field == "write_password" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Default initial_password behavior changed",
				Detail:   "Starting from version 1.0.50 the initial password is not written to state by default, if you wish to keep the initial password configuration in state, please add the input write_password and set it to true. Example configuration: `write_password = true`.",
			})

		} else if field == "write_kubeconfig" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Default kubeconfig behavior changed",
				Detail:   "Starting from version 1.0.50, kubeconfig will no longer be written to the Terraform state by default for the civo_kubernetes resource. This change is made to enhance security by preventing sensitive information from being stored in state files. If you want to retain kubeconfig in your state file, please update your configuration by adding the `write_kubeconfig` parameter and setting it to `true`. Example configuration: `write_kubeconfig = true`.",
			})
		}
	}
	return diags
}

// CustomError captures a specific portion of the full API error
type CustomError struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("%s - %s", e.Code, e.Reason)
}

var jsonRegex *regexp.Regexp
var once sync.Once
var regexErr error

func getJSONRegex() (*regexp.Regexp, error) {

	once.Do(func() {
		jsonRegex, regexErr = regexp.Compile(`\{.*\}`)
	})
	return jsonRegex, regexErr

}

// extractJSON uses regex to find JSON content within a string
func extractJSON(s string) (string, error) {
	re, err := getJSONRegex()
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %v", err)
	}
	match := re.FindString(s)
	if match == "" {
		return "", fmt.Errorf("no JSON object found in the string")
	}
	return match, nil
}

// ParseErrorResponse extracts and parses the JSON error response
func ParseErrorResponse(errorMsg string) (*CustomError, error) {
	jsonStr, err := extractJSON(errorMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to extract JSON: %v", err)
	}

	var customErr CustomError
	err = json.Unmarshal([]byte(jsonStr), &customErr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse error response: %v", err)
	}
	return &customErr, nil
}

// ValidateUUID checks if a given string is a UUID or not
func ValidateUUID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := uuid.Parse(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q must be a valid UUID", k))
	}
	return
}
