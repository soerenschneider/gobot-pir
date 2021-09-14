package internal

import (
	"os"
	"testing"
)

func Test_fromEnvBool(t *testing.T) {
	key := "asdjnasdogsagsadgjsdgsdgasdgjsdg"
	os.Setenv(key, "true")
	type args struct {
		name string
		def  bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "default",
			args: args{
				name: "asdjfasdighasgasdgasdg",
				def:  true,
			},
			want: true,
		},
		{
			name: "default",
			args: args{
				name: "asdjfasdighasgasdgasdg",
				def:  false,
			},
			want: false,
		},
		{
			name: "test",
			args: args{
				name: key,
				def:  false,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromEnvBool(tt.args.name, tt.args.def); got != tt.want {
				t.Errorf("fromEnvBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
