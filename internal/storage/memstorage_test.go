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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &StoreMem{
				monitor: tt.fields.monitor,
			}
			sm.SetGMvalue(tt.args.gmname, tt.args.gm)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &StoreMem{
				monitor: tt.fields.monitor,
			}
			sm.SetCMvalue(tt.args.cmname, tt.args.cm)
		})
	}
}
