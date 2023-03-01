package reporter

import (
	"fmt"
	"testing"
)

func Test_internalError_addMessage(t *testing.T) {
	type fields struct {
		errorNumber   int
		funcLocations []functionLocation
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "case0",
			fields: fields{errorNumber: 1, funcLocations: []functionLocation{{
				functionName: "getRedis",
				functionPath: "smart-qa/path_replay"},
			}},
			want: "errornumber(1)-r\ninternalerror(smart-qa/path_replay;getRedis;)-r\ninternalerrorstring(;)-r\ndisablenumber(0)-r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &internalError{
				errorNumber:        tt.fields.errorNumber,
				errorFuncLocations: tt.fields.funcLocations,
			}
			fmt.Println(c.addMessage())
			if got := c.addMessage(); got != tt.want {
				t.Errorf("addMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_internalError_Report(t *testing.T) {
	type fields struct {
		errorNumber   int
		funcLocations []functionLocation
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "case0",
			fields: fields{errorNumber: 1, funcLocations: []functionLocation{{
				functionName: "getRedis",
				functionPath: "smart-qa/path_replay"},
			}},
			want: "errornumber(1)-r\ninternalerror(smart-qa/path_replay;getRedis;)-r\ninternalerrorstring(;)-r\ndisablenumber(0)-r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &internalError{
				errorNumber:        tt.fields.errorNumber,
				errorFuncLocations: tt.fields.funcLocations,
			}
			if got := c.Report(); got != tt.want {
				t.Errorf("Report() = %v, want %v", got, tt.want)
			}
		})
	}
}
