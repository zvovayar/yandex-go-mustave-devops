package main

import (
	"net/http"
	"testing"

	internal "github.com/zvovayar/yandex-go-mustave-devops/internal/http"
)

func TestUpdateGaugeMetric(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			internal.UpdateGaugeMetric(tt.args.w, tt.args.r)
		})
	}
}

func TestUpdateCounterMetric(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			internal.UpdateCounterMetric(tt.args.w, tt.args.r)
		})
	}
}
