package binder

import (
	"log"
	"reflect"
	"testing"
)

func TestFromQuery_Bind(t *testing.T) {
	type s struct {
		Int             int    `json:"int"`
		Int8            int8   `json:"int_8"`
		Int16           int16  `json:"int_16"`
		Int32           int32  `json:"int_32"`
		Int64           int64  `json:"int_64"`
		String          string `json:"string"`
		Bool            bool   `json:"bool"`
		CapitalizedBool bool   `json:"capitalized_bool"`
	}

	want := s{
		Int:             4,
		Int8:            8,
		Int16:           16,
		Int32:           32,
		Int64:           64,
		String:          "value",
		Bool:            true,
		CapitalizedBool: true,
	}

	query := map[string][]string{
		"int":              {"4"},
		"int_8":            {"8"},
		"int_16":           {"16"},
		"int_32":           {"32"},
		"int_64":           {"64"},
		"string":           {"value"},
		"bool":             {"true"},
		"capitalized_bool": {"True"},
	}

	var got s
	err := FromQuery(&got, query)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected value, got:%+v, want:%+v", got, want)
	}
}

func TestFromQuery_BindSlice(t *testing.T) {
	type s struct {
		IntSlice         []int    `json:"int_slice"`
		StringSlice      []string `json:"string_slice"`
		MissedValueSlice []int    `json:"missed_value_slice"`
	}

	want := s{
		IntSlice:         []int{1, 2, 3, 4},
		StringSlice:      []string{"one", "two"},
		MissedValueSlice: []int{1, 3, 4},
	}

	query := map[string][]string{
		"int_slice":          {"1,2,3,4"},
		"string_slice":       {"one,two", "three"},
		"missed_value_slice": {"1,,3,4,"},
	}

	var got s
	err := FromQuery(&got, query)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected value, got:%+v, want:%+v", got, want)
	}
}

func TestFromQuery_BindError(t *testing.T) {
	err1 := FromQuery(struct{}{}, map[string][]string{})
	if err1 == nil || err1.Error() != "invalid interface type" {
		log.Fatal("no expected error")
	}
}
