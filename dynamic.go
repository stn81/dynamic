package dynamic

import (
	"encoding/json"
	"reflect"
)

// DynamicFielder is the dynamic fielder interface
// Struct which implement this interface will have dynamic field support
type DynamicFielder interface {
	NewDynamicField(fieldName string) interface{}
}

var DynamicType = reflect.TypeOf(&Type{})

// IsDynamic return true if the type is dynamic
func IsDynamic(typ reflect.Type) bool {
	return typ == DynamicType
}

// Type defines the implementation of dynamic type
type Type struct {
	Value interface{}     `json:"-"`
	raw   json.RawMessage `json:"-"`
}

// New create a new dynamic instance
func New(v interface{}) *Type {
	return &Type{Value: v}
}

// GetValue return the value hold in dynamic, return nil if dynamic self is nil.
func GetValue(t *Type) interface{} {
	if t != nil {
		return t.Value
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *Type) UnmarshalJSON(data []byte) error {
	t.raw = data
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (t *Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

// GetRawMessage return the json.RawMessage hold by dynamic
func (t *Type) GetRawMessage() json.RawMessage {
	return t.raw
}
