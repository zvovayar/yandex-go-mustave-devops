package storage

import (
	"bufio"
	"os"
	"testing"
)

func TestMemPStorage_GetGMvalue(t *testing.T) {
	type fields struct {
		sm            StoreMem
		chanPStoreMem chan StoreMem
		file          *os.File
		writer        *bufio.Writer
	}
	type args struct {
		gmname string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Gauge
	}{
		{
			name: "Get testSetGet134",
			fields: fields{
				sm: StoreMem{
					monitor: Monitor{
						Gmetrics: []Gauge{1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
						Cmetrics: []Counter{0, 1},
					},
				},
				chanPStoreMem: make(chan StoreMem),
				file:          &os.File{},
				writer:        &bufio.Writer{},
			},
			args: args{
				gmname: "testSetGet134",
			},
			want: 30,
		},
		{
			name: "Get testSetGet13 unrecognized",
			fields: fields{
				sm: StoreMem{
					monitor: Monitor{
						Gmetrics: []Gauge{1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
						Cmetrics: []Counter{0, 1},
					},
				},
				chanPStoreMem: make(chan StoreMem),
				file:          &os.File{},
				writer:        &bufio.Writer{},
			},
			args: args{
				gmname: "testSetGet13",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mps := &MemPStorage{
				sm:            tt.fields.sm,
				chanPStoreMem: tt.fields.chanPStoreMem,
				file:          tt.fields.file,
				writer:        tt.fields.writer,
			}
			if got := mps.GetGMvalue(tt.args.gmname); got != tt.want {
				t.Errorf("MemPStorage.GetGMvalue() = %v, want %v", got, tt.want)
			}
		})
	}
}
