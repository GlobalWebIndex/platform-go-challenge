// proto/core/favourite/v1/favourite.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: proto/core/favourite/v1/favourite.proto

package favourite_pb

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	v11 "x-gwi/proto/core/asset/v1"
	v1 "x-gwi/proto/core/idx/v1"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// FavouriteAsset
type FavouriteAsset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id - User's Favourite Asset
	Id *v1.IDX `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	// md - User's Favourite Asset MetaData
	Md *v11.AssetMetaData `protobuf:"bytes,4,opt,name=md,proto3,oneof" json:"md,omitempty"`
	// id_user (from)
	IdUser *v1.IDX `protobuf:"bytes,5,opt,name=id_user,json=idUser,proto3,oneof" json:"id_user,omitempty"`
	// id_asset (to)
	IdAsset *v1.IDX `protobuf:"bytes,6,opt,name=id_asset,json=idAsset,proto3,oneof" json:"id_asset,omitempty"`
	// asset - User's Favourite Asset Instance
	Asset *v11.AssetInstance `protobuf:"bytes,7,opt,name=asset,proto3,oneof" json:"asset,omitempty"`
}

func (x *FavouriteAsset) Reset() {
	*x = FavouriteAsset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_core_favourite_v1_favourite_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FavouriteAsset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavouriteAsset) ProtoMessage() {}

func (x *FavouriteAsset) ProtoReflect() protoreflect.Message {
	mi := &file_proto_core_favourite_v1_favourite_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FavouriteAsset.ProtoReflect.Descriptor instead.
func (*FavouriteAsset) Descriptor() ([]byte, []int) {
	return file_proto_core_favourite_v1_favourite_proto_rawDescGZIP(), []int{0}
}

func (x *FavouriteAsset) GetId() *v1.IDX {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *FavouriteAsset) GetMd() *v11.AssetMetaData {
	if x != nil {
		return x.Md
	}
	return nil
}

func (x *FavouriteAsset) GetIdUser() *v1.IDX {
	if x != nil {
		return x.IdUser
	}
	return nil
}

func (x *FavouriteAsset) GetIdAsset() *v1.IDX {
	if x != nil {
		return x.IdAsset
	}
	return nil
}

func (x *FavouriteAsset) GetAsset() *v11.AssetInstance {
	if x != nil {
		return x.Asset
	}
	return nil
}

var File_proto_core_favourite_v1_favourite_proto protoreflect.FileDescriptor

var file_proto_core_favourite_v1_favourite_proto_rawDesc = []byte{
	0x0a, 0x27, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x66, 0x61, 0x76,
	0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72,
	0x69, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e,
	0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x61, 0x73, 0x73, 0x65, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x69, 0x64, 0x78, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x64, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe3, 0x02, 0x0a, 0x0e, 0x46,
	0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x35, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x64, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x44,
	0x58, 0x42, 0x0d, 0xe0, 0x41, 0x01, 0xfa, 0x42, 0x07, 0x8a, 0x01, 0x04, 0x08, 0x00, 0x10, 0x00,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x37, 0x0a, 0x02, 0x6d, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x73,
	0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61,
	0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x02, 0x6d, 0x64, 0x88, 0x01, 0x01, 0x12, 0x34, 0x0a,
	0x07, 0x69, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x64, 0x78, 0x2e,
	0x76, 0x31, 0x2e, 0x49, 0x44, 0x58, 0x48, 0x01, 0x52, 0x06, 0x69, 0x64, 0x55, 0x73, 0x65, 0x72,
	0x88, 0x01, 0x01, 0x12, 0x36, 0x0a, 0x08, 0x69, 0x64, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x69, 0x64, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x44, 0x58, 0x48, 0x02, 0x52,
	0x07, 0x69, 0x64, 0x41, 0x73, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x12, 0x3d, 0x0a, 0x05, 0x61,
	0x73, 0x73, 0x65, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x48, 0x03,
	0x52, 0x05, 0x61, 0x73, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x6d,
	0x64, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x69, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x0b, 0x0a,
	0x09, 0x5f, 0x69, 0x64, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x61,
	0x73, 0x73, 0x65, 0x74, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03,
	0x42, 0x2c, 0x5a, 0x2a, 0x78, 0x2d, 0x67, 0x77, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x76,
	0x31, 0x3b, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_core_favourite_v1_favourite_proto_rawDescOnce sync.Once
	file_proto_core_favourite_v1_favourite_proto_rawDescData = file_proto_core_favourite_v1_favourite_proto_rawDesc
)

func file_proto_core_favourite_v1_favourite_proto_rawDescGZIP() []byte {
	file_proto_core_favourite_v1_favourite_proto_rawDescOnce.Do(func() {
		file_proto_core_favourite_v1_favourite_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_core_favourite_v1_favourite_proto_rawDescData)
	})
	return file_proto_core_favourite_v1_favourite_proto_rawDescData
}

var file_proto_core_favourite_v1_favourite_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_core_favourite_v1_favourite_proto_goTypes = []interface{}{
	(*FavouriteAsset)(nil),    // 0: proto.core.favourite.v1.FavouriteAsset
	(*v1.IDX)(nil),            // 1: proto.core.idx.v1.IDX
	(*v11.AssetMetaData)(nil), // 2: proto.core.asset.v1.AssetMetaData
	(*v11.AssetInstance)(nil), // 3: proto.core.asset.v1.AssetInstance
}
var file_proto_core_favourite_v1_favourite_proto_depIdxs = []int32{
	1, // 0: proto.core.favourite.v1.FavouriteAsset.id:type_name -> proto.core.idx.v1.IDX
	2, // 1: proto.core.favourite.v1.FavouriteAsset.md:type_name -> proto.core.asset.v1.AssetMetaData
	1, // 2: proto.core.favourite.v1.FavouriteAsset.id_user:type_name -> proto.core.idx.v1.IDX
	1, // 3: proto.core.favourite.v1.FavouriteAsset.id_asset:type_name -> proto.core.idx.v1.IDX
	3, // 4: proto.core.favourite.v1.FavouriteAsset.asset:type_name -> proto.core.asset.v1.AssetInstance
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_core_favourite_v1_favourite_proto_init() }
func file_proto_core_favourite_v1_favourite_proto_init() {
	if File_proto_core_favourite_v1_favourite_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_core_favourite_v1_favourite_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FavouriteAsset); i {
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
	file_proto_core_favourite_v1_favourite_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_core_favourite_v1_favourite_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_core_favourite_v1_favourite_proto_goTypes,
		DependencyIndexes: file_proto_core_favourite_v1_favourite_proto_depIdxs,
		MessageInfos:      file_proto_core_favourite_v1_favourite_proto_msgTypes,
	}.Build()
	File_proto_core_favourite_v1_favourite_proto = out.File
	file_proto_core_favourite_v1_favourite_proto_rawDesc = nil
	file_proto_core_favourite_v1_favourite_proto_goTypes = nil
	file_proto_core_favourite_v1_favourite_proto_depIdxs = nil
}
