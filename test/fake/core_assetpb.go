//nolint:revive,exhaustruct,gomnd
package fake

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/protobuf/types/known/structpb"

	assetpb "x-gwi/proto/core/asset/v1"
	"x-gwi/service"
)

// assetpb

func keyAsset(seed int64) string {
	if seed == 0 {
		return ""
	}

	return fmt.Sprintf("fk-a:%d", seed)
}

// AssetCore
func FakeAssetCore(seed int64) *assetpb.AssetCore {
	if seed == 0 {
		return &assetpb.AssetCore{}
	}

	return &assetpb.AssetCore{
		Qid:      FakeShareQID(seed, service.NameAsset, keyAsset(seed)),
		Md:       FakeMetaDescription(seed),
		Mmd:      FakeMetaMultiDescription(seed),
		Chart:    FakeAssetChart(seed),
		Insight:  FakeAssetInsight(seed),
		Audience: FakeAssetAudience(seed),
	}
}

// AssetChart
func FakeAssetChart(seed int64) *assetpb.AssetChart {
	if seed == 0 {
		return &assetpb.AssetChart{}
	}

	f := gofakeit.New(seed + 41)

	return &assetpb.AssetChart{
		Title:   PtrStr(f.SentenceSimple()),
		Md:      FakeMetaDescription(seed + 77),
		Mmd:     FakeMetaMultiDescription(seed + 88),
		Data:    &structpb.Struct{},
		DataRaw: []byte{},
		Options: &structpb.Struct{},
	}
}

// AssetInsight
func FakeAssetInsight(seed int64) *assetpb.AssetInsight {
	if seed == 0 {
		return &assetpb.AssetInsight{}
	}

	f := gofakeit.New(seed + 32)

	return &assetpb.AssetInsight{
		Sentence: PtrStr(f.HipsterSentence(9)),
		Md:       FakeMetaDescription(seed + 64),
		Mmd:      FakeMetaMultiDescription(seed + 128),
	}
}

// AssetAudience
func FakeAssetAudience(seed int64) *assetpb.AssetAudience {
	if seed == 0 {
		return &assetpb.AssetAudience{}
	}

	f := gofakeit.New(seed + 25)

	return &assetpb.AssetAudience{
		Md:           FakeMetaDescription(seed + 53),
		Mmd:          FakeMetaMultiDescription(seed + 87),
		Gender:       PtrStr(f.Gender()),
		Genders:      []string{f.Gender()},
		CountryCode:  PtrStr(f.CountryAbr()),
		CountryCodes: []string{f.CountryAbr()},
		AgeMin:       PtrUInt32(uint32(f.UintRange(0, 100))),
		AgeMax:       PtrUInt32(uint32(f.UintRange(0, 100))),
		HoursMin:     PtrUInt32(uint32(f.UintRange(0, 24))),
		HoursMax:     PtrUInt32(uint32(f.UintRange(0, 24))),
		PurchasesMin: PtrUInt32(uint32(f.UintRange(0, 100_000))),
		PurchasesMax: PtrUInt32(uint32(f.UintRange(0, 100_000))),
	}
}
