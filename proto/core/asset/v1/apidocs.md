# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [proto/core/asset/v1/asset.proto](#proto_core_asset_v1_asset-proto)
    - [AssetAudience](#proto-core-asset-v1-AssetAudience)
    - [AssetChart](#proto-core-asset-v1-AssetChart)
    - [AssetInsight](#proto-core-asset-v1-AssetInsight)
    - [AssetInstance](#proto-core-asset-v1-AssetInstance)
    - [AssetMetaData](#proto-core-asset-v1-AssetMetaData)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proto_core_asset_v1_asset-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/core/asset/v1/asset.proto



<a name="proto-core-asset-v1-AssetAudience"></a>

### AssetAudience
AssetAudience -  (which is a series of characteristics, for that exercise
lets focus on gender (Male, Female), birth country, age groups, hours spent
daily on social media, number of purchases last month) e.g. Males from 24-35
that spent more than 3 hours on social media daily.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| md | [AssetMetaData](#proto-core-asset-v1-AssetMetaData) | optional | md - Asset Audience MetaData |
| gender | [string](#string) | optional | gender - (male | female) |
| genders | [string](#string) | repeated | genders - [](male, female) |
| country_code | [string](#string) | optional | country_code - iso 2 birth country code (auto lowercase) |
| country_codes | [string](#string) | repeated | country_codes - [](0-300 elements, unique, code max_len 2) optional |
| age_min | [uint32](#uint32) | optional | age_min - age groups (1-100) |
| age_max | [uint32](#uint32) | optional | age_max - age groups (1-100) |
| hours_min | [uint32](#uint32) | optional | hours_min - hours on social (1-24) hours spent daily on social media |
| hours_max | [uint32](#uint32) | optional | hours_max - hours on social (1-24) hours spent daily on social media |
| purchases_min | [uint32](#uint32) | optional | purchases_min - (1-100_000) number of purchases last month |
| purchases_max | [uint32](#uint32) | optional | purchases_max - (1-100_000) number of purchases last month |






<a name="proto-core-asset-v1-AssetChart"></a>

### AssetChart
AssetChart - (that has a small title, axes titles and data)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| md | [AssetMetaData](#proto-core-asset-v1-AssetMetaData) | optional | md - Asset Chart MetaData |
| title | [string](#string) | optional | title - (0-256 characters) optional |
| data | [google.protobuf.Struct](#google-protobuf-Struct) | optional | data |
| data_raw | [bytes](#bytes) | optional | data_raw - raw binary data (65536) |
| options | [google.protobuf.Struct](#google-protobuf-Struct) | optional | options |






<a name="proto-core-asset-v1-AssetInsight"></a>

### AssetInsight
AssetInsight - (a small piece of text that provides some insight into a
topic, e.g. &#34;40% of millenials spend more than 3hours on social media daily&#34;)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| md | [AssetMetaData](#proto-core-asset-v1-AssetMetaData) | optional | md - Asset Insight MetaData |
| sentence | [string](#string) | optional | sentence - (1-256 characters) required |






<a name="proto-core-asset-v1-AssetInstance"></a>

### AssetInstance
AssetInstance


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [proto.core.idx.v1.IDX](#proto-core-idx-v1-IDX) |  | id - AssetIDX |
| md | [AssetMetaData](#proto-core-asset-v1-AssetMetaData) | optional | md - AssetMetaData |
| chart | [AssetChart](#proto-core-asset-v1-AssetChart) | optional | chart - Chart (that has a small title, axes titles and data) |
| insight | [AssetInsight](#proto-core-asset-v1-AssetInsight) | optional | insight - Insight (a small piece of text that provides some insight into a topic, e.g. &#34;40% of millenials spend more than 3hours on social media daily&#34;) |
| audience | [AssetAudience](#proto-core-asset-v1-AssetAudience) | optional | audience - Audience (which is a series of characteristics, for that exercise lets focus on gender (Male, Female), birth country, age groups, hours spent daily on social media, number of purchases last month) e.g. Males from 24-35 that spent more than 3 hours on social media daily. |






<a name="proto-core-asset-v1-AssetMetaData"></a>

### AssetMetaData
AssetMetaData - title, topic, label, description, tags


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title | [string](#string) | optional | title - (0-256 characters) optional |
| topic | [string](#string) | optional | topic - (0-256 characters) optional |
| label | [string](#string) | optional | label - (0-256 characters) optional |
| description | [string](#string) | optional | description - (0-1024 characters) optional |
| tags | [string](#string) | repeated | tags - (0-32 elements, unique, tag max_len 32) optional |





 

 

 

 



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

