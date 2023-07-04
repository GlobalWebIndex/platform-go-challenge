# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [proto/serv/_user/v1/user_srv.proto](#proto_serv__user_v1_user_srv-proto)
    - [CreateRequest](#proto-serv-_user-v1-CreateRequest)
    - [CreateResponse](#proto-serv-_user-v1-CreateResponse)
    - [DeleteRequest](#proto-serv-_user-v1-DeleteRequest)
    - [DeleteResponse](#proto-serv-_user-v1-DeleteResponse)
    - [GetRequest](#proto-serv-_user-v1-GetRequest)
    - [GetResponse](#proto-serv-_user-v1-GetResponse)
    - [ListRequest](#proto-serv-_user-v1-ListRequest)
    - [ListResponse](#proto-serv-_user-v1-ListResponse)
    - [UpdateRequest](#proto-serv-_user-v1-UpdateRequest)
    - [UpdateResponse](#proto-serv-_user-v1-UpdateResponse)
  
    - [UserService](#proto-serv-_user-v1-UserService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proto_serv__user_v1_user_srv-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/serv/_user/v1/user_srv.proto



<a name="proto-serv-_user-v1-CreateRequest"></a>

### CreateRequest
CreateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) |  | user |






<a name="proto-serv-_user-v1-CreateResponse"></a>

### CreateResponse
CreateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) | optional | user |






<a name="proto-serv-_user-v1-DeleteRequest"></a>

### DeleteRequest
DeleteRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qid | [proto.core._share.v1.ShareQID](#proto-core-_share-v1-ShareQID) |  | qid |






<a name="proto-serv-_user-v1-DeleteResponse"></a>

### DeleteResponse
DeleteResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) |  | user |






<a name="proto-serv-_user-v1-GetRequest"></a>

### GetRequest
GetRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qid | [proto.core._share.v1.ShareQID](#proto-core-_share-v1-ShareQID) |  | qid |






<a name="proto-serv-_user-v1-GetResponse"></a>

### GetResponse
GetResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) |  | user |






<a name="proto-serv-_user-v1-ListRequest"></a>

### ListRequest
ListRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) |  | user |






<a name="proto-serv-_user-v1-ListResponse"></a>

### ListResponse
ListResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) |  | user |






<a name="proto-serv-_user-v1-UpdateRequest"></a>

### UpdateRequest
UpdateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) |  | user |
| update_mask | [google.protobuf.FieldMask](#google-protobuf-FieldMask) |  | update_mask https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/patch_feature |






<a name="proto-serv-_user-v1-UpdateResponse"></a>

### UpdateResponse
UpdateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [proto.core._user.v1.UserCore](#proto-core-_user-v1-UserCore) |  | user |





 

 

 


<a name="proto-serv-_user-v1-UserService"></a>

### UserService
UserService
a microservice to handle companies. It provides the following operations:
• Create, • Patch, • Delete, • Get (one)
standard methods: List, Get, Create, Update, and Delete

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [CreateRequest](#proto-serv-_user-v1-CreateRequest) | [CreateResponse](#proto-serv-_user-v1-CreateResponse) | Create |
| Get | [GetRequest](#proto-serv-_user-v1-GetRequest) | [GetResponse](#proto-serv-_user-v1-GetResponse) | Get |
| Gett | [GetRequest](#proto-serv-_user-v1-GetRequest) | [GetResponse](#proto-serv-_user-v1-GetResponse) | Gett |
| Update | [UpdateRequest](#proto-serv-_user-v1-UpdateRequest) | [UpdateResponse](#proto-serv-_user-v1-UpdateResponse) | Update |
| Delete | [DeleteRequest](#proto-serv-_user-v1-DeleteRequest) | [DeleteResponse](#proto-serv-_user-v1-DeleteResponse) | Delete |
| List | [ListRequest](#proto-serv-_user-v1-ListRequest) | [ListResponse](#proto-serv-_user-v1-ListResponse) stream | List/Filter - stream user matching request |

 



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

