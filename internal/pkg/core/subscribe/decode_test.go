package subscribe

import (
	"testing"
)

func TestSub2links(t *testing.T) {
	type args struct {
		subtext string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Sub2links(tt.args.subtext)
		})
	}
}
