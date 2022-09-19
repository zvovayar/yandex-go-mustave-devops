// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: internal/proto/monitor.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Result_CodeError int32

const (
	Result_OK  Result_CodeError = 0
	Result_ERR Result_CodeError = 1
)

// Enum value maps for Result_CodeError.
var (
	Result_CodeError_name = map[int32]string{
		0: "OK",
		1: "ERR",
	}
	Result_CodeError_value = map[string]int32{
		"OK":  0,
		"ERR": 1,
	}
)

func (x Result_CodeError) Enum() *Result_CodeError {
	p := new(Result_CodeError)
	*p = x
	return p
}

func (x Result_CodeError) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Result_CodeError) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_proto_monitor_proto_enumTypes[0].Descriptor()
}

func (Result_CodeError) Type() protoreflect.EnumType {
	return &file_internal_proto_monitor_proto_enumTypes[0]
}

func (x Result_CodeError) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Result_CodeError.Descriptor instead.
func (Result_CodeError) EnumDescriptor() ([]byte, []int) {
	return file_internal_proto_monitor_proto_rawDescGZIP(), []int{2, 0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID    string  `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`         // имя метрики
	MType string  `protobuf:"bytes,2,opt,name=MType,proto3" json:"MType,omitempty"`   // параметр, принимающий значение gauge или counter
	Delta int64   `protobuf:"varint,3,opt,name=Delta,proto3" json:"Delta,omitempty"`  // значение метрики в случае передачи counter
	Value float64 `protobuf:"fixed64,4,opt,name=Value,proto3" json:"Value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string  `protobuf:"bytes,5,opt,name=Hash,proto3" json:"Hash,omitempty"`     // значение хеш-функции
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_monitor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_monitor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_internal_proto_monitor_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Metric) GetMType() string {
	if x != nil {
		return x.MType
	}
	return ""
}

func (x *Metric) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Metric) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type BatchMetrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=Metrics,proto3" json:"Metrics,omitempty"`
	Count   int32     `protobuf:"varint,2,opt,name=Count,proto3" json:"Count,omitempty"`
}

func (x *BatchMetrics) Reset() {
	*x = BatchMetrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_monitor_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchMetrics) ProtoMessage() {}

func (x *BatchMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_monitor_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchMetrics.ProtoReflect.Descriptor instead.
func (*BatchMetrics) Descriptor() ([]byte, []int) {
	return file_internal_proto_monitor_proto_rawDescGZIP(), []int{1}
}

func (x *BatchMetrics) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

func (x *BatchMetrics) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code Result_CodeError `protobuf:"varint,1,opt,name=Code,proto3,enum=metrics.Result_CodeError" json:"Code,omitempty"`
	Text string           `protobuf:"bytes,2,opt,name=Text,proto3" json:"Text,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_monitor_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_monitor_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_internal_proto_monitor_proto_rawDescGZIP(), []int{2}
}

func (x *Result) GetCode() Result_CodeError {
	if x != nil {
		return x.Code
	}
	return Result_OK
}

func (x *Result) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

var File_internal_proto_monitor_proto protoreflect.FileDescriptor

var file_internal_proto_monitor_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x6e, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49,
	0x44, 0x12, 0x14, 0x0a, 0x05, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x44, 0x65, 0x6c, 0x74, 0x61,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x44, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a,
	0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x22, 0x4f, 0x0a, 0x0c, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x29, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x69, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x2d, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x19, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x04, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x54, 0x65, 0x78, 0x74, 0x22, 0x1c, 0x0a, 0x09, 0x43, 0x6f, 0x64, 0x65, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x45, 0x52,
	0x52, 0x10, 0x01, 0x32, 0xa6, 0x01, 0x0a, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x2d, 0x0a,
	0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0f, 0x2e, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x1a, 0x0f, 0x2e, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x30, 0x0a, 0x0c,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0f, 0x2e, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x1a, 0x0f, 0x2e,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x3c,
	0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x12, 0x15, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x1a, 0x0f, 0x2e, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x10, 0x5a, 0x0e,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_monitor_proto_rawDescOnce sync.Once
	file_internal_proto_monitor_proto_rawDescData = file_internal_proto_monitor_proto_rawDesc
)

func file_internal_proto_monitor_proto_rawDescGZIP() []byte {
	file_internal_proto_monitor_proto_rawDescOnce.Do(func() {
		file_internal_proto_monitor_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_monitor_proto_rawDescData)
	})
	return file_internal_proto_monitor_proto_rawDescData
}

var file_internal_proto_monitor_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_proto_monitor_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_proto_monitor_proto_goTypes = []interface{}{
	(Result_CodeError)(0), // 0: metrics.Result.CodeError
	(*Metric)(nil),        // 1: metrics.Metric
	(*BatchMetrics)(nil),  // 2: metrics.BatchMetrics
	(*Result)(nil),        // 3: metrics.Result
}
var file_internal_proto_monitor_proto_depIdxs = []int32{
	1, // 0: metrics.BatchMetrics.Metrics:type_name -> metrics.Metric
	0, // 1: metrics.Result.Code:type_name -> metrics.Result.CodeError
	1, // 2: metrics.Users.GetMetric:input_type -> metrics.Metric
	1, // 3: metrics.Users.UpdateMetric:input_type -> metrics.Metric
	2, // 4: metrics.Users.UpdateBatchMetrics:input_type -> metrics.BatchMetrics
	1, // 5: metrics.Users.GetMetric:output_type -> metrics.Metric
	3, // 6: metrics.Users.UpdateMetric:output_type -> metrics.Result
	3, // 7: metrics.Users.UpdateBatchMetrics:output_type -> metrics.Result
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_proto_monitor_proto_init() }
func file_internal_proto_monitor_proto_init() {
	if File_internal_proto_monitor_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_monitor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_monitor_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchMetrics); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_monitor_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_proto_monitor_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_monitor_proto_goTypes,
		DependencyIndexes: file_internal_proto_monitor_proto_depIdxs,
		EnumInfos:         file_internal_proto_monitor_proto_enumTypes,
		MessageInfos:      file_internal_proto_monitor_proto_msgTypes,
	}.Build()
	File_internal_proto_monitor_proto = out.File
	file_internal_proto_monitor_proto_rawDesc = nil
	file_internal_proto_monitor_proto_goTypes = nil
	file_internal_proto_monitor_proto_depIdxs = nil
}
