package main

import (
	"net/http"
	"testing"

	"github.com/zvovayar/yandex-go-mustave-devops/cmd/server/handlers"
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
			handlers.UpdateGaugeMetric(tt.args.w, tt.args.r)
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
			handlers.UpdateCounterMetric(tt.args.w, tt.args.r)
		})
	}
}
