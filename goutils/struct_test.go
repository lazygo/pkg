package goutils

import "testing"

type structTestUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Token    string `json:"-"`
	hidden   string
	NoTag    string
}

func TestStruct2Map(t *testing.T) {
	u := structTestUser{ID: 1, Name: "lzp", Password: "pwd", Age: 18, Token: "tk", hidden: "h", NoTag: "nt"}

	result, excluded := Struct2Map(u, "password")

	if result["id"] != 1 || result["name"] != "lzp" || result["age"] != 18 {
		t.Errorf("result = %v", result)
	}
	if _, ok := result["password"]; ok {
		t.Errorf("password should be excluded: %v", result)
	}
	if _, ok := result["-"]; ok {
		t.Errorf(`tag "-" should be excluded: %v`, result)
	}
	if _, ok := result["hidden"]; ok {
		t.Errorf("unexported field should be skipped: %v", result)
	}
	if _, ok := result["NoTag"]; ok {
		t.Errorf("field without json tag should be skipped: %v", result)
	}
	if excluded["password"] != "pwd" {
		t.Errorf("excluded = %v", excluded)
	}
}

func TestStruct2MapPointer(t *testing.T) {
	result, _ := Struct2Map(&structTestUser{ID: 7, Name: "a"})
	if result["id"] != 7 {
		t.Errorf("result = %v", result)
	}

	// nil 指针不 panic，返回空 map
	result, excluded := Struct2Map((*structTestUser)(nil))
	if len(result) != 0 || len(excluded) != 0 {
		t.Errorf("Struct2Map(nil) = %v, %v", result, excluded)
	}
}

func TestStruct2MapNotStruct(t *testing.T) {
	result, excluded := Struct2Map("not a struct")
	if len(result) != 0 || len(excluded) != 0 {
		t.Errorf("Struct2Map(string) = %v, %v", result, excluded)
	}
}
