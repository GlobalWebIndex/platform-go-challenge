package fake

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	opinionpb "x-gwi/proto/core/opinion/v1"
	"x-gwi/service"
)

// opinionpb

func keyOpinion(from, to int64) string {
	// return fmt.Sprintf("fk-uao:(%s,%s)", keyUser(from), keyAsset(to))
	// on unique idx from to -> ""
	// return ""
	return fmt.Sprintf("fk-uao:(%s,%s)", keyUser(from), keyAsset(to))
}

// FavouriteCore
func FakeOpinionCore(seed, from, to int64) *opinionpb.OpinionCore {
	if seed == 0 {
		return &opinionpb.OpinionCore{}
	}
	f := gofakeit.NewCrypto()
	return &opinionpb.OpinionCore{
		Qid:         FakeShareQID(seed, service.NameOpinion, keyOpinion(from, to)),
		QidFromUser: FakeShareQID(from, service.NameUser, keyUser(from)),
		QidToAsset:  FakeShareQID(to, service.NameAsset, keyAsset(to)),
		Md:          FakeMetaDescription(f.Int64()),
		Mmd:         FakeMetaMultiDescription(f.Int64()),
		IsFavourite: PtrBool(f.Bool(), f.Bool()),
		Stars:       PtrInt32(int32(f.IntRange(0, 10))),
		// Asset:       FakeAssetCore(to),
	}
}
