package atghelper

import (
	"path"
	"reflect"
	"runtime"
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

var commentCode = []byte(`
// Hello world
func CommentTest2(a int, b int) (*CommentTry, string) {
	if a == 0 {
		return nil, ""
	}
	fmt.Println("http://www.bytedance.com")
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
	fmt.Println("http://www.bytedance.com")
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

func TestIsFileExist(t *testing.T) {

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"case1", args{"/not/exist"}, false},
		{"case2", args{"./"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFileExist(tt.args.path); got != tt.want {
				t.Errorf("IsFileExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		slice []string
		item  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"case1", args{GoKeywords, "interface"}, true},
		{"case2", args{GoKeywords, "notexist"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.slice, tt.args.item); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRepoByRelativePath(t *testing.T) {
	type args struct {
		relativePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"case1", args{"github.com/myuser/myrepo/mypackage"}, "github.com/myuser/myrepo"},
		{"case2", args{"strings"}, "strings"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRepoByRelativePath(tt.args.relativePath); got != tt.want {
				t.Errorf("GetRepoByRelativePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDirectory(t *testing.T) {
	type args struct {
		dirName string
	}
	tests := []struct {
		name string
		args args
	}{
		{"case1", args{"not/exist"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemoveDirectory(tt.args.dirName)
		})
	}
}

func TestRandStringBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"test", args{10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandStringBytes(tt.args.n); len(got) != tt.want {
				t.Errorf("RandStringBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveCStyleComments(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"test1", args{[]byte("/*comment*/\nfunc Add(){}")}, []byte("\nfunc Add(){}")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveCStyleComments(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveCStyleComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveGoComments(t *testing.T) {
	_, filePath, _, _ := runtime.Caller(0)
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{path.Join(filePath, "../../atg/template/atgv2.go")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RemoveGoComments(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveGoComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
