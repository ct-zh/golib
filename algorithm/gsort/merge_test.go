package gsort

import (
	"reflect"
	"testing"
)

func Test_mergeSort(t *testing.T) {
	type args struct {
		a []int
	}
	tests := []struct {
		args args
		want []int
	}{
		{
			args: args{a: []int{5, 4, 3, 2, 1}},
			want: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			mergeSort(tt.args.a)
			if !reflect.DeepEqual(tt.args.a, tt.want) {
				t.Fatal("error", tt.args.a)
			} else {
				t.Logf("ok")
			}
		})
	}
}
