package common

import (
	"testing"
)

func TestIpToStr(t *testing.T) {
	type args struct {
		ip int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{ip: 2130706433},
			want: "127.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IpToStr(tt.args.ip); got != tt.want {
				t.Errorf("IpToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIpToInt(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			args:    args{ip: "127.0.0.1"},
			want:    2130706433,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IpToInt(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("IpToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IpToInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
