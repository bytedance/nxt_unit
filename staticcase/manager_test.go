package staticcase

import "testing"

func Test_checkFileName(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test",
			args{
				"smartunit_test.go",
			},
			false,
		},
		{
			"code",
			args{
				"smartunit.go",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkFileName(tt.args.fileName); got != tt.want {
				t.Errorf("checkFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
