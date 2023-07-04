# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [proto/core/opinion/v1/opinion.proto](#proto_core_opinion_v1_opinion-proto)
    - [OpinionCore](#proto-core-opinion-v1-OpinionCore)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proto_core_opinion_v1_opinion-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/core/opinion/v1/opinion.proto



<a name="proto-core-opinion-v1-OpinionCore"></a>

### OpinionCore
OpinionCore


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qid | [proto.core._share.v1.ShareQID](#proto-core-_share-v1-ShareQID) |  | qid - Favourite StoreQID |
| md | [proto.core._share.v1.MetaDescription](#proto-core-_share-v1-MetaDescription) | optional | md - MetaDescription |
| mmd | [proto.core._share.v1.MetaMultiDescription](#proto-core-_share-v1-MetaMultiDescription) | optional | mmd - MetaMultiDescription |
| qid_from_user | [proto.core._share.v1.ShareQID](#proto-core-_share-v1-ShareQID) |  | qid_from_user (from) |
| qid_to_asset | [proto.core._share.v1.ShareQID](#proto-core-_share-v1-ShareQID) |  | qid_to_asset (to) |
| is_favourite | [bool](#bool) | optional | is_favourite |
| stars | [int32](#int32) | optional | stars |
| asset | [proto.core.asset.v1.AssetCore](#proto-core-asset-v1-AssetCore) | optional | asset - REF User&#39;s Favourite Asset |





 

 

 

 



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

