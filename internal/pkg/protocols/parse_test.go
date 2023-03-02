package protocols

import (
	"reflect"
	"testing"
)

func TestParseSSLink(t *testing.T) {
	type args struct {
		link string
	}
	tests := []struct {
		name string
		args args
		want *ShadowSocks
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
			if got := ParseSSLink(tt.args.link); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSSLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
