package utils

// ValidateNameSize is a functo check the size of a name
// func ValidateNameSize

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// FileSizeLimit limits the size of file to be used by the user
	FileSizeLimit = int64(20 * 1024 * 1024) // 20 MB
)

// VersionInfo stores Provider's version Info
type VersionInfo struct {
	ProviderSelections map[string]string `json:"provider_selections"`
}

// IgnoreCaseDiff is a DiffSuppressFunc that suppresses diffs caused only by
// case differences (e.g. "FRA1" vs "fra1"). This is used for the region field
// across all resources so that changing the case of a region identifier does not
// trigger an unnecessary update or replacement.
func IgnoreCaseDiff(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	return strings.EqualFold(oldValue, newValue)
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

// CheckFileSize function checks if the file the file size is less than the allowed limit(current: 20MB)
func CheckFileSize(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("error getting the file info: %w", err)
	}
	if fileInfo.Size() > FileSizeLimit {
		return fmt.Errorf("file size exceeds the allowed limit of %d bytes", FileSizeLimit)
	}
	return nil
}

// RegionRef points at a schema field holding the ID of another region-scoped
// resource, plus a probe that reports whether that ID is visible in a given
// region. It lets a resource work out which region it belongs in when the user
// did not spell out `region` themselves.
type RegionRef struct {
	// Field is the schema key holding the referenced resource's ID.
	Field string
	// Kind names the referenced resource type, for log and error messages.
	Kind string
	// Probe reports whether id exists in the region c is scoped to. False with a
	// nil error means definitely absent; an error means the lookup itself failed
	// and says nothing either way. A 200 carrying an empty record counts as
	// absent, so a lenient endpoint cannot make us settle on the wrong region.
	Probe func(c *civogo.Client, id string) (bool, error)
}

// NetworkRef builds a RegionRef for a field holding a network ID.
func NetworkRef(field string) RegionRef {
	return RegionRef{Field: field, Kind: "network", Probe: func(c *civogo.Client, id string) (bool, error) {
		n, err := c.GetVPCNetwork(id)
		if err != nil {
			return false, err
		}
		return n != nil && n.ID != "", nil
	}}
}

// FirewallRef builds a RegionRef for a field holding a firewall ID.
func FirewallRef(field string) RegionRef {
	return RegionRef{Field: field, Kind: "firewall", Probe: func(c *civogo.Client, id string) (bool, error) {
		f, err := c.FindVPCFirewall(id)
		if err != nil {
			return false, err
		}
		return f != nil && f.ID != "", nil
	}}
}

// InstanceRef builds a RegionRef for a field holding an instance ID.
func InstanceRef(field string) RegionRef {
	return RegionRef{Field: field, Kind: "instance", Probe: func(c *civogo.Client, id string) (bool, error) {
		i, err := c.GetInstance(id)
		if err != nil {
			return false, err
		}
		return i != nil && i.ID != "", nil
	}}
}

// VolumeRef builds a RegionRef for a field holding a volume ID.
func VolumeRef(field string) RegionRef {
	return RegionRef{Field: field, Kind: "volume", Probe: func(c *civogo.Client, id string) (bool, error) {
		v, err := c.GetVolume(id)
		if err != nil {
			return false, err
		}
		return v != nil && v.ID != "", nil
	}}
}

// RegionOption configures how RegionalClient picks a region.
type RegionOption func(*regionSpec)

type regionSpec struct {
	region string
	d      *schema.ResourceData
	refs   []RegionRef
}

// WithRegion pins the client to a region the caller already knows.
func WithRegion(region string) RegionOption {
	return func(s *regionSpec) { s.region = region }
}

// ResolveRegion takes the region from the resource's own "region" field, and
// when that is empty infers it from the first ref whose ID is set.
//
// An inferred region is written back into d, so the Read that follows a create,
// and every operation after it, scope themselves from state without repeating
// the lookup. Pass refs only for resources whose "region" field is Computed;
// otherwise that write-back reads as a diff against an empty configuration.
func ResolveRegion(d *schema.ResourceData, refs ...RegionRef) RegionOption {
	return func(s *regionSpec) {
		s.d = d
		s.refs = refs
	}
}

// RegionalClient returns a shallow copy of the API client scoped to a region.
//
// Resources must use this rather than mutating apiClient.Region: the provider
// shares one *civogo.Client across every resource and Terraform runs their
// operations concurrently, so mutating it races and can send requests to the
// wrong region (issue #395).
//
// The region is picked in this order:
//
//  1. WithRegion, when the caller already knows it
//  2. the resource's own "region" field
//  3. the region of the first ResolveRegion ref whose ID is set, via an API lookup
//  4. otherwise the client unchanged, leaving the provider's region or the
//     account default to apply
//
// Step 3 is what makes a firewall that only names a network_id land in that
// network's region. Step 4 is deliberate: configs that set no region anywhere
// and rely on the account default must keep working, and a resource with no
// region and no reference is unambiguous anyway.
func RegionalClient(apiClient *civogo.Client, opts ...RegionOption) (*civogo.Client, error) {
	spec := &regionSpec{}
	for _, opt := range opts {
		opt(spec)
	}

	if spec.region != "" {
		return regionScoped(apiClient, spec.region), nil
	}
	if spec.d == nil {
		return apiClient, nil
	}

	if region, ok := spec.d.GetOk("region"); ok && region.(string) != "" {
		return regionScoped(apiClient, region.(string)), nil
	}

	for _, ref := range spec.refs {
		id, ok := spec.d.GetOk(ref.Field)
		if !ok || id.(string) == "" {
			continue
		}

		region, err := findRegionOf(apiClient, ref, id.(string))
		if err != nil {
			return nil, err
		}
		log.Printf("[INFO] resolved region %q from %s %s (%q was not set)", region, ref.Kind, id, ref.Field)

		if err := spec.d.Set("region", region); err != nil {
			return nil, fmt.Errorf("could not record the resolved region %q: %w", region, err)
		}
		return regionScoped(apiClient, region), nil
	}

	return apiClient, nil
}

// findRegionOf searches the account's regions for the one holding id.
func findRegionOf(apiClient *civogo.Client, ref RegionRef, id string) (string, error) {
	regions, err := apiClient.ListRegions()
	if err != nil {
		return "", fmt.Errorf("could not list regions to find %s %s: %w", ref.Kind, id, err)
	}

	// The client's own region first when it has one: the likeliest answer, so the
	// common case costs a single extra call.
	codes := make([]string, 0, len(regions)+1)
	if apiClient.Region != "" {
		codes = append(codes, apiClient.Region)
	}
	for _, r := range regions {
		if !strings.EqualFold(r.Code, apiClient.Region) {
			codes = append(codes, r.Code)
		}
	}

	// A probe that errors says nothing about whether the resource is there, so
	// keep looking but hold on to the failure. Without it, one unhealthy region
	// would be reported as "not found anywhere", sending people after the wrong
	// problem.
	var probeErr error
	for _, code := range codes {
		found, err := ref.Probe(regionScoped(apiClient, code), id)
		if err != nil {
			probeErr = fmt.Errorf("looking in %s: %w", code, err)
			continue
		}
		if found {
			return code, nil
		}
	}

	msg := fmt.Sprintf(
		"could not find %s %s in any of your regions (%s); set `region` on this resource to say where it belongs",
		ref.Kind, id, strings.Join(codes, ", "))
	if probeErr != nil {
		return "", fmt.Errorf("%s. A lookup also failed, which may be the real cause: %w", msg, probeErr)
	}
	return "", errors.New(msg)
}

func regionScoped(apiClient *civogo.Client, region string) *civogo.Client {
	c := *apiClient
	c.Region = region
	return &c
}

// inUseErrorCodes are the API error codes returned when a resource cannot be
// deleted yet because another resource still references it. During a
// `terraform destroy` these are transient: the referencing resource is being
// torn down at the same time, so the delete only needs retrying until the
// reference is gone.
var inUseErrorCodes = map[string]bool{
	"database_firewall_inuse_by_cluster":          true,
	"database_firewall_inuse_by_instance":         true,
	"database_firewall_used_by_loadbalancer":      true,
	"database_network_inuse_by_cluster":           true,
	"database_network_inuse_by_database":          true,
	"database_network_inuse_by_instance":          true,
	"database_network_inuse_by_instance_snapshot": true,
	"database_network_inuse_by_volumes":           true,
}

// APIErrorCode returns the Civo API error code carried by err, or "" if err
// does not wrap an API error response. civogo only surfaces the code for
// responses it has no dedicated error for; those it recognises are reduced to
// their reason string, so use IsResourceInUseError rather than calling this
// directly.
func APIErrorCode(err error) string {
	if err == nil {
		return ""
	}
	customErr, parseErr := ParseErrorResponse(err.Error())
	if parseErr != nil {
		return ""
	}
	return customErr.Code
}

// IsResourceInUseError reports whether err is a conflict caused by another
// resource still referencing the one being deleted, and so is worth retrying.
func IsResourceInUseError(err error) bool {
	if err == nil {
		return false
	}
	// The two conflicts civogo has dedicated errors for; the code does not
	// survive in the message for these, so match on the sentinel.
	if errors.Is(err, civogo.DatabaseNetworkInUseByVolumes) ||
		errors.Is(err, civogo.DatabaseNetworkDeleteWithInstanceError) {
		return true
	}
	return inUseErrorCodes[APIErrorCode(err)]
}
