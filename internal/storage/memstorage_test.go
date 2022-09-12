package storage

import (
	"testing"
)

func TestStoreMem_GetGMvalue(t *testing.T) {
	type fields struct {
		monitor Monitor
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
				monitor: Monitor{
					Gmetrics: []Gauge{1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
					Cmetrics: []Counter{0, 1},
				},
			},
			args: args{
				gmname: "testSetGet134",
			},
			want: 30,
		},
		{
			name: "Get testSetGet13 (unrecognized)",
			fields: fields{
				monitor: Monitor{
					Gmetrics: []Gauge{1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
					Cmetrics: []Counter{0, 1},
				},
			},
			args: args{
				gmname: "testSetGet13",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &StoreMem{
				monitor: tt.fields.monitor,
			}
			if got := sm.GetGMvalue(tt.args.gmname); got != tt.want {
				t.Errorf("StoreMem.GetGMvalue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreMem_GetCMvalue(t *testing.T) {
	type fields struct {
		monitor Monitor
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
				monitor: Monitor{
					Gmetrics: []Gauge{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 10},
					Cmetrics: []Counter{0, 1},
				},
			},
			args: args{
				cmname: "testSetGet33",
			},
			want: 1,
		},
		{
			name: "Get testSetGet3 (unrecognized)",
			fields: fields{
				monitor: Monitor{
					Gmetrics: []Gauge{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 10},
					Cmetrics: []Counter{1, 2},
				},
			},
			args: args{
				cmname: "testSetGet3",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &StoreMem{
				monitor: tt.fields.monitor,
			}
			if got := sm.GetCMvalue(tt.args.cmname); got != tt.want {
				t.Errorf("StoreMem.GetCMvalue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreMem_SetGMvalue(t *testing.T) {
	type fields struct {
		monitor Monitor
	}
	type args struct {
		gmname string
		gm     Gauge
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Set Frees",
			fields: fields{
				monitor: Monitor{
					Gmetrics: []Gauge{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
					Cmetrics: []Counter{},
				},
			},
			args: args{
				gmname: "Frees",
				gm:     5,
			},
		},
		{
			name: "Set new gmetric",
			fields: fields{
				monitor: Monitor{
					Gmetrics: []Gauge{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
					Cmetrics: []Counter{},
				},
			},
			args: args{
				gmname: "NewMetric",
				gm:     53,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &StoreMem{
				monitor: tt.fields.monitor,
			}
			sm.SetGMvalue(tt.args.gmname, tt.args.gm)

			if sm.monitor.Gmetrics[Gmetricnames[tt.args.gmname]] != tt.args.gm {
				t.Errorf("Fail set %v value", tt.args.gmname)
			}
		})
	}
}

func TestStoreMem_SetCMvalue(t *testing.T) {
	type fields struct {
		monitor Monitor
	}
	type args struct {
		cmname string
		cm     Counter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Set testSetGet33",
			fields: fields{
				monitor: Monitor{
					Gmetrics: []Gauge{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
					Cmetrics: []Counter{1, 2},
				},
			},
			args: args{
				cmname: "testSetGet33",
				cm:     5,
			},
		},
		{
			name: "Set new gmetric",
			fields: fields{
				monitor: Monitor{
					Gmetrics: []Gauge{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
					Cmetrics: []Counter{1, 2},
				},
			},
			args: args{
				cmname: "NewMetric",
				cm:     7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &StoreMem{
				monitor: tt.fields.monitor,
			}
			sm.SetCMvalue(tt.args.cmname, tt.args.cm)

			if sm.monitor.Cmetrics[Cmetricnames[tt.args.cmname]] != 7 {
				t.Errorf("Fail set %v value set %v got %v", tt.args.cmname,
					tt.args.cm,
					sm.monitor.Cmetrics[Cmetricnames[tt.args.cmname]])
			}

		})
	}
}
