package datalist

import (
	"math"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestValueMatchesListAndSetValues(t *testing.T) {
	stringSet := &schema.Schema{Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeString}}
	stringList := &schema.Schema{Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}}
	var nilSet *schema.Set

	tests := []struct {
		name   string
		schema *schema.Schema
		value  interface{}
		want   bool
	}{
		{name: "set as *schema.Set", schema: stringSet, value: schema.NewSet(schema.HashString, []interface{}{"prod", "web"}), want: true},
		{name: "set as []interface{}", schema: stringSet, value: []interface{}{"prod", "web"}, want: true},
		{name: "set as []string", schema: stringSet, value: []string{"prod", "web"}, want: true},
		{name: "set as []string without a match", schema: stringSet, value: []string{"staging"}, want: false},
		{name: "set as nil []string", schema: stringSet, value: []string(nil), want: false},
		{name: "set as nil *schema.Set", schema: stringSet, value: nilSet, want: false},
		{name: "set as nil", schema: stringSet, value: nil, want: false},
		{name: "list as []interface{}", schema: stringList, value: []interface{}{"prod"}, want: true},
		{name: "list as []string", schema: stringList, value: []string{"prod"}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, valueMatches(tt.schema, tt.value, "prod", "exact"))
		})
	}
}

func TestApplyFiltersOnStringSliceSetField(t *testing.T) {
	// FlattenRecord functions such as civo_instances' pass tags as a plain []string
	recordSchema := map[string]*schema.Schema{
		"name": {Type: schema.TypeString},
		"tags": {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeString}},
	}
	records := []map[string]interface{}{
		{"name": "web-1", "tags": []string{"prod", "web"}},
		{"name": "db-1", "tags": []string{"staging"}},
		{"name": "untagged", "tags": []string(nil)},
	}

	filters, err := expandFilters(recordSchema, []interface{}{
		map[string]interface{}{"key": "tags", "values": []interface{}{"prod"}, "all": false, "match_by": "exact"},
	})
	assert.NoError(t, err)

	got := applyFilters(recordSchema, records, filters)

	if assert.Len(t, got, 1) {
		assert.Equal(t, "web-1", got[0]["name"])
	}
}

type tag string

func TestListElementsNormalisesContainersAndElements(t *testing.T) {
	prod := "prod"
	tags := []string{"prod"}
	var nilPtr *string
	var nilErr error = (*customErr)(nil)

	tests := []struct {
		name  string
		value interface{}
		want  []interface{}
	}{
		{name: "nil", value: nil, want: nil},
		{name: "nil *schema.Set", value: (*schema.Set)(nil), want: nil},
		{name: "*schema.Set", value: schema.NewSet(schema.HashString, []interface{}{"prod"}), want: []interface{}{"prod"}},
		{name: "[]interface{} with nil and sized int", value: []interface{}{"prod", nil, int64(2)}, want: []interface{}{"prod", 2}},
		{name: "[]string", value: []string{"prod", "web"}, want: []interface{}{"prod", "web"}},
		{name: "nil []string", value: []string(nil), want: []interface{}{}},
		{name: "named string type", value: []tag{"prod"}, want: []interface{}{"prod"}},
		{name: "pointers, skipping nil", value: []*string{nil, &prod}, want: []interface{}{"prod"}},
		{name: "interface holding a typed nil", value: []error{nilErr}, want: []interface{}{}},
		{name: "pointer to slice", value: &tags, want: []interface{}{"prod"}},
		{name: "nil pointer to slice", value: (*[]string)(nil), want: nil},
		{name: "array", value: [2]string{"prod", "web"}, want: []interface{}{"prod", "web"}},
		{name: "sized and unsigned ints", value: []interface{}{int8(1), uint(2), uint64(3)}, want: []interface{}{1, 2, 3}},
		{name: "uint64 too large for int is skipped", value: []uint64{1, math.MaxUint64}, want: []interface{}{1}},
		{name: "float32", value: []float32{0.5}, want: []interface{}{0.5}},
		{name: "bool", value: []bool{true}, want: []interface{}{true}},
		{name: "nil element pointer", value: []*string{nilPtr}, want: []interface{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, listElements(tt.value))
		})
	}
}

type customErr struct{}

func (*customErr) Error() string { return "custom" }

func TestValueMatchesOnTypedSlices(t *testing.T) {
	stringSet := &schema.Schema{Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeString}}
	intList := &schema.Schema{Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeInt}}

	assert.True(t, valueMatches(stringSet, []tag{"prod"}, "prod", "exact"), "named string type")
	assert.True(t, valueMatches(intList, []int64{1, 2}, 2, "exact"), "sized ints")
	assert.False(t, valueMatches(intList, []int64{1, 2}, 3, "exact"))
}

func TestListElementsFailsLoudlyOnBadShapes(t *testing.T) {
	// a FlattenRecord that returns the wrong shape should fail loudly, not quietly match nothing
	assert.PanicsWithValue(t, "datalist: expected a list or set value, got string", func() {
		listElements("prod,web")
	})
	assert.PanicsWithValue(t, "datalist: list elements must be strings, numbers or bools, got map[string]interface {}", func() {
		listElements([]map[string]interface{}{{"key": "prod"}})
	})
}

func TestApplyFiltersOnStringSliceSetFieldMatchModes(t *testing.T) {
	recordSchema := map[string]*schema.Schema{
		"name": {Type: schema.TypeString},
		"tags": {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeString}},
	}
	records := []map[string]interface{}{
		{"name": "web-1", "tags": []string{"prod", "web"}},
		{"name": "web-2", "tags": []string{"prod"}},
		{"name": "db-1", "tags": []string{"staging"}},
	}

	tests := []struct {
		name   string
		filter map[string]interface{}
		want   []string
	}{
		{name: "regex", filter: map[string]interface{}{"key": "tags", "values": []interface{}{"^stag"}, "all": false, "match_by": "re"}, want: []string{"db-1"}},
		{name: "substring", filter: map[string]interface{}{"key": "tags", "values": []interface{}{"we"}, "all": false, "match_by": "substring"}, want: []string{"web-1"}},
		{name: "all values", filter: map[string]interface{}{"key": "tags", "values": []interface{}{"prod", "web"}, "all": true, "match_by": "exact"}, want: []string{"web-1"}},
		{name: "any value", filter: map[string]interface{}{"key": "tags", "values": []interface{}{"web", "staging"}, "all": false, "match_by": "exact"}, want: []string{"web-1", "db-1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters, err := expandFilters(recordSchema, []interface{}{tt.filter})
			assert.NoError(t, err)

			var got []string
			for _, r := range applyFilters(recordSchema, records, filters) {
				got = append(got, r["name"].(string))
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
