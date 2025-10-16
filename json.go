package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"strconv"
	"strings"
)

// OrderedObject is an ordered sequence of name/value members in a JSON object.
//
// RFC 8259 defines an object as an "unordered collection".
// JSON implementations need not make "ordering of object members visible"
// to applications nor will they agree on the semantic meaning of an object if
// "the names within an object are not unique". For maximum compatibility,
// applications should avoid relying on ordering or duplicity of object names.
type OrderedObject[V any] []ObjectMember[V]

// ObjectMember is a JSON object member.
type ObjectMember[V any] struct {
	Value V
	Name  string
}

// MarshalJSONTo encodes obj as a JSON object into enc.
//
//nolint:wrapcheck // Example code from json/v2 documentation
func (obj *OrderedObject[V]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for i := range *obj {
		member := &(*obj)[i]
		if err := json.MarshalEncode(enc, &member.Name); err != nil {
			return err
		}
		if err := json.MarshalEncode(enc, &member.Value); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

// UnmarshalJSONFrom decodes a JSON object from dec into obj.
//
//nolint:wrapcheck // Example code from json/v2 documentation
func (obj *OrderedObject[V]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if k := dec.PeekKind(); k != '{' {
		return fmt.Errorf("expected object start, but encountered %v", k)
	}
	if _, err := dec.ReadToken(); err != nil {
		return err
	}
	for dec.PeekKind() != '}' {
		*obj = append(*obj, ObjectMember[V]{})
		member := &(*obj)[len(*obj)-1]
		if err := json.UnmarshalDecode(dec, &member.Name); err != nil {
			return err
		}
		if err := json.UnmarshalDecode(dec, &member.Value); err != nil {
			return err
		}
	}
	if _, err := dec.ReadToken(); err != nil {
		return err
	}
	return nil
}

// extractBudgetKeysAndValues extracts the keys and values from data.budget in their original order.
func extractBudgetKeysAndValues(jsonData []byte) ([]string, map[string]any, error) {
	// Parse the outer structure with ordered budget
	type DataWrapper struct {
		Data struct {
			Budget OrderedObject[any] `json:"budget"`
		} `json:"data"`
	}

	var wrapper DataWrapper
	if err := json.Unmarshal(jsonData, &wrapper); err != nil {
		return nil, nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract keys in order and build a map for easy value lookup
	keys := make([]string, len(wrapper.Data.Budget))
	values := make(map[string]any, len(wrapper.Data.Budget))
	for i, member := range wrapper.Data.Budget {
		keys[i] = member.Name
		values[member.Name] = member.Value
	}

	return keys, values, nil
}

// inspectJSONValue returns a Nushell-style description of a JSON value.
func inspectJSONValue(v any) string {
	switch val := v.(type) {
	case map[string]any:
		fieldCount := len(val)
		if fieldCount == 1 {
			return strings.Trim(fmt.Sprint(val), "map[]")
		}
		return fmt.Sprintf("{record %d fields}", fieldCount)
	case []any:
		itemCount := len(val)
		// Check if it's a table (array of objects) or a list (array of primitives)
		if itemCount > 0 {
			if _, isMap := val[0].(map[string]any); isMap {
				if itemCount == 1 {
					return fmt.Sprint(val)
				}
				return fmt.Sprintf("[table %d rows]", itemCount)
			}
		}
		if itemCount == 1 {
			return fmt.Sprint(val)
		}
		return fmt.Sprintf("[list %d items]", itemCount)
	case string:
		return formatMonthYear(val)
	case float64:
		// Check if it's an integer
		if val == float64(int64(val)) {
			return strconv.FormatInt(int64(val), 10)
		}
		return fmt.Sprintf("%v", val)
	case bool:
		return strconv.FormatBool(val)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}
