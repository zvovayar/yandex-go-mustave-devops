package storage

import (
	"bufio"
	"context"
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

func TestMemPStorage_GetCMvalue(t *testing.T) {
	type fields struct {
		sm            StoreMem
		chanPStoreMem chan StoreMem
		file          *os.File
		writer        *bufio.Writer
	}
	type args struct {
		cmname string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Counter
	}{
		{
			name: "Get testSetGet33",
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
				cmname: "testSetGet33",
			},
			want: 1,
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
				cmname: "testSetGet13",
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
			if got := mps.GetCMvalue(tt.args.cmname); got != tt.want {
				t.Errorf("MemPStorage.GetCMvalue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemPStorage_SaveBatch(t *testing.T) {
	type fields struct {
		sm            StoreMem
		chanPStoreMem chan StoreMem
		file          *os.File
		writer        *bufio.Writer
	}
	type args struct {
		ctx    context.Context
		batchM []Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test save batch",
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
				ctx: nil,
				batchM: []Metrics{
					{
						ID:    "BuckHashSys",
						MType: "gauge",
						Delta: new(int64),
						Value: new(float64),
						Hash:  "",
					},
					{
						ID:    "PollCount",
						MType: "counter",
						Delta: new(int64),
						Value: new(float64),
						Hash:  "",
					},
				},
			},
			wantErr: false,
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
			if err := mps.SaveBatch(tt.args.ctx, tt.args.batchM); (err != nil) != tt.wantErr {
				t.Errorf("MemPStorage.SaveBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemPStorage_ClosePersistanceStorage(t *testing.T) {
	type fields struct {
		sm            StoreMem
		chanPStoreMem chan StoreMem
		file          *os.File
		writer        *bufio.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test close resistance storage",
			fields: fields{
				sm: StoreMem{
					monitor: Monitor{
						Gmetrics: []Gauge{},
						Cmetrics: []Counter{},
					},
				},
				chanPStoreMem: make(chan StoreMem),
				file:          &os.File{},
				writer:        &bufio.Writer{},
			},
			wantErr: true,
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
			if err := mps.ClosePersistanceStorage(); (err != nil) != tt.wantErr {
				t.Errorf("MemPStorage.ClosePersistanceStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemPStorage_NewPersistanceStorage(t *testing.T) {
	type fields struct {
		sm            StoreMem
		chanPStoreMem chan StoreMem
		file          *os.File
		writer        *bufio.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test save file monitor",
			fields: fields{
				sm:            StoreMem{},
				chanPStoreMem: make(chan StoreMem),
				file:          &os.File{},
				writer:        &bufio.Writer{},
			},
			wantErr: false,
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
			if err := mps.NewPersistanceStorage(); (err != nil) != tt.wantErr {
				t.Errorf("MemPStorage.NewPersistanceStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
