package storepb2

import (
	share "x-gwi/proto/core/_share/v1"
	userpb "x-gwi/proto/core/_user/v1"
	assetpb "x-gwi/proto/core/asset/v1"
	favouritepb "x-gwi/proto/core/favourite/v1"
	opinionpb "x-gwi/proto/core/opinion/v1"
)

// .WithRevision(

// StoreAQL store Arango (synch 2)
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
//
//nolint:tagliatelle,revive,stylecheck,govet
type StoreAQL struct {
	// key - AQL _key - doc storage key unique per db.collection
	Key string `json:"_key,omitempty"`
	// id - AQL _id - DocumentID string Format: collection/_key
	Id string `json:"_id,omitempty"`
	// rev - AQL _rev - revision string
	Rev string `json:"_rev,omitempty"`
	// old_rev - AQL _oldRev old revision string
	OldRev string `json:"_old_rev,omitempty"`
	// qid - StoreQID
	Qid *share.ShareQID `json:"qid,omitempty"`
	// user - UserCore
	User *userpb.UserCore `json:"user,omitempty"`
	// asset - AssetCore
	Asset *assetpb.AssetCore `json:"asset,omitempty"`
	// favourite - FavouriteCore
	Favourite *favouritepb.FavouriteCore `json:"favourite,omitempty"`
	// opinion - OpinionCore
	Opinion *opinionpb.OpinionCore `json:"opinion,omitempty"`
}
