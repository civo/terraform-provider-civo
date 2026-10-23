package datalist

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func floatApproxEquals(a, b float64) bool {
	return math.Abs(a-b) < 0.000001
}

func valueMatches(s *schema.Schema, value interface{}, filterValue interface{}, matchBy string) bool {
	switch s.Type {
	case schema.TypeString:
		switch matchBy {
		case "exact":
			return strings.EqualFold(filterValue.(string), value.(string))
		case "substring":
			return strings.Contains(value.(string), filterValue.(string))
		case "re":
			return filterValue.(*regexp.Regexp).MatchString(value.(string))
		}

	case schema.TypeBool:
		return filterValue.(bool) == value.(bool)

	case schema.TypeInt:
		return filterValue.(int) == value.(int)

	case schema.TypeFloat:
		return floatApproxEquals(filterValue.(float64), value.(float64))

	case schema.TypeList, schema.TypeSet:
		result := false
		for _, listValue := range listElements(value) {
			valueDoesMatch := valueMatches(s.Elem.(*schema.Schema), listValue, filterValue, matchBy)
			result = result || valueDoesMatch
		}
		return result
	}

	return false
}

// listElements returns the elements of a list or set field as the string, int, float64 or bool
// values valueMatches compares against. FlattenRecord functions may return a *schema.Set, a
// []interface{}, a typed slice such as []string, an array or a pointer to one of those, so accept
// all of them and normalise every element the same way.
func listElements(value interface{}) []interface{} {
	if set, ok := value.(*schema.Set); ok {
		if set == nil {
			return nil
		}
		value = set.List()
	}

	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		// A FlattenRecord bug: fail loudly rather than quietly matching nothing
		panic(fmt.Sprintf("datalist: expected a list or set value, got %T", value))
	}

	elements := make([]interface{}, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		if element, ok := primitiveValue(rv.Index(i)); ok {
			elements = append(elements, element)
		}
	}
	return elements
}

// primitiveValue converts one list element to the string, int, float64 or bool valueMatches
// compares against, so named types (type Tag string), sized or unsigned ints, float32 and
// pointers work too. It reports false for nil elements and for integers too large for int,
// which can never equal a filter value.
func primitiveValue(v reflect.Value) (interface{}, bool) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil, false
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.String(), true
	case reflect.Bool:
		return v.Bool(), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := v.Int()
		if n < math.MinInt || n > math.MaxInt {
			return nil, false
		}
		return int(n), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n := v.Uint()
		if n > math.MaxInt {
			return nil, false
		}
		return int(n), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	case reflect.Invalid:
		return nil, false
	}

	panic(fmt.Sprintf("datalist: list elements must be strings, numbers or bools, got %s", v.Type()))
}

func compareValues(s *schema.Schema, value1 interface{}, value2 interface{}) int {
	switch s.Type {
	case schema.TypeString:
		return strings.Compare(value1.(string), value2.(string))

	case schema.TypeBool:
		boolValue1 := value1.(bool)
		boolValue2 := value2.(bool)
		if boolValue1 == boolValue2 {
			return 0
		} else if !boolValue1 {
			return -1
		} else {
			return 1
		}

	case schema.TypeInt:
		intValue1 := value1.(int)
		intValue2 := value2.(int)
		if intValue1 < intValue2 {
			return -1
		} else if intValue1 > intValue2 {
			return 1
		} else {
			return 0
		}

	case schema.TypeFloat:
		floatValue1 := value1.(float64)
		floatValue2 := value2.(float64)
		if floatApproxEquals(floatValue1, floatValue2) {
			return 0
		} else if floatValue1 < floatValue2 {
			return -1
		} else if floatValue1 > floatValue2 {
			return 1
		} else {
			return 0
		}

	default:
		panic("Illegal state: Unsupported value type for sort")
	}
}
