package fake

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	favouritepb "x-gwi/proto/core/favourite/v1"
	"x-gwi/service"
)

// favouritepb

func keyFavourite(from, to int64) string {
	return fmt.Sprintf("fk-uaf:(%s,%s)", keyUser(from), keyAsset(to))
}

// FavouriteCore
func FakeFavouriteCore(seed, from, to int64) *favouritepb.FavouriteCore {
	if seed == 0 {
		return &favouritepb.FavouriteCore{}
	}
	f := gofakeit.NewCrypto()
	return &favouritepb.FavouriteCore{
		Qid:         FakeShareQID(seed, service.NameFavourite, keyFavourite(from, to)),
		QidFromUser: FakeShareQID(from, service.NameUser, keyUser(from)),
		QidToAsset:  FakeShareQID(to, service.NameAsset, keyAsset(to)),
		Md:          FakeMetaDescription(f.Int64()),
		Mmd:         FakeMetaMultiDescription(f.Int64()),
		// Asset:       FakeAssetCore(to),
	}
}
