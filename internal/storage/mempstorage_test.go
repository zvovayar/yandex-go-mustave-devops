package storage

import (
	"bufio"
	"os"
	"reflect"
	"testing"
)

func TestMemPStorage_InitMemPStorage(t *testing.T) {
	type fields struct {
		sm            StoreMem
		chanPStoreMem chan StoreMem
		file          *os.File
		writer        *bufio.Writer
	}
	type args struct {
		ch chan StoreMem
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   chan StoreMem
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mps := &MemPStorage{
				sm:            tt.fields.sm,
				chanPStoreMem: tt.fields.chanPStoreMem,
				file:          tt.fields.file,
				writer:        tt.fields.writer,
			}
			if got := mps.InitMemPStorage(tt.args.ch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemPStorage.InitMemPStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
