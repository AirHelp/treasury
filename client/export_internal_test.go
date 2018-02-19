package client

import "testing"

func Test_validPrefix(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "correct prefix",
			args: args{
				input: "test/key/",
			},
			want: true,
		},
		{
			name: "correct key",
			args: args{
				input: "test/key/key",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validPrefix(tt.args.input); got != tt.want {
				t.Errorf("prefixCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
