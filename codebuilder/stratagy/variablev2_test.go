package stratagy

import (
	"context"
	"reflect"
	"testing"
)

func TestGetSpecialVariableV2(t *testing.T) {
	type args struct {
		ctx context.Context
		v   reflect.Value
	}
	//variableContext := atgconstant.VariableContext{}
	//ctx := context.Background()
	//vtx := contexthelper.SetVariableContext(ctx, variableContext)
	// TODO: below usage is related with redis initialization.
	// var cl redis.Client
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		// Below test is running on the local machine because we cannot introduce the package
		//{name: "case3",
		//	args:  args{ctx: vtx, v: reflect.ValueOf(&app.RequestContext{})},
		//	want:  "&app.RequestContext{}",
		//	want1: true,
		//},
		//{name: "case4",
		//	args:  args{ctx: vtx, v: reflect.ValueOf(&gin.Context{})},
		//	want:  "&gin.Context{Request: &http.Request{URL: &url.URL{Path: \"test_path\",}}}",
		//	want1: true,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetSpecialVariableV2(tt.args.ctx, tt.args.v)
			if got != tt.want {
				t.Errorf("GetSpecialVariableV2() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetSpecialVariableV2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
