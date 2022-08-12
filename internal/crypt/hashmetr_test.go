package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

func TestMakeMAC(t *testing.T) {
	type args struct {
		message []byte
		key     []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "check crypt hash MAC",
			args: args{
				message: []byte{0, 1, 2},
				key:     []byte{1, 1, 1},
			},
			want: []byte{173, 211, 104, 12, 213, 198, 8, 245, 209, 93, 211, 92, 45, 77, 109, 152, 9, 201, 185, 19, 222, 232, 216, 41, 3, 62, 97, 176, 186, 77, 21, 233},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeMAC(tt.args.message, tt.args.key)
			assert.Equal(t, tt.want, got, "they should be equal")

		})
	}
}

func TestValidMAC(t *testing.T) {
	type args struct {
		message    []byte
		messageMAC []byte
		key        []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test good MAC",
			args: args{
				message:    []byte{0, 1, 2},
				messageMAC: []byte{173, 211, 104, 12, 213, 198, 8, 245, 209, 93, 211, 92, 45, 77, 109, 152, 9, 201, 185, 19, 222, 232, 216, 41, 3, 62, 97, 176, 186, 77, 21, 233},
				key:        []byte{1, 1, 1},
			},
			want: true,
		},
		{
			name: "test bad MAC",
			args: args{
				message:    []byte{0, 2, 2},
				messageMAC: []byte{173, 211, 104, 12, 213, 198, 8, 245, 209, 93, 211, 92, 45, 77, 109, 152, 9, 201, 185, 19, 222, 232, 216, 41, 3, 62, 97, 176, 186, 77, 21, 233},
				key:        []byte{1, 1, 1},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidMAC(tt.args.message, tt.args.messageMAC, tt.args.key)
			assert.Equal(t, tt.want, got, "they should be equal")
		})
	}
}

func TestMetricsCrypt_MakeHashMetrics(t *testing.T) {
	type fields struct {
		M   inst.Metrics
		key string
		Msg string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "check crypt hash metrics for gauge",
			fields: fields{
				M: inst.Metrics{
					ID:    "metric name",
					MType: "gauge",
					Delta: new(int64),
					Value: new(float64),
					Hash:  "000",
				},
				key: "1234567890",
				Msg: "123123",
			},
			args: args{
				key: "0987654321",
			},
			want: "c6a34ac700e6fc843776ce518f1f86f7ab6be421efa4be35d157cee5a11d0cbb",
		},
		{
			name: "check crypt hash metrics for counter",
			fields: fields{
				M: inst.Metrics{
					ID:    "metric name",
					MType: "counter",
					Delta: new(int64),
					Value: new(float64),
					Hash:  "000",
				},
				key: "1234567890",
				Msg: "metric name:counter:0",
			},
			args: args{
				key: "0987654321",
			},
			want: "2bd92f8adb956d16e149bfa8b0538a287ad0b81ef22c396a543b385a866cd752",
		},
		{
			name: "incorrect MType",
			fields: fields{
				M: inst.Metrics{
					ID:    "metric name",
					MType: "digit",
					Delta: new(int64),
					Value: new(float64),
					Hash:  "000",
				},
				key: "1234567890",
				Msg: "12312312",
			},
			args: args{
				key: "0987654321",
			},
			want: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &MetricsCrypt{
				M:   tt.fields.M,
				key: tt.fields.key,
				Msg: tt.fields.Msg,
			}
			got := mc.MakeHashMetrics(tt.args.key)
			assert.Equal(t, tt.want, got, "they should be equal")
		})
	}
}

func TestMetricsCrypt_ControlHashMetrics(t *testing.T) {
	type fields struct {
		M   inst.Metrics
		key string
		Msg string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "check correct counter hash",
			fields: fields{
				M: inst.Metrics{
					ID:    "metric name",
					MType: "counter",
					Delta: new(int64),
					Value: new(float64),
					Hash:  "2bd92f8adb956d16e149bfa8b0538a287ad0b81ef22c396a543b385a866cd752",
				},

				key: "1234567890",
				Msg: "metric name:counter:0",
			},
			args: args{
				key: "1234567890",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &MetricsCrypt{
				M:   tt.fields.M,
				key: tt.fields.key,
				Msg: tt.fields.Msg,
			}
			got := mc.ControlHashMetrics(tt.args.key)
			assert.Equal(t, tt.want, got, "they should be equal")
		})
	}
}
