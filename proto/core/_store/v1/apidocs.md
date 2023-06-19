# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [proto/core/_store/v1/store.proto](#proto_core__store_v1_store-proto)
    - [Store](#proto-core-_store-v1-Store)
    - [StoreAQL](#proto-core-_store-v1-StoreAQL)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proto_core__store_v1_store-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/core/_store/v1/store.proto



<a name="proto-core-_store-v1-Store"></a>

### Store
Store


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qid | [proto.core._share.v1.ShareQID](#proto-core-_share-v1-ShareQID) |  | qid - StoreQID |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) | optional | user - UserCore |
| asset | [proto.core.asset.v1.AssetCore](#proto-core-asset-v1-AssetCore) | optional | asset - AssetCore |
| favourite | [proto.core.favourite.v1.FavouriteCore](#proto-core-favourite-v1-FavouriteCore) | optional | favourite - FavouriteCore |
| opinion | [proto.core.opinion.v1.OpinionCore](#proto-core-opinion-v1-OpinionCore) | optional | opinion - OpinionCore |






<a name="proto-core-_store-v1-StoreAQL"></a>

### StoreAQL
StoreAQL store Arango
DocumentMeta contains all meta data used to identifier a document.
type DocumentMeta struct {
	Key    string     `json:&#34;_key,omitempty&#34;`
	ID     DocumentID `json:&#34;_id,omitempty&#34;`
	Rev    string     `json:&#34;_rev,omitempty&#34;`
	OldRev string     `json:&#34;_oldRev,omitempty&#34;`
}
DocumentID references a document in a collection.
Format: collection/_key
type DocumentID string
ArangoID is a generic Arango ID struct representation
type ArangoID struct {
	ID               string `json:&#34;qid,omitempty&#34;`
	GloballyUniqueId string `json:&#34;globallyUniqueId,omitempty&#34;`
}
REV rev = 4 [json_name = &#34;_rev&#34;];

WARNING! generated pb files keep json_name for protobuf but not for json
want: json=_key,proto3&#34; json:&#34;_key (from json_name = &#34;_key&#34;)
Key string `protobuf:&#34;bytes,4,opt,name=key,json=_key,proto3&#34; json:&#34;_key,omitempty&#34;`
got:  json=_key,proto3&#34; json:&#34;key  (from json_name = &#34;_key&#34;)
Key string `protobuf:&#34;bytes,4,opt,name=key,json=_key,proto3&#34; json:&#34;key,omitempty&#34;`
solution: 1) edit generated proto/core/_store/v1/store.pb.go; 
or 2) synch dedicated type in app/storage/storepb2/storepb2.go


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | key - AQL _key - doc storage key unique per db.collection |
| id | [string](#string) |  | id - AQL _id - DocumentID string Format: collection/_key |
| rev | [string](#string) |  | rev - AQL _rev - revision string |
| old_rev | [string](#string) |  | old_rev - AQL _oldRev old revision string |
| qid | [proto.core._share.v1.ShareQID](#proto-core-_share-v1-ShareQID) |  | qid - StoreQID |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) | optional | user - UserCore |
| asset | [proto.core.asset.v1.AssetCore](#proto-core-asset-v1-AssetCore) | optional | asset - AssetCore |
| favourite | [proto.core.favourite.v1.FavouriteCore](#proto-core-favourite-v1-FavouriteCore) | optional | favourite - FavouriteCore |
| opinion | [proto.core.opinion.v1.OpinionCore](#proto-core-opinion-v1-OpinionCore) | optional | opinion - OpinionCore |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

