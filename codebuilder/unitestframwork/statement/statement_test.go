package statement

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestMockedStatementResult(t *testing.T) {
	eleTypeStr := "*under.varType"
	index := strings.LastIndex(eleTypeStr, "/")
	if index != -1 {
		eleTypeStr = eleTypeStr[index+1:]
		rs := strings.Split(eleTypeStr, ".")
		if len(rs) > 0 {
			name := rs[1]
			ch, _ := utf8.DecodeRuneInString(name)
			fmt.Println(unicode.IsUpper(ch))
		}
	}
}
