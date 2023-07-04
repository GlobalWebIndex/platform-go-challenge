// proto/core/asset/v1/asset.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.2
// source: proto/core/asset/v1/asset.proto

package assetpb

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
	v1 "x-gwi/proto/core/_share/v1"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// AssetCore
type AssetCore struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// qid - Asset StoreQID
	Qid *v1.ShareQID `protobuf:"bytes,3,opt,name=qid,proto3" json:"qid,omitempty"`
	// md - MetaDescription
	Md *v1.MetaDescription `protobuf:"bytes,5,opt,name=md,proto3,oneof" json:"md,omitempty"`
	// mmd - MetaMultiDescription
	Mmd *v1.MetaMultiDescription `protobuf:"bytes,6,opt,name=mmd,proto3,oneof" json:"mmd,omitempty"`
	// chart - Chart (that has a small title, axes titles and data)
	Chart *AssetChart `protobuf:"bytes,7,opt,name=chart,proto3,oneof" json:"chart,omitempty"`
	// insight - Insight (a small piece of text that provides some insight into a
	// topic, e.g. "40% of millenials spend more than 3hours on social media
	// daily")
	Insight *AssetInsight `protobuf:"bytes,8,opt,name=insight,proto3,oneof" json:"insight,omitempty"`
	// audience - Audience (which is a series of characteristics, for that
	// exercise lets focus on gender (Male, Female), birth country, age groups,
	// hours spent daily on social media, number of purchases last month) e.g.
	// Males from 24-35 that spent more than 3 hours on social media daily.
	Audience *AssetAudience `protobuf:"bytes,9,opt,name=audience,proto3,oneof" json:"audience,omitempty"`
}

func (x *AssetCore) Reset() {
	*x = AssetCore{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_core_asset_v1_asset_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetCore) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetCore) ProtoMessage() {}

func (x *AssetCore) ProtoReflect() protoreflect.Message {
	mi := &file_proto_core_asset_v1_asset_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetCore.ProtoReflect.Descriptor instead.
func (*AssetCore) Descriptor() ([]byte, []int) {
	return file_proto_core_asset_v1_asset_proto_rawDescGZIP(), []int{0}
}

func (x *AssetCore) GetQid() *v1.ShareQID {
	if x != nil {
		return x.Qid
	}
	return nil
}

func (x *AssetCore) GetMd() *v1.MetaDescription {
	if x != nil {
		return x.Md
	}
	return nil
}

func (x *AssetCore) GetMmd() *v1.MetaMultiDescription {
	if x != nil {
		return x.Mmd
	}
	return nil
}

func (x *AssetCore) GetChart() *AssetChart {
	if x != nil {
		return x.Chart
	}
	return nil
}

func (x *AssetCore) GetInsight() *AssetInsight {
	if x != nil {
		return x.Insight
	}
	return nil
}

func (x *AssetCore) GetAudience() *AssetAudience {
	if x != nil {
		return x.Audience
	}
	return nil
}

// AssetChart - (that has a small title, axes titles and data)
type AssetChart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// title - (0-256 characters) optional
	Title *string `protobuf:"bytes,4,opt,name=title,proto3,oneof" json:"title,omitempty"`
	// md - MetaDescription
	Md *v1.MetaDescription `protobuf:"bytes,5,opt,name=md,proto3,oneof" json:"md,omitempty"`
	// mmd - MetaMultiDescription
	Mmd *v1.MetaMultiDescription `protobuf:"bytes,6,opt,name=mmd,proto3,oneof" json:"mmd,omitempty"`
	// data
	Data *structpb.Struct `protobuf:"bytes,7,opt,name=data,proto3,oneof" json:"data,omitempty"`
	// data_raw - raw binary data (65536)
	DataRaw []byte `protobuf:"bytes,8,opt,name=data_raw,json=dataRaw,proto3,oneof" json:"data_raw,omitempty"`
	// options
	Options *structpb.Struct `protobuf:"bytes,9,opt,name=options,proto3,oneof" json:"options,omitempty"`
}

func (x *AssetChart) Reset() {
	*x = AssetChart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_core_asset_v1_asset_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetChart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetChart) ProtoMessage() {}

func (x *AssetChart) ProtoReflect() protoreflect.Message {
	mi := &file_proto_core_asset_v1_asset_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetChart.ProtoReflect.Descriptor instead.
func (*AssetChart) Descriptor() ([]byte, []int) {
	return file_proto_core_asset_v1_asset_proto_rawDescGZIP(), []int{1}
}

func (x *AssetChart) GetTitle() string {
	if x != nil && x.Title != nil {
		return *x.Title
	}
	return ""
}

func (x *AssetChart) GetMd() *v1.MetaDescription {
	if x != nil {
		return x.Md
	}
	return nil
}

func (x *AssetChart) GetMmd() *v1.MetaMultiDescription {
	if x != nil {
		return x.Mmd
	}
	return nil
}

func (x *AssetChart) GetData() *structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *AssetChart) GetDataRaw() []byte {
	if x != nil {
		return x.DataRaw
	}
	return nil
}

func (x *AssetChart) GetOptions() *structpb.Struct {
	if x != nil {
		return x.Options
	}
	return nil
}

// AssetInsight - (a small piece of text that provides some insight into a
// topic, e.g. "40% of millenials spend more than 3hours on social media daily")
type AssetInsight struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// sentence - (1-256 characters) required
	Sentence *string `protobuf:"bytes,4,opt,name=sentence,proto3,oneof" json:"sentence,omitempty"`
	// md - MetaDescription
	Md *v1.MetaDescription `protobuf:"bytes,5,opt,name=md,proto3,oneof" json:"md,omitempty"`
	// mmd - MetaMultiDescription
	Mmd *v1.MetaMultiDescription `protobuf:"bytes,6,opt,name=mmd,proto3,oneof" json:"mmd,omitempty"`
}

func (x *AssetInsight) Reset() {
	*x = AssetInsight{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_core_asset_v1_asset_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetInsight) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetInsight) ProtoMessage() {}

func (x *AssetInsight) ProtoReflect() protoreflect.Message {
	mi := &file_proto_core_asset_v1_asset_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetInsight.ProtoReflect.Descriptor instead.
func (*AssetInsight) Descriptor() ([]byte, []int) {
	return file_proto_core_asset_v1_asset_proto_rawDescGZIP(), []int{2}
}

func (x *AssetInsight) GetSentence() string {
	if x != nil && x.Sentence != nil {
		return *x.Sentence
	}
	return ""
}

func (x *AssetInsight) GetMd() *v1.MetaDescription {
	if x != nil {
		return x.Md
	}
	return nil
}

func (x *AssetInsight) GetMmd() *v1.MetaMultiDescription {
	if x != nil {
		return x.Mmd
	}
	return nil
}

// AssetAudience -  (which is a series of characteristics, for that exercise
// lets focus on gender (Male, Female), birth country, age groups, hours spent
// daily on social media, number of purchases last month) e.g. Males from 24-35
// that spent more than 3 hours on social media daily.
type AssetAudience struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// md - MetaDescription
	Md *v1.MetaDescription `protobuf:"bytes,5,opt,name=md,proto3,oneof" json:"md,omitempty"`
	// mmd - MetaMultiDescription
	Mmd *v1.MetaMultiDescription `protobuf:"bytes,6,opt,name=mmd,proto3,oneof" json:"mmd,omitempty"`
	// gender - (male | female)
	Gender *string `protobuf:"bytes,7,opt,name=gender,proto3,oneof" json:"gender,omitempty"`
	// genders - [](male, female)
	Genders []string `protobuf:"bytes,8,rep,name=genders,proto3" json:"genders,omitempty"`
	// country_code - iso 2 birth country code (auto lowercase)
	CountryCode *string `protobuf:"bytes,9,opt,name=country_code,json=countryCode,proto3,oneof" json:"country_code,omitempty"`
	// country_codes - [](0-300 elements, unique, code max_len 2) optional
	CountryCodes []string `protobuf:"bytes,10,rep,name=country_codes,json=countryCodes,proto3" json:"country_codes,omitempty"`
	// age_min - age groups (1-100)
	AgeMin *uint32 `protobuf:"varint,11,opt,name=age_min,json=ageMin,proto3,oneof" json:"age_min,omitempty"`
	// age_max - age groups (1-100)
	AgeMax *uint32 `protobuf:"varint,12,opt,name=age_max,json=ageMax,proto3,oneof" json:"age_max,omitempty"`
	// hours_min - hours on social (1-24) hours spent daily on social media
	HoursMin *uint32 `protobuf:"varint,13,opt,name=hours_min,json=hoursMin,proto3,oneof" json:"hours_min,omitempty"`
	// hours_max - hours on social (1-24) hours spent daily on social media
	HoursMax *uint32 `protobuf:"varint,14,opt,name=hours_max,json=hoursMax,proto3,oneof" json:"hours_max,omitempty"`
	// purchases_min - (1-100_000) number of purchases last month
	PurchasesMin *uint32 `protobuf:"varint,15,opt,name=purchases_min,json=purchasesMin,proto3,oneof" json:"purchases_min,omitempty"`
	// purchases_max - (1-100_000) number of purchases last month
	PurchasesMax *uint32 `protobuf:"varint,16,opt,name=purchases_max,json=purchasesMax,proto3,oneof" json:"purchases_max,omitempty"`
}

func (x *AssetAudience) Reset() {
	*x = AssetAudience{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_core_asset_v1_asset_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetAudience) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetAudience) ProtoMessage() {}

func (x *AssetAudience) ProtoReflect() protoreflect.Message {
	mi := &file_proto_core_asset_v1_asset_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetAudience.ProtoReflect.Descriptor instead.
func (*AssetAudience) Descriptor() ([]byte, []int) {
	return file_proto_core_asset_v1_asset_proto_rawDescGZIP(), []int{3}
}

func (x *AssetAudience) GetMd() *v1.MetaDescription {
	if x != nil {
		return x.Md
	}
	return nil
}

func (x *AssetAudience) GetMmd() *v1.MetaMultiDescription {
	if x != nil {
		return x.Mmd
	}
	return nil
}

func (x *AssetAudience) GetGender() string {
	if x != nil && x.Gender != nil {
		return *x.Gender
	}
	return ""
}

func (x *AssetAudience) GetGenders() []string {
	if x != nil {
		return x.Genders
	}
	return nil
}

func (x *AssetAudience) GetCountryCode() string {
	if x != nil && x.CountryCode != nil {
		return *x.CountryCode
	}
	return ""
}

func (x *AssetAudience) GetCountryCodes() []string {
	if x != nil {
		return x.CountryCodes
	}
	return nil
}

func (x *AssetAudience) GetAgeMin() uint32 {
	if x != nil && x.AgeMin != nil {
		return *x.AgeMin
	}
	return 0
}

func (x *AssetAudience) GetAgeMax() uint32 {
	if x != nil && x.AgeMax != nil {
		return *x.AgeMax
	}
	return 0
}

func (x *AssetAudience) GetHoursMin() uint32 {
	if x != nil && x.HoursMin != nil {
		return *x.HoursMin
	}
	return 0
}

func (x *AssetAudience) GetHoursMax() uint32 {
	if x != nil && x.HoursMax != nil {
		return *x.HoursMax
	}
	return 0
}

func (x *AssetAudience) GetPurchasesMin() uint32 {
	if x != nil && x.PurchasesMin != nil {
		return *x.PurchasesMin
	}
	return 0
}

func (x *AssetAudience) GetPurchasesMax() uint32 {
	if x != nil && x.PurchasesMax != nil {
		return *x.PurchasesMax
	}
	return 0
}

var File_proto_core_asset_v1_asset_proto protoreflect.FileDescriptor

var file_proto_core_asset_v1_asset_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x73, 0x73,
	0x65, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x73,
	0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xd2, 0x03, 0x0a, 0x09, 0x41, 0x73, 0x73, 0x65, 0x74, 0x43, 0x6f, 0x72, 0x65, 0x12, 0x3f,
	0x0a, 0x03, 0x71, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x51, 0x49, 0x44, 0x42, 0x0d, 0xe0, 0x41, 0x01,
	0xfa, 0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00, 0x52, 0x03, 0x71, 0x69, 0x64, 0x12,
	0x3a, 0x0a, 0x02, 0x6d, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x48, 0x00, 0x52, 0x02, 0x6d, 0x64, 0x88, 0x01, 0x01, 0x12, 0x41, 0x0a, 0x03, 0x6d,
	0x6d, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x65, 0x74, 0x61, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x48, 0x01, 0x52, 0x03, 0x6d, 0x6d, 0x64, 0x88, 0x01, 0x01, 0x12, 0x3a,
	0x0a, 0x05, 0x63, 0x68, 0x61, 0x72, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x43, 0x68, 0x61, 0x72, 0x74, 0x48, 0x02,
	0x52, 0x05, 0x63, 0x68, 0x61, 0x72, 0x74, 0x88, 0x01, 0x01, 0x12, 0x40, 0x0a, 0x07, 0x69, 0x6e,
	0x73, 0x69, 0x67, 0x68, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x69, 0x67, 0x68, 0x74, 0x48, 0x03,
	0x52, 0x07, 0x69, 0x6e, 0x73, 0x69, 0x67, 0x68, 0x74, 0x88, 0x01, 0x01, 0x12, 0x43, 0x0a, 0x08,
	0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x73, 0x73, 0x65,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x41, 0x75, 0x64, 0x69, 0x65, 0x6e,
	0x63, 0x65, 0x48, 0x04, 0x52, 0x08, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x88, 0x01,
	0x01, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x6d, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6d, 0x6d, 0x64,
	0x42, 0x08, 0x0a, 0x06, 0x5f, 0x63, 0x68, 0x61, 0x72, 0x74, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x69,
	0x6e, 0x73, 0x69, 0x67, 0x68, 0x74, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x61, 0x75, 0x64, 0x69, 0x65,
	0x6e, 0x63, 0x65, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a,
	0x04, 0x08, 0x04, 0x10, 0x05, 0x22, 0x8d, 0x03, 0x0a, 0x0a, 0x41, 0x73, 0x73, 0x65, 0x74, 0x43,
	0x68, 0x61, 0x72, 0x74, 0x12, 0x29, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08, 0x72, 0x06, 0x18, 0x80, 0x02,
	0xd0, 0x01, 0x01, 0x48, 0x00, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x3a, 0x0a, 0x02, 0x6d, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x48, 0x01, 0x52, 0x02, 0x6d, 0x64, 0x88, 0x01, 0x01, 0x12, 0x41, 0x0a, 0x03, 0x6d,
	0x6d, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x65, 0x74, 0x61, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x48, 0x02, 0x52, 0x03, 0x6d, 0x6d, 0x64, 0x88, 0x01, 0x01, 0x12, 0x30,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x48, 0x03, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x88, 0x01, 0x01,
	0x12, 0x1e, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x72, 0x61, 0x77, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0c, 0x48, 0x04, 0x52, 0x07, 0x64, 0x61, 0x74, 0x61, 0x52, 0x61, 0x77, 0x88, 0x01, 0x01,
	0x12, 0x36, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x48, 0x05, 0x52, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x6d, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6d, 0x6d,
	0x64, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x64,
	0x61, 0x74, 0x61, 0x5f, 0x72, 0x61, 0x77, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a,
	0x04, 0x08, 0x03, 0x10, 0x04, 0x22, 0xeb, 0x01, 0x0a, 0x0c, 0x41, 0x73, 0x73, 0x65, 0x74, 0x49,
	0x6e, 0x73, 0x69, 0x67, 0x68, 0x74, 0x12, 0x2e, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x74, 0x65, 0x6e,
	0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0xe0, 0x41, 0x02, 0xfa, 0x42, 0x07,
	0x72, 0x05, 0x10, 0x01, 0x18, 0x80, 0x02, 0x48, 0x00, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x74, 0x65,
	0x6e, 0x63, 0x65, 0x88, 0x01, 0x01, 0x12, 0x3a, 0x0a, 0x02, 0x6d, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x01, 0x52, 0x02, 0x6d, 0x64, 0x88,
	0x01, 0x01, 0x12, 0x41, 0x0a, 0x03, 0x6d, 0x6d, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x4d, 0x75, 0x6c, 0x74, 0x69,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x02, 0x52, 0x03, 0x6d,
	0x6d, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x73, 0x65, 0x6e, 0x74, 0x65, 0x6e,
	0x63, 0x65, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x6d, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6d, 0x6d,
	0x64, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08,
	0x03, 0x10, 0x04, 0x22, 0xcd, 0x06, 0x0a, 0x0d, 0x41, 0x73, 0x73, 0x65, 0x74, 0x41, 0x75, 0x64,
	0x69, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x3a, 0x0a, 0x02, 0x6d, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x02, 0x6d, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x41, 0x0a, 0x03, 0x6d, 0x6d, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x44,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x01, 0x52, 0x03, 0x6d, 0x6d,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x36, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x19, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x13, 0x72, 0x11, 0x52, 0x04,
	0x6d, 0x61, 0x6c, 0x65, 0x52, 0x06, 0x66, 0x65, 0x6d, 0x61, 0x6c, 0x65, 0xd0, 0x01, 0x01, 0x48,
	0x02, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x3d, 0x0a, 0x07,
	0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x42, 0x23, 0xe0,
	0x41, 0x01, 0xfa, 0x42, 0x1d, 0x92, 0x01, 0x1a, 0x10, 0x1e, 0x18, 0x01, 0x22, 0x12, 0x72, 0x10,
	0x52, 0x00, 0x52, 0x04, 0x6d, 0x61, 0x6c, 0x65, 0x52, 0x06, 0x66, 0x65, 0x6d, 0x61, 0x6c, 0x65,
	0x28, 0x01, 0x52, 0x07, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x12, 0x36, 0x0a, 0x0c, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08, 0x72, 0x06, 0x98, 0x01, 0x02, 0xd0, 0x01,
	0x01, 0x48, 0x03, 0x52, 0x0b, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x64, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x3b, 0x0a, 0x0d, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x63,
	0x6f, 0x64, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x09, 0x42, 0x16, 0xe0, 0x41, 0x01, 0xfa,
	0x42, 0x10, 0x92, 0x01, 0x0d, 0x10, 0xac, 0x02, 0x18, 0x01, 0x22, 0x04, 0x72, 0x02, 0x18, 0x02,
	0x28, 0x01, 0x52, 0x0c, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x73,
	0x12, 0x2c, 0x0a, 0x07, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x69, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0d, 0x42, 0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08, 0x2a, 0x06, 0x18, 0x64, 0x28, 0x01, 0x40,
	0x01, 0x48, 0x04, 0x52, 0x06, 0x61, 0x67, 0x65, 0x4d, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x2c,
	0x0a, 0x07, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x42,
	0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08, 0x2a, 0x06, 0x18, 0x64, 0x28, 0x01, 0x40, 0x01, 0x48,
	0x05, 0x52, 0x06, 0x61, 0x67, 0x65, 0x4d, 0x61, 0x78, 0x88, 0x01, 0x01, 0x12, 0x30, 0x0a, 0x09,
	0x68, 0x6f, 0x75, 0x72, 0x73, 0x5f, 0x6d, 0x69, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0d, 0x42,
	0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08, 0x2a, 0x06, 0x18, 0x18, 0x28, 0x01, 0x40, 0x01, 0x48,
	0x06, 0x52, 0x08, 0x68, 0x6f, 0x75, 0x72, 0x73, 0x4d, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x30,
	0x0a, 0x09, 0x68, 0x6f, 0x75, 0x72, 0x73, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x0d, 0x42, 0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08, 0x2a, 0x06, 0x18, 0x18, 0x28, 0x01, 0x40,
	0x01, 0x48, 0x07, 0x52, 0x08, 0x68, 0x6f, 0x75, 0x72, 0x73, 0x4d, 0x61, 0x78, 0x88, 0x01, 0x01,
	0x12, 0x3a, 0x0a, 0x0d, 0x70, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x73, 0x5f, 0x6d, 0x69,
	0x6e, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x10, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x0a, 0x2a,
	0x08, 0x18, 0xa0, 0x8d, 0x06, 0x28, 0x01, 0x40, 0x01, 0x48, 0x08, 0x52, 0x0c, 0x70, 0x75, 0x72,
	0x63, 0x68, 0x61, 0x73, 0x65, 0x73, 0x4d, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x3a, 0x0a, 0x0d,
	0x70, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x73, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x10, 0x20,
	0x01, 0x28, 0x0d, 0x42, 0x10, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x0a, 0x2a, 0x08, 0x18, 0xa0, 0x8d,
	0x06, 0x28, 0x01, 0x40, 0x01, 0x48, 0x09, 0x52, 0x0c, 0x70, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73,
	0x65, 0x73, 0x4d, 0x61, 0x78, 0x88, 0x01, 0x01, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x6d, 0x64, 0x42,
	0x06, 0x0a, 0x04, 0x5f, 0x6d, 0x6d, 0x64, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x67, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x63,
	0x6f, 0x64, 0x65, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x69, 0x6e, 0x42,
	0x0a, 0x0a, 0x08, 0x5f, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x61, 0x78, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x68, 0x6f, 0x75, 0x72, 0x73, 0x5f, 0x6d, 0x69, 0x6e, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x68, 0x6f,
	0x75, 0x72, 0x73, 0x5f, 0x6d, 0x61, 0x78, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x70, 0x75, 0x72, 0x63,
	0x68, 0x61, 0x73, 0x65, 0x73, 0x5f, 0x6d, 0x69, 0x6e, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x70, 0x75,
	0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x73, 0x5f, 0x6d, 0x61, 0x78, 0x4a, 0x04, 0x08, 0x01, 0x10,
	0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x4a, 0x04, 0x08,
	0x04, 0x10, 0x05, 0x42, 0x23, 0x5a, 0x21, 0x78, 0x2d, 0x67, 0x77, 0x69, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2f, 0x76, 0x31,
	0x3b, 0x61, 0x73, 0x73, 0x65, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_core_asset_v1_asset_proto_rawDescOnce sync.Once
	file_proto_core_asset_v1_asset_proto_rawDescData = file_proto_core_asset_v1_asset_proto_rawDesc
)

func file_proto_core_asset_v1_asset_proto_rawDescGZIP() []byte {
	file_proto_core_asset_v1_asset_proto_rawDescOnce.Do(func() {
		file_proto_core_asset_v1_asset_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_core_asset_v1_asset_proto_rawDescData)
	})
	return file_proto_core_asset_v1_asset_proto_rawDescData
}

var file_proto_core_asset_v1_asset_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_core_asset_v1_asset_proto_goTypes = []interface{}{
	(*AssetCore)(nil),               // 0: proto.core.asset.v1.AssetCore
	(*AssetChart)(nil),              // 1: proto.core.asset.v1.AssetChart
	(*AssetInsight)(nil),            // 2: proto.core.asset.v1.AssetInsight
	(*AssetAudience)(nil),           // 3: proto.core.asset.v1.AssetAudience
	(*v1.ShareQID)(nil),             // 4: proto.core._share.v1.ShareQID
	(*v1.MetaDescription)(nil),      // 5: proto.core._share.v1.MetaDescription
	(*v1.MetaMultiDescription)(nil), // 6: proto.core._share.v1.MetaMultiDescription
	(*structpb.Struct)(nil),         // 7: google.protobuf.Struct
}
var file_proto_core_asset_v1_asset_proto_depIdxs = []int32{
	4,  // 0: proto.core.asset.v1.AssetCore.qid:type_name -> proto.core._share.v1.ShareQID
	5,  // 1: proto.core.asset.v1.AssetCore.md:type_name -> proto.core._share.v1.MetaDescription
	6,  // 2: proto.core.asset.v1.AssetCore.mmd:type_name -> proto.core._share.v1.MetaMultiDescription
	1,  // 3: proto.core.asset.v1.AssetCore.chart:type_name -> proto.core.asset.v1.AssetChart
	2,  // 4: proto.core.asset.v1.AssetCore.insight:type_name -> proto.core.asset.v1.AssetInsight
	3,  // 5: proto.core.asset.v1.AssetCore.audience:type_name -> proto.core.asset.v1.AssetAudience
	5,  // 6: proto.core.asset.v1.AssetChart.md:type_name -> proto.core._share.v1.MetaDescription
	6,  // 7: proto.core.asset.v1.AssetChart.mmd:type_name -> proto.core._share.v1.MetaMultiDescription
	7,  // 8: proto.core.asset.v1.AssetChart.data:type_name -> google.protobuf.Struct
	7,  // 9: proto.core.asset.v1.AssetChart.options:type_name -> google.protobuf.Struct
	5,  // 10: proto.core.asset.v1.AssetInsight.md:type_name -> proto.core._share.v1.MetaDescription
	6,  // 11: proto.core.asset.v1.AssetInsight.mmd:type_name -> proto.core._share.v1.MetaMultiDescription
	5,  // 12: proto.core.asset.v1.AssetAudience.md:type_name -> proto.core._share.v1.MetaDescription
	6,  // 13: proto.core.asset.v1.AssetAudience.mmd:type_name -> proto.core._share.v1.MetaMultiDescription
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_proto_core_asset_v1_asset_proto_init() }
func file_proto_core_asset_v1_asset_proto_init() {
	if File_proto_core_asset_v1_asset_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_core_asset_v1_asset_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetCore); i {
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
		file_proto_core_asset_v1_asset_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetChart); i {
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
		file_proto_core_asset_v1_asset_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetInsight); i {
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
		file_proto_core_asset_v1_asset_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetAudience); i {
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
	file_proto_core_asset_v1_asset_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_proto_core_asset_v1_asset_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_proto_core_asset_v1_asset_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_proto_core_asset_v1_asset_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_core_asset_v1_asset_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_core_asset_v1_asset_proto_goTypes,
		DependencyIndexes: file_proto_core_asset_v1_asset_proto_depIdxs,
		MessageInfos:      file_proto_core_asset_v1_asset_proto_msgTypes,
	}.Build()
	File_proto_core_asset_v1_asset_proto = out.File
	file_proto_core_asset_v1_asset_proto_rawDesc = nil
	file_proto_core_asset_v1_asset_proto_goTypes = nil
	file_proto_core_asset_v1_asset_proto_depIdxs = nil
}
