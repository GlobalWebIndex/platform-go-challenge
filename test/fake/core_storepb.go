package fake

import (
	storepb "x-gwi/proto/core/_store/v1"
	"x-gwi/service"
)

// storepb

// Store
func FakeStore(seed int64, cn service.CoreName) *storepb.Store {
	if seed == 0 {
		return &storepb.Store{}
	}
	return &storepb.Store{
		// Qid: FakeShareQID(seed,cn,key),
		User:  FakeUserCore(seed),
		Asset: FakeAssetCore(seed),
		// Favourite: FakeFavouriteCore(seed, fom, to),
		// Opinion: FakeOpinionCore(seed,from,to),
	}
}

// StoreAQL
func FakeStoreAQL(seed int64, cn service.CoreName) *storepb.StoreAQL {
	if seed == 0 {
		return &storepb.StoreAQL{}
	}
	return &storepb.StoreAQL{
		XKey:    "",
		XId:     "",
		XRev:    "",
		XOldRev: "",
		XFrom:   "",
		XTo:     "",
		// Qid: FakeShareQID(seed,cn,key),
		User:  FakeUserCore(seed),
		Asset: FakeAssetCore(seed),
		// Favourite: FakeFavouriteCore(seed, fom, to),
		// Opinion: FakeOpinionCore(seed,from,to),
	}
}
