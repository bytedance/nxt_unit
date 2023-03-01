package atghelper

import (
	"reflect"
	"testing"

	"gotest.tools/assert"
)

func TestRemoveTestCaseFile(t *testing.T) {
	RemoveTestCaseFile("_testdata", "testdocommentaction8e93ac6efdcb22eb2c7650b4a86dd377_test.go")
}

func TestUpperCaseFirstLetter(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"a"}, "A"},
		{"test1", args{""}, ""},
		{"test1", args{"doComment"}, "DoComment"},
		{"test1", args{"hello"}, "Hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpperCaseFirstLetter(tt.args.s); got != tt.want {
				t.Errorf("UpperCaseFirstLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpperCaseFirstLetter2(t *testing.T) {
	a := "sdwa"
	UpperCaseFirstLetter(a)
	assert.Equal(t, a, "sdwa")
}

var cCode = []byte(`
/*
 * I am C style comments
 */
func init() {
  fmt.Println("test");
}
`)

var commentCode = []byte(`
// Hello world
func CommentTest2(a int, b int) (*CommentTry, string) {
	if a == 0 {
		return nil, ""
	}
	fmt.Println("htt://www.tiktok.com")
	fmt.Println("\\'//abcse")
	fmt.Println("//")
	fmt.Println("// single-line comments")
	return &CommentTry{
		A: "su",
		B: 0,  // yesss
	}, ""
}
`)

var codeAfter = []byte(`

func CommentTest2(a int, b int) (*CommentTry, string) {
	if a == 0 {
		return nil, ""
	}
	fmt.Println("htt://www.tiktok.com")
	fmt.Println("\\'//abcse")
	fmt.Println("//")
	fmt.Println("// single-line comments")
	return &CommentTry{
		A: "su",
		B: 0,  
	}, ""
}
`)

func TestRemoveStyleComments(t *testing.T) {
	if string(RemoveGoStyleComments(commentCode, []string{"yesss", "Hello world"})) != string(codeAfter) {
		t.Error("Remove Cpp style comments failure!")
	}
}

func TestGetPkgName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case0",
			args: args{path: "errors"},
			want: "errors",
		},
		{
			name: "case1",
			args: args{path: "errors"},
			want: "errors",
		},
		{
			name: "case2",
			args: args{path: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPkgName(tt.args.path); got != tt.want {
				t.Errorf("GetPkgName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplacePkgName(t *testing.T) {
	type args struct {
		s               string
		pkgName         string
		originalPkgName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{s: "abc2.hello", pkgName: "su", originalPkgName: "abc2"}, want: "su.hello"},
		{name: "case2", args: args{s: "*abc2.hello", pkgName: "su", originalPkgName: "abc2"}, want: "*su.hello"},
		{name: "case3", args: args{s: "[]test1.SmartAtg", pkgName: "su", originalPkgName: "test1"}, want: "[]su.SmartAtg"},
		{name: "case4", args: args{s: "*abc2.hello", pkgName: "su", originalPkgName: "abc2"}, want: "*su.hello"},
		{name: "case5", args: args{s: "map[abc2]atg.hello", pkgName: "su", originalPkgName: "atg"}, want: "map[abc2]su.hello"},
		{name: "case6", args: args{s: "hello", pkgName: "su", originalPkgName: "abc2"}, want: "hello"},
		{name: "case7", args: args{s: "[]*pack.Comment", pkgName: "su", originalPkgName: "pack"}, want: "[]*su.Comment"},
		{name: "case8", args: args{s: "[][]smartunit.hello", pkgName: "su", originalPkgName: "smartunit"}, want: "[][]su.hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplacePkgName(tt.args.s, tt.args.pkgName, tt.args.originalPkgName); got != tt.want {
				t.Errorf("ReplacePkgName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTypeExported(t *testing.T) {
	lowerCaseArray := make([]string, 0)
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want bool
	}{{name: "case1", args: args{reflect.TypeOf(lowerCaseArray)}, want: true}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsTypeExported(tt.args.t), "IsTypeExported(%v)", tt.args.t)
		})
	}
}

func TestReplacePkgNameForMap(t *testing.T) {
	type args struct {
		s               string
		pkgName         string
		originalPkgName string
		fromLeft        bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{s: "abc2.hello", pkgName: "su", originalPkgName: "abc2", fromLeft: true}, want: "su.hello"},
		{name: "case2", args: args{s: "*abc2.hello", pkgName: "su", originalPkgName: "abc2", fromLeft: true}, want: "*su.hello"},
		{name: "case3", args: args{s: "[]test1.SmartAtg", pkgName: "su", originalPkgName: "test1", fromLeft: true}, want: "[]su.SmartAtg"},
		{name: "case4", args: args{s: "*abc2.hello", pkgName: "su", originalPkgName: "abc2", fromLeft: true}, want: "*su.hello"},
		{name: "case5", args: args{s: "map[abc2]atg.hello", pkgName: "su", originalPkgName: "atg", fromLeft: true}, want: "map[abc2]su.hello"},
		{name: "case6", args: args{s: "hello", pkgName: "su", originalPkgName: "abc2", fromLeft: true}, want: "hello"},
		{name: "case7", args: args{s: "[]*pack.Comment", pkgName: "su", originalPkgName: "pack", fromLeft: true}, want: "[]*su.Comment"},
		{name: "case8", args: args{s: "[][]smartunit.hello", pkgName: "su", originalPkgName: "smartunit", fromLeft: true}, want: "[][]su.hello"},
		{name: "case5", args: args{s: "map[atg123.hello]atg.hello", pkgName: "su", originalPkgName: "atg", fromLeft: false}, want: "map[atg123]su.hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ReplacePkgNameForMap(tt.args.s, tt.args.pkgName, tt.args.originalPkgName, tt.args.fromLeft), "ReplacePkgNameForMap(%v, %v, %v, %v)", tt.args.s, tt.args.pkgName, tt.args.originalPkgName, tt.args.fromLeft)
		})
	}
}
