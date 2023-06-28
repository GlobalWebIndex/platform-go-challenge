// proto/core/_store/v1/store.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.2
// source: proto/core/_store/v1/store.proto

package storepb

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	v1 "x-gwi/proto/core/_share/v1"
	v11 "x-gwi/proto/core/_user/v1"
	v12 "x-gwi/proto/core/asset/v1"
	v13 "x-gwi/proto/core/favourite/v1"
	v14 "x-gwi/proto/core/opinion/v1"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Store
type Store struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// qid - StoreQID
	Qid *v1.ShareQID `protobuf:"bytes,3,opt,name=qid,proto3" json:"qid,omitempty"`
	// user - UserCore
	User *v11.UserCore `protobuf:"bytes,10,opt,name=user,proto3,oneof" json:"user,omitempty"`
	// asset - AssetCore
	Asset *v12.AssetCore `protobuf:"bytes,11,opt,name=asset,proto3,oneof" json:"asset,omitempty"`
	// favourite - FavouriteCore
	Favourite *v13.FavouriteCore `protobuf:"bytes,12,opt,name=favourite,proto3,oneof" json:"favourite,omitempty"`
	// opinion - OpinionCore
	Opinion *v14.OpinionCore `protobuf:"bytes,13,opt,name=opinion,proto3,oneof" json:"opinion,omitempty"`
}

func (x *Store) Reset() {
	*x = Store{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_core__store_v1_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Store) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Store) ProtoMessage() {}

func (x *Store) ProtoReflect() protoreflect.Message {
	mi := &file_proto_core__store_v1_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Store.ProtoReflect.Descriptor instead.
func (*Store) Descriptor() ([]byte, []int) {
	return file_proto_core__store_v1_store_proto_rawDescGZIP(), []int{0}
}

func (x *Store) GetQid() *v1.ShareQID {
	if x != nil {
		return x.Qid
	}
	return nil
}

func (x *Store) GetUser() *v11.UserCore {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *Store) GetAsset() *v12.AssetCore {
	if x != nil {
		return x.Asset
	}
	return nil
}

func (x *Store) GetFavourite() *v13.FavouriteCore {
	if x != nil {
		return x.Favourite
	}
	return nil
}

func (x *Store) GetOpinion() *v14.OpinionCore {
	if x != nil {
		return x.Opinion
	}
	return nil
}

// StoreAQL store Arango
// DocumentMeta contains all meta data used to identifier a document.
//
//	type DocumentMeta struct {
//		Key    string     `json:"_key,omitempty"`
//		ID     DocumentID `json:"_id,omitempty"`
//		Rev    string     `json:"_rev,omitempty"`
//		OldRev string     `json:"_oldRev,omitempty"`
//	}
//
// DocumentID references a document in a collection.
// Format: collection/_key
// type DocumentID string
// ArangoID is a generic Arango ID struct representation
//
//	type ArangoID struct {
//		ID               string `json:"qid,omitempty"`
//		GloballyUniqueId string `json:"globallyUniqueId,omitempty"`
//	}
//
// REV rev = 4 [json_name = "_rev"];
//
// WARNING! generated pb files keep json_name for protobuf but not for json
// want: json=_key,proto3" json:"_key (from json_name = "_key")
// Key string `protobuf:"bytes,4,opt,name=key,json=_key,proto3" json:"_key,omitempty"`
// got:  json=_key,proto3" json:"key  (from json_name = "_key")
// Key string `protobuf:"bytes,4,opt,name=key,json=_key,proto3" json:"key,omitempty"`
// solution: 1) edit generated proto/core/_store/v1/store.pb.go;
// or 2) synch dedicated type in app/storage/storepb2/storepb2.go
type StoreAQL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// _key - AQL _key - doc storage key unique per db.collection
	XKey string `protobuf:"bytes,4,opt,name=_key,proto3" json:"_key,omitempty"`
	// _id - AQL _id - DocumentID string Format: collection/_key
	XId string `protobuf:"bytes,5,opt,name=_id,proto3" json:"_id,omitempty"`
	// _rev - AQL _rev - revision string
	XRev string `protobuf:"bytes,6,opt,name=_rev,proto3" json:"_rev,omitempty"`
	// _oldRev - AQL _oldRev old revision string
	XOldRev string `protobuf:"bytes,7,opt,name=_oldRev,proto3" json:"_oldRev,omitempty"`
	// _from - AQL _from - edge storage
	XFrom string `protobuf:"bytes,8,opt,name=_from,proto3" json:"_from,omitempty"`
	// _to - AQL _to - edge storage
	XTo string `protobuf:"bytes,9,opt,name=_to,proto3" json:"_to,omitempty"`
	// qid - StoreQID
	Qid *v1.ShareQID `protobuf:"bytes,3,opt,name=qid,proto3" json:"qid,omitempty"`
	// user - UserCore
	User *v11.UserCore `protobuf:"bytes,10,opt,name=user,proto3,oneof" json:"user,omitempty"`
	// asset - AssetCore
	Asset *v12.AssetCore `protobuf:"bytes,11,opt,name=asset,proto3,oneof" json:"asset,omitempty"`
	// favourite - FavouriteCore
	Favourite *v13.FavouriteCore `protobuf:"bytes,12,opt,name=favourite,proto3,oneof" json:"favourite,omitempty"`
	// opinion - OpinionCore
	Opinion *v14.OpinionCore `protobuf:"bytes,13,opt,name=opinion,proto3,oneof" json:"opinion,omitempty"`
}

func (x *StoreAQL) Reset() {
	*x = StoreAQL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_core__store_v1_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreAQL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreAQL) ProtoMessage() {}

func (x *StoreAQL) ProtoReflect() protoreflect.Message {
	mi := &file_proto_core__store_v1_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreAQL.ProtoReflect.Descriptor instead.
func (*StoreAQL) Descriptor() ([]byte, []int) {
	return file_proto_core__store_v1_store_proto_rawDescGZIP(), []int{1}
}

func (x *StoreAQL) GetXKey() string {
	if x != nil {
		return x.XKey
	}
	return ""
}

func (x *StoreAQL) GetXId() string {
	if x != nil {
		return x.XId
	}
	return ""
}

func (x *StoreAQL) GetXRev() string {
	if x != nil {
		return x.XRev
	}
	return ""
}

func (x *StoreAQL) GetXOldRev() string {
	if x != nil {
		return x.XOldRev
	}
	return ""
}

func (x *StoreAQL) GetXFrom() string {
	if x != nil {
		return x.XFrom
	}
	return ""
}

func (x *StoreAQL) GetXTo() string {
	if x != nil {
		return x.XTo
	}
	return ""
}

func (x *StoreAQL) GetQid() *v1.ShareQID {
	if x != nil {
		return x.Qid
	}
	return nil
}

func (x *StoreAQL) GetUser() *v11.UserCore {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *StoreAQL) GetAsset() *v12.AssetCore {
	if x != nil {
		return x.Asset
	}
	return nil
}

func (x *StoreAQL) GetFavourite() *v13.FavouriteCore {
	if x != nil {
		return x.Favourite
	}
	return nil
}

func (x *StoreAQL) GetOpinion() *v14.OpinionCore {
	if x != nil {
		return x.Opinion
	}
	return nil
}

var File_proto_core__store_v1_store_proto protoreflect.FileDescriptor

var file_proto_core__store_v1_store_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x5f, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76,
	0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x76, 0x31,
	0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69,
	0x74, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x70, 0x69,
	0x6e, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xe2, 0x03, 0x0a, 0x05, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x3f, 0x0a,
	0x03, 0x71, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x51, 0x49, 0x44, 0x42, 0x0d, 0xe0, 0x41, 0x01, 0xfa,
	0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00, 0x52, 0x03, 0x71, 0x69, 0x64, 0x12, 0x45,
	0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x72, 0x65, 0x42, 0x0d, 0xe0, 0x41, 0x01,
	0xfa, 0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00, 0x48, 0x00, 0x52, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x48, 0x0a, 0x05, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74,
	0x43, 0x6f, 0x72, 0x65, 0x42, 0x0d, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x07, 0x8a, 0x01, 0x04, 0x08,
	0x00, 0x10, 0x00, 0x48, 0x01, 0x52, 0x05, 0x61, 0x73, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x12,
	0x58, 0x0a, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76,
	0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x72, 0x65, 0x42, 0x0d, 0xe0, 0x41, 0x01, 0xfa,
	0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00, 0x48, 0x02, 0x52, 0x09, 0x66, 0x61, 0x76,
	0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x50, 0x0a, 0x07, 0x6f, 0x70, 0x69,
	0x6e, 0x69, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x4f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x72, 0x65, 0x42, 0x0d,
	0xe0, 0x41, 0x01, 0xfa, 0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00, 0x48, 0x03, 0x52,
	0x07, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x42, 0x0c,
	0x0a, 0x0a, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x42, 0x0a, 0x0a, 0x08,
	0x5f, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04,
	0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x04, 0x10, 0x05, 0x4a, 0x04, 0x08, 0x05, 0x10, 0x06,
	0x4a, 0x04, 0x08, 0x06, 0x10, 0x07, 0x4a, 0x04, 0x08, 0x07, 0x10, 0x08, 0x4a, 0x04, 0x08, 0x08,
	0x10, 0x09, 0x4a, 0x04, 0x08, 0x09, 0x10, 0x0a, 0x22, 0x9d, 0x05, 0x0a, 0x08, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x41, 0x51, 0x4c, 0x12, 0x22, 0x0a, 0x04, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08, 0x72, 0x06, 0x18, 0x80, 0x01,
	0xd0, 0x01, 0x01, 0x52, 0x04, 0x5f, 0x6b, 0x65, 0x79, 0x12, 0x20, 0x0a, 0x03, 0x5f, 0x69, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08, 0x72, 0x06,
	0x18, 0x80, 0x01, 0xd0, 0x01, 0x01, 0x52, 0x03, 0x5f, 0x69, 0x64, 0x12, 0x22, 0x0a, 0x04, 0x5f,
	0x72, 0x65, 0x76, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42,
	0x08, 0x72, 0x06, 0x18, 0x80, 0x01, 0xd0, 0x01, 0x01, 0x52, 0x04, 0x5f, 0x72, 0x65, 0x76, 0x12,
	0x28, 0x0a, 0x07, 0x5f, 0x6f, 0x6c, 0x64, 0x52, 0x65, 0x76, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08, 0x72, 0x06, 0x18, 0x80, 0x01, 0xd0, 0x01, 0x01,
	0x52, 0x07, 0x5f, 0x6f, 0x6c, 0x64, 0x52, 0x65, 0x76, 0x12, 0x24, 0x0a, 0x05, 0x5f, 0x66, 0x72,
	0x6f, 0x6d, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x08,
	0x72, 0x06, 0x18, 0x80, 0x01, 0xd0, 0x01, 0x01, 0x52, 0x05, 0x5f, 0x66, 0x72, 0x6f, 0x6d, 0x12,
	0x20, 0x0a, 0x03, 0x5f, 0x74, 0x6f, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xe0, 0x41,
	0x01, 0xfa, 0x42, 0x08, 0x72, 0x06, 0x18, 0x80, 0x01, 0xd0, 0x01, 0x01, 0x52, 0x03, 0x5f, 0x74,
	0x6f, 0x12, 0x3f, 0x0a, 0x03, 0x71, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x51, 0x49, 0x44, 0x42, 0x0d,
	0xe0, 0x41, 0x01, 0xfa, 0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00, 0x52, 0x03, 0x71,
	0x69, 0x64, 0x12, 0x45, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x72, 0x65, 0x42,
	0x0d, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00, 0x48, 0x00,
	0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x48, 0x0a, 0x05, 0x61, 0x73, 0x73,
	0x65, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x73, 0x73, 0x65, 0x74, 0x43, 0x6f, 0x72, 0x65, 0x42, 0x0d, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x07,
	0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00, 0x48, 0x01, 0x52, 0x05, 0x61, 0x73, 0x73, 0x65, 0x74,
	0x88, 0x01, 0x01, 0x12, 0x58, 0x0a, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x46, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x72, 0x65, 0x42, 0x0d,
	0xe0, 0x41, 0x01, 0xfa, 0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00, 0x48, 0x02, 0x52,
	0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x50, 0x0a,
	0x07, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x6f, 0x70, 0x69, 0x6e,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x43, 0x6f,
	0x72, 0x65, 0x42, 0x0d, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10,
	0x00, 0x48, 0x03, 0x52, 0x07, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42,
	0x07, 0x0a, 0x05, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x61, 0x73, 0x73,
	0x65, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65,
	0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x4a, 0x04, 0x08, 0x01,
	0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x42, 0x24, 0x5a, 0x22, 0x78, 0x2d, 0x67, 0x77,
	0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x5f, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_core__store_v1_store_proto_rawDescOnce sync.Once
	file_proto_core__store_v1_store_proto_rawDescData = file_proto_core__store_v1_store_proto_rawDesc
)

func file_proto_core__store_v1_store_proto_rawDescGZIP() []byte {
	file_proto_core__store_v1_store_proto_rawDescOnce.Do(func() {
		file_proto_core__store_v1_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_core__store_v1_store_proto_rawDescData)
	})
	return file_proto_core__store_v1_store_proto_rawDescData
}

var file_proto_core__store_v1_store_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_core__store_v1_store_proto_goTypes = []interface{}{
	(*Store)(nil),             // 0: proto.core._store.v1.Store
	(*StoreAQL)(nil),          // 1: proto.core._store.v1.StoreAQL
	(*v1.ShareQID)(nil),       // 2: proto.core._share.v1.ShareQID
	(*v11.UserCore)(nil),      // 3: proto.core._user.v1.UserCore
	(*v12.AssetCore)(nil),     // 4: proto.core.asset.v1.AssetCore
	(*v13.FavouriteCore)(nil), // 5: proto.core.favourite.v1.FavouriteCore
	(*v14.OpinionCore)(nil),   // 6: proto.core.opinion.v1.OpinionCore
}
var file_proto_core__store_v1_store_proto_depIdxs = []int32{
	2,  // 0: proto.core._store.v1.Store.qid:type_name -> proto.core._share.v1.ShareQID
	3,  // 1: proto.core._store.v1.Store.user:type_name -> proto.core._user.v1.UserCore
	4,  // 2: proto.core._store.v1.Store.asset:type_name -> proto.core.asset.v1.AssetCore
	5,  // 3: proto.core._store.v1.Store.favourite:type_name -> proto.core.favourite.v1.FavouriteCore
	6,  // 4: proto.core._store.v1.Store.opinion:type_name -> proto.core.opinion.v1.OpinionCore
	2,  // 5: proto.core._store.v1.StoreAQL.qid:type_name -> proto.core._share.v1.ShareQID
	3,  // 6: proto.core._store.v1.StoreAQL.user:type_name -> proto.core._user.v1.UserCore
	4,  // 7: proto.core._store.v1.StoreAQL.asset:type_name -> proto.core.asset.v1.AssetCore
	5,  // 8: proto.core._store.v1.StoreAQL.favourite:type_name -> proto.core.favourite.v1.FavouriteCore
	6,  // 9: proto.core._store.v1.StoreAQL.opinion:type_name -> proto.core.opinion.v1.OpinionCore
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_proto_core__store_v1_store_proto_init() }
func file_proto_core__store_v1_store_proto_init() {
	if File_proto_core__store_v1_store_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_core__store_v1_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Store); i {
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
		file_proto_core__store_v1_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreAQL); i {
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
	file_proto_core__store_v1_store_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_proto_core__store_v1_store_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_core__store_v1_store_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_core__store_v1_store_proto_goTypes,
		DependencyIndexes: file_proto_core__store_v1_store_proto_depIdxs,
		MessageInfos:      file_proto_core__store_v1_store_proto_msgTypes,
	}.Build()
	File_proto_core__store_v1_store_proto = out.File
	file_proto_core__store_v1_store_proto_rawDesc = nil
	file_proto_core__store_v1_store_proto_goTypes = nil
	file_proto_core__store_v1_store_proto_depIdxs = nil
}
