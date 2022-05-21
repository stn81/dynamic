package dynamic_test

import (
	"encoding/json"
	"fmt"

	"github.com/stn81/dynamic"
)

type aContent struct {
	Value int `json:"value"`
}

type bContent struct {
	Values []int `json:"items"`
}

type subTypeA struct {
	Name string `json:"name"`
}

type subTypeB struct {
	Age int `json:"age"`
}

type nestedContent struct {
	SubType string        `json:"sub_type"`
	Data    *dynamic.Type `json:"value,omitempty"`
}

func (c *nestedContent) NewDynamicField(fieldName string) interface{} {
	switch c.SubType {
	case "subTypeA":
		return &subTypeA{}
	case "subTypeB":
		return &subTypeB{}
	}
	return nil
}

type jsonValue struct {
	Type    string        `json:"type"`
	Content *dynamic.Type `json:"content,omitempty"`
}

func (jc *jsonValue) NewDynamicField(fieldName string) interface{} {
	switch jc.Type {
	case "a":
		return &aContent{}
	case "b":
		return &bContent{}
	case "nested":
		return &nestedContent{}
	}
	return nil
}

func ExampleMarshalA() {
	obj := &jsonValue{
		Type:    "a",
		Content: dynamic.New(&aContent{16}),
	}
	data, _ := json.Marshal(obj)
	fmt.Println(string(data))

	obj = &jsonValue{
		Type:    "b",
		Content: dynamic.New(&bContent{Values: []int{1, 2, 3}}),
	}
	data, _ = json.Marshal(obj)
	fmt.Println(string(data))
	// Output:
	// {"type":"a","content":{"value":16}}
	// {"type":"b","content":{"items":[1,2,3]}}
}

func ExampleUnmarshal() {
	input := []byte(`{"type":"b","content":{"items":[1,2,3]}}`)
	obj := &jsonValue{}
	_ = dynamic.ParseJSON(input, obj)
	content, ok := obj.Content.Value.(*bContent)
	fmt.Println(obj.Type)
	fmt.Println(ok)
	fmt.Println(content.Values)
	// Output:
	// b
	// true
	// [1 2 3]
}
