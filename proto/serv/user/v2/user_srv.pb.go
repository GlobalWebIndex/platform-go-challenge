// proto/serv/user/v2/user_srv.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: proto/serv/user/v2/user_srv.proto

package user_srvpb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	v11 "x-gwi/proto/core/_store/v1"
	v12 "x-gwi/proto/core/favourite/v1"
	v13 "x-gwi/proto/core/opinion/v1"
	v1 "x-gwi/proto/core/user/v1"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_proto_serv_user_v2_user_srv_proto protoreflect.FileDescriptor

var file_proto_serv_user_v2_user_srv_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x2f, 0x76, 0x32, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x73, 0x72, 0x76, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63,
	0x6f, 0x72, 0x65, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x23, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x6f, 0x70, 0x69,
	0x6e, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x32, 0xba, 0x06, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x6c, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x20,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x1a, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x3a, 0x01, 0x2a, 0x22, 0x13, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x32, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x12, 0x7f, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x58, 0x1a, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x36, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x30, 0x5a, 0x15, 0x3a, 0x01, 0x2a, 0x22, 0x10, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x32, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x74, 0x12, 0x17, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x32, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x74, 0x2f, 0x7b, 0x75, 0x75,
	0x69, 0x64, 0x7d, 0x12, 0x48, 0x0a, 0x04, 0x47, 0x65, 0x74, 0x74, 0x12, 0x1e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x58, 0x1a, 0x20, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x7f, 0x0a,
	0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x1a, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x31, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x2b, 0x3a, 0x01, 0x2a, 0x32, 0x26, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x32, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x2e,
	0x75, 0x75, 0x69, 0x64, 0x7d, 0x2f, 0x7b, 0x69, 0x64, 0x2e, 0x72, 0x65, 0x76, 0x7d, 0x12, 0x6e,
	0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x58, 0x1a, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x22, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x1c, 0x2a, 0x1a, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x32, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x2f, 0x7b, 0x75, 0x75, 0x69, 0x64, 0x7d, 0x12, 0x83,
	0x01, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65,
	0x73, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x5f,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44,
	0x58, 0x1a, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x66,
	0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f,
	0x75, 0x72, 0x69, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x22, 0x26, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x20, 0x12, 0x1e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x32, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x2f, 0x7b, 0x75, 0x75, 0x69, 0x64, 0x7d, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x75, 0x72, 0x69, 0x74,
	0x65, 0x73, 0x30, 0x01, 0x12, 0x7b, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x70, 0x69, 0x6e,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x49, 0x44, 0x58, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x70, 0x69,
	0x6e, 0x69, 0x6f, 0x6e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1e, 0x12, 0x1c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x32, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f,
	0x7b, 0x75, 0x75, 0x69, 0x64, 0x7d, 0x2f, 0x6f, 0x70, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x73, 0x30,
	0x01, 0x42, 0x25, 0x5a, 0x23, 0x78, 0x2d, 0x67, 0x77, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x76, 0x32, 0x3b, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x73, 0x72, 0x76, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_proto_serv_user_v2_user_srv_proto_goTypes = []interface{}{
	(*v1.UserInstance)(nil),    // 0: proto.core.user.v1.UserInstance
	(*v11.StoreIDX)(nil),       // 1: proto.core._store.v1.StoreIDX
	(*v12.FavouriteAsset)(nil), // 2: proto.core.favourite.v1.FavouriteAsset
	(*v13.OpinionAsset)(nil),   // 3: proto.core.opinion.v1.OpinionAsset
}
var file_proto_serv_user_v2_user_srv_proto_depIdxs = []int32{
	0, // 0: proto.serv.user.v2.UserService.Create:input_type -> proto.core.user.v1.UserInstance
	1, // 1: proto.serv.user.v2.UserService.Get:input_type -> proto.core._store.v1.StoreIDX
	1, // 2: proto.serv.user.v2.UserService.Gett:input_type -> proto.core._store.v1.StoreIDX
	0, // 3: proto.serv.user.v2.UserService.Update:input_type -> proto.core.user.v1.UserInstance
	1, // 4: proto.serv.user.v2.UserService.Delete:input_type -> proto.core._store.v1.StoreIDX
	1, // 5: proto.serv.user.v2.UserService.ListFavourites:input_type -> proto.core._store.v1.StoreIDX
	1, // 6: proto.serv.user.v2.UserService.ListOpinions:input_type -> proto.core._store.v1.StoreIDX
	0, // 7: proto.serv.user.v2.UserService.Create:output_type -> proto.core.user.v1.UserInstance
	0, // 8: proto.serv.user.v2.UserService.Get:output_type -> proto.core.user.v1.UserInstance
	0, // 9: proto.serv.user.v2.UserService.Gett:output_type -> proto.core.user.v1.UserInstance
	0, // 10: proto.serv.user.v2.UserService.Update:output_type -> proto.core.user.v1.UserInstance
	0, // 11: proto.serv.user.v2.UserService.Delete:output_type -> proto.core.user.v1.UserInstance
	2, // 12: proto.serv.user.v2.UserService.ListFavourites:output_type -> proto.core.favourite.v1.FavouriteAsset
	3, // 13: proto.serv.user.v2.UserService.ListOpinions:output_type -> proto.core.opinion.v1.OpinionAsset
	7, // [7:14] is the sub-list for method output_type
	0, // [0:7] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_serv_user_v2_user_srv_proto_init() }
func file_proto_serv_user_v2_user_srv_proto_init() {
	if File_proto_serv_user_v2_user_srv_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_serv_user_v2_user_srv_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_serv_user_v2_user_srv_proto_goTypes,
		DependencyIndexes: file_proto_serv_user_v2_user_srv_proto_depIdxs,
	}.Build()
	File_proto_serv_user_v2_user_srv_proto = out.File
	file_proto_serv_user_v2_user_srv_proto_rawDesc = nil
	file_proto_serv_user_v2_user_srv_proto_goTypes = nil
	file_proto_serv_user_v2_user_srv_proto_depIdxs = nil
}
