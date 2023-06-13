// proto/serv/favourite/v1/favourite_srv.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: proto/serv/favourite/v1/favourite_srv.proto

package favourite_srvpb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	reflect "reflect"
	sync "sync"
	v1 "x-gwi/proto/core/favourite/v1"
	v11 "x-gwi/proto/core/idx/v1"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// CreateRequest
type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// favourite
	Favourite *v1.FavouriteAsset `protobuf:"bytes,4,opt,name=favourite,proto3" json:"favourite,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{0}
}

func (x *CreateRequest) GetFavourite() *v1.FavouriteAsset {
	if x != nil {
		return x.Favourite
	}
	return nil
}

// CreateResponse
type CreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// favourite
	Favourite *v1.FavouriteAsset `protobuf:"bytes,4,opt,name=favourite,proto3" json:"favourite,omitempty"`
}

func (x *CreateResponse) Reset() {
	*x = CreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateResponse) ProtoMessage() {}

func (x *CreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateResponse.ProtoReflect.Descriptor instead.
func (*CreateResponse) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{1}
}

func (x *CreateResponse) GetFavourite() *v1.FavouriteAsset {
	if x != nil {
		return x.Favourite
	}
	return nil
}

// GetRequest
type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id
	Id *v11.IDX `protobuf:"bytes,4,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{2}
}

func (x *GetRequest) GetId() *v11.IDX {
	if x != nil {
		return x.Id
	}
	return nil
}

// GetResponse
type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// favourite
	Favourite *v1.FavouriteAsset `protobuf:"bytes,4,opt,name=favourite,proto3" json:"favourite,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{3}
}

func (x *GetResponse) GetFavourite() *v1.FavouriteAsset {
	if x != nil {
		return x.Favourite
	}
	return nil
}

// UpdateRequest
type UpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id
	Id *v11.IDX `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	// favourite
	Favourite *v1.FavouriteAsset `protobuf:"bytes,4,opt,name=favourite,proto3" json:"favourite,omitempty"`
	// update_mask
	// https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/patch_feature
	UpdateMask *fieldmaskpb.FieldMask `protobuf:"bytes,5,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
}

func (x *UpdateRequest) Reset() {
	*x = UpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRequest) ProtoMessage() {}

func (x *UpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRequest.ProtoReflect.Descriptor instead.
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateRequest) GetId() *v11.IDX {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *UpdateRequest) GetFavourite() *v1.FavouriteAsset {
	if x != nil {
		return x.Favourite
	}
	return nil
}

func (x *UpdateRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

// UpdateResponse
type UpdateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// favourite
	Favourite *v1.FavouriteAsset `protobuf:"bytes,4,opt,name=favourite,proto3" json:"favourite,omitempty"`
}

func (x *UpdateResponse) Reset() {
	*x = UpdateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateResponse) ProtoMessage() {}

func (x *UpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateResponse.ProtoReflect.Descriptor instead.
func (*UpdateResponse) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateResponse) GetFavourite() *v1.FavouriteAsset {
	if x != nil {
		return x.Favourite
	}
	return nil
}

// DeleteRequest
type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id
	Id *v11.IDX `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteRequest) GetId() *v11.IDX {
	if x != nil {
		return x.Id
	}
	return nil
}

// DeleteResponse
type DeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// favourite
	Favourite *v1.FavouriteAsset `protobuf:"bytes,4,opt,name=favourite,proto3" json:"favourite,omitempty"`
}

func (x *DeleteResponse) Reset() {
	*x = DeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteResponse) ProtoMessage() {}

func (x *DeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteResponse.ProtoReflect.Descriptor instead.
func (*DeleteResponse) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteResponse) GetFavourite() *v1.FavouriteAsset {
	if x != nil {
		return x.Favourite
	}
	return nil
}

// ListRequest
type ListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id_user (from)
	IdUser *v11.IDX `protobuf:"bytes,5,opt,name=id_user,json=idUser,proto3" json:"id_user,omitempty"`
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRequest.ProtoReflect.Descriptor instead.
func (*ListRequest) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{8}
}

func (x *ListRequest) GetIdUser() *v11.IDX {
	if x != nil {
		return x.IdUser
	}
	return nil
}

// ListResponse
type ListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// favourite
	Favourite *v1.FavouriteAsset `protobuf:"bytes,4,opt,name=favourite,proto3" json:"favourite,omitempty"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResponse.ProtoReflect.Descriptor instead.
func (*ListResponse) Descriptor() ([]byte, []int) {
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP(), []int{9}
}

func (x *ListResponse) GetFavourite() *v1.FavouriteAsset {
	if x != nil {
		return x.Favourite
	}
	return nil
}

var File_proto_serv_favourite_v1_favourite_srv_proto protoreflect.FileDescriptor

var file_proto_serv_favourite_v1_favourite_srv_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x2f, 0x66, 0x61, 0x76,
	0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72,
	0x69, 0x74, 0x65, 0x5f, 0x73, 0x72, 0x76, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72,
	0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f,
	0x72, 0x65, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x2f,
	0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x69, 0x64, 0x78, 0x2f,
	0x76, 0x31, 0x2f, 0x69, 0x64, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x68, 0x0a, 0x0d,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x45, 0x0a,
	0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x66, 0x61,
	0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x75,
	0x72, 0x69, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x75,
	0x72, 0x69, 0x74, 0x65, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03,
	0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x22, 0x69, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x09, 0x66, 0x61, 0x76, 0x6f,
	0x75, 0x72, 0x69, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69,
	0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x41,
	0x73, 0x73, 0x65, 0x74, 0x52, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x4a,
	0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x03, 0x10,
	0x04, 0x22, 0x46, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x26, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x64, 0x78, 0x2e, 0x76, 0x31, 0x2e,
	0x49, 0x44, 0x58, 0x52, 0x02, 0x69, 0x64, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08,
	0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x22, 0x66, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x09, 0x66, 0x61, 0x76, 0x6f,
	0x75, 0x72, 0x69, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69,
	0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x41,
	0x73, 0x73, 0x65, 0x74, 0x52, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x4a,
	0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x03, 0x10,
	0x04, 0x22, 0xc7, 0x01, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x64, 0x78,
	0x2e, 0x76, 0x31, 0x2e, 0x49, 0x44, 0x58, 0x52, 0x02, 0x69, 0x64, 0x12, 0x45, 0x0a, 0x09, 0x66,
	0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x66, 0x61, 0x76, 0x6f,
	0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69,
	0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69,
	0x74, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x6d, 0x61, 0x73,
	0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d,
	0x61, 0x73, 0x6b, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x73, 0x6b, 0x4a,
	0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x22, 0x69, 0x0a, 0x0e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a,
	0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x66, 0x61,
	0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x75,
	0x72, 0x69, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x75,
	0x72, 0x69, 0x74, 0x65, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03,
	0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x22, 0x43, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x69, 0x64, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x44, 0x58, 0x52, 0x02, 0x69, 0x64, 0x4a,
	0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x22, 0x69, 0x0a, 0x0e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a,
	0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x66, 0x61,
	0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x75,
	0x72, 0x69, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x75,
	0x72, 0x69, 0x74, 0x65, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03,
	0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x22, 0x50, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x07, 0x69, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x69, 0x64, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x44, 0x58, 0x52, 0x06,
	0x69, 0x64, 0x55, 0x73, 0x65, 0x72, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02,
	0x10, 0x03, 0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x22, 0x67, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x09, 0x66, 0x61, 0x76, 0x6f,
	0x75, 0x72, 0x69, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69,
	0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x41,
	0x73, 0x73, 0x65, 0x74, 0x52, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x4a,
	0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x03, 0x10,
	0x04, 0x32, 0xb7, 0x06, 0x0a, 0x10, 0x46, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x86, 0x01, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x12, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e, 0x66,
	0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x2b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x3a, 0x09, 0x66, 0x61, 0x76, 0x6f,
	0x75, 0x72, 0x69, 0x74, 0x65, 0x22, 0x18, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x66,
	0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12,
	0x96, 0x01, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72,
	0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x44, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x3e, 0x5a, 0x1b, 0x3a, 0x02, 0x69, 0x64,
	0x22, 0x15, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72,
	0x69, 0x74, 0x65, 0x2f, 0x67, 0x65, 0x74, 0x12, 0x1f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x67, 0x65, 0x74, 0x2f, 0x7b,
	0x69, 0x64, 0x2e, 0x75, 0x75, 0x69, 0x64, 0x7d, 0x12, 0x51, 0x0a, 0x04, 0x47, 0x65, 0x74, 0x74,
	0x12, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e, 0x66, 0x61,
	0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x99, 0x01, 0x0a, 0x06,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f,
	0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x38, 0x3a,
	0x09, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x32, 0x2b, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x2e, 0x75, 0x75, 0x69, 0x64, 0x7d, 0x2f, 0x7b,
	0x69, 0x64, 0x2e, 0x72, 0x65, 0x76, 0x7d, 0x12, 0x85, 0x01, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e,
	0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x2a, 0x22, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x2e, 0x75, 0x75, 0x69, 0x64, 0x7d, 0x12,
	0x89, 0x01, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e, 0x66, 0x61, 0x76, 0x6f,
	0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2c, 0x12, 0x2a, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65,
	0x2f, 0x6c, 0x69, 0x73, 0x74, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x7b, 0x69, 0x64, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x75, 0x75, 0x69, 0x64, 0x7d, 0x30, 0x01, 0x42, 0x2f, 0x5a, 0x2d, 0x78,
	0x2d, 0x67, 0x77, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x2f,
	0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x66, 0x61, 0x76,
	0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x73, 0x72, 0x76, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_serv_favourite_v1_favourite_srv_proto_rawDescOnce sync.Once
	file_proto_serv_favourite_v1_favourite_srv_proto_rawDescData = file_proto_serv_favourite_v1_favourite_srv_proto_rawDesc
)

func file_proto_serv_favourite_v1_favourite_srv_proto_rawDescGZIP() []byte {
	file_proto_serv_favourite_v1_favourite_srv_proto_rawDescOnce.Do(func() {
		file_proto_serv_favourite_v1_favourite_srv_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_serv_favourite_v1_favourite_srv_proto_rawDescData)
	})
	return file_proto_serv_favourite_v1_favourite_srv_proto_rawDescData
}

var file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_proto_serv_favourite_v1_favourite_srv_proto_goTypes = []interface{}{
	(*CreateRequest)(nil),         // 0: proto.serv.favourite.v1.CreateRequest
	(*CreateResponse)(nil),        // 1: proto.serv.favourite.v1.CreateResponse
	(*GetRequest)(nil),            // 2: proto.serv.favourite.v1.GetRequest
	(*GetResponse)(nil),           // 3: proto.serv.favourite.v1.GetResponse
	(*UpdateRequest)(nil),         // 4: proto.serv.favourite.v1.UpdateRequest
	(*UpdateResponse)(nil),        // 5: proto.serv.favourite.v1.UpdateResponse
	(*DeleteRequest)(nil),         // 6: proto.serv.favourite.v1.DeleteRequest
	(*DeleteResponse)(nil),        // 7: proto.serv.favourite.v1.DeleteResponse
	(*ListRequest)(nil),           // 8: proto.serv.favourite.v1.ListRequest
	(*ListResponse)(nil),          // 9: proto.serv.favourite.v1.ListResponse
	(*v1.FavouriteAsset)(nil),     // 10: proto.core.favourite.v1.FavouriteAsset
	(*v11.IDX)(nil),               // 11: proto.core.idx.v1.IDX
	(*fieldmaskpb.FieldMask)(nil), // 12: google.protobuf.FieldMask
}
var file_proto_serv_favourite_v1_favourite_srv_proto_depIdxs = []int32{
	10, // 0: proto.serv.favourite.v1.CreateRequest.favourite:type_name -> proto.core.favourite.v1.FavouriteAsset
	10, // 1: proto.serv.favourite.v1.CreateResponse.favourite:type_name -> proto.core.favourite.v1.FavouriteAsset
	11, // 2: proto.serv.favourite.v1.GetRequest.id:type_name -> proto.core.idx.v1.IDX
	10, // 3: proto.serv.favourite.v1.GetResponse.favourite:type_name -> proto.core.favourite.v1.FavouriteAsset
	11, // 4: proto.serv.favourite.v1.UpdateRequest.id:type_name -> proto.core.idx.v1.IDX
	10, // 5: proto.serv.favourite.v1.UpdateRequest.favourite:type_name -> proto.core.favourite.v1.FavouriteAsset
	12, // 6: proto.serv.favourite.v1.UpdateRequest.update_mask:type_name -> google.protobuf.FieldMask
	10, // 7: proto.serv.favourite.v1.UpdateResponse.favourite:type_name -> proto.core.favourite.v1.FavouriteAsset
	11, // 8: proto.serv.favourite.v1.DeleteRequest.id:type_name -> proto.core.idx.v1.IDX
	10, // 9: proto.serv.favourite.v1.DeleteResponse.favourite:type_name -> proto.core.favourite.v1.FavouriteAsset
	11, // 10: proto.serv.favourite.v1.ListRequest.id_user:type_name -> proto.core.idx.v1.IDX
	10, // 11: proto.serv.favourite.v1.ListResponse.favourite:type_name -> proto.core.favourite.v1.FavouriteAsset
	0,  // 12: proto.serv.favourite.v1.FavouriteService.Create:input_type -> proto.serv.favourite.v1.CreateRequest
	2,  // 13: proto.serv.favourite.v1.FavouriteService.Get:input_type -> proto.serv.favourite.v1.GetRequest
	2,  // 14: proto.serv.favourite.v1.FavouriteService.Gett:input_type -> proto.serv.favourite.v1.GetRequest
	4,  // 15: proto.serv.favourite.v1.FavouriteService.Update:input_type -> proto.serv.favourite.v1.UpdateRequest
	6,  // 16: proto.serv.favourite.v1.FavouriteService.Delete:input_type -> proto.serv.favourite.v1.DeleteRequest
	8,  // 17: proto.serv.favourite.v1.FavouriteService.List:input_type -> proto.serv.favourite.v1.ListRequest
	1,  // 18: proto.serv.favourite.v1.FavouriteService.Create:output_type -> proto.serv.favourite.v1.CreateResponse
	3,  // 19: proto.serv.favourite.v1.FavouriteService.Get:output_type -> proto.serv.favourite.v1.GetResponse
	3,  // 20: proto.serv.favourite.v1.FavouriteService.Gett:output_type -> proto.serv.favourite.v1.GetResponse
	5,  // 21: proto.serv.favourite.v1.FavouriteService.Update:output_type -> proto.serv.favourite.v1.UpdateResponse
	7,  // 22: proto.serv.favourite.v1.FavouriteService.Delete:output_type -> proto.serv.favourite.v1.DeleteResponse
	9,  // 23: proto.serv.favourite.v1.FavouriteService.List:output_type -> proto.serv.favourite.v1.ListResponse
	18, // [18:24] is the sub-list for method output_type
	12, // [12:18] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_proto_serv_favourite_v1_favourite_srv_proto_init() }
func file_proto_serv_favourite_v1_favourite_srv_proto_init() {
	if File_proto_serv_favourite_v1_favourite_srv_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
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
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateResponse); i {
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
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
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
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
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
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRequest); i {
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
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateResponse); i {
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
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteResponse); i {
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
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRequest); i {
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
		file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListResponse); i {
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
			RawDescriptor: file_proto_serv_favourite_v1_favourite_srv_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_serv_favourite_v1_favourite_srv_proto_goTypes,
		DependencyIndexes: file_proto_serv_favourite_v1_favourite_srv_proto_depIdxs,
		MessageInfos:      file_proto_serv_favourite_v1_favourite_srv_proto_msgTypes,
	}.Build()
	File_proto_serv_favourite_v1_favourite_srv_proto = out.File
	file_proto_serv_favourite_v1_favourite_srv_proto_rawDesc = nil
	file_proto_serv_favourite_v1_favourite_srv_proto_goTypes = nil
	file_proto_serv_favourite_v1_favourite_srv_proto_depIdxs = nil
}
