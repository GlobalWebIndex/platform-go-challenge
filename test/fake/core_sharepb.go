package fake

import (
	"github.com/brianvoe/gofakeit/v6"

	"x-gwi/app/x/id"
	sharepb "x-gwi/proto/core/_share/v1"
	"x-gwi/service"
)

// sharepb

// ShareQID
func FakeShareQID(seed int64, cn service.CoreName, key string) *sharepb.ShareQID {
	if seed == 0 {
		return &sharepb.ShareQID{}
	}
	// "github.com/brianvoe/gofakeit/v6"
	// f := gofakeit.New(seed)
	// f.UUID()
	return &sharepb.ShareQID{
		Kind: cn.String(),
		Key:  key,
		Uid:  id.XiD().String(),
		Uuid: id.UUID().String(),
		Rev:  id.Rev(),
	}
}

// MetaDescription
func FakeMetaDescription(seed int64) *sharepb.MetaDescription {
	if seed == 0 {
		return &sharepb.MetaDescription{}
	}
	f := gofakeit.New(seed)
	return &sharepb.MetaDescription{
		Title:       PtrStr(f.SentenceSimple()),
		Topic:       PtrStr(f.Phrase()),
		Label:       PtrStr(f.Noun()),
		Description: PtrStr(f.HipsterSentence(9)),
	}
}

// MetaMultiDescription
func FakeMetaMultiDescription(seed int64) *sharepb.MetaMultiDescription {
	if seed == 0 {
		return &sharepb.MetaMultiDescription{}
	}
	f := gofakeit.New(seed)
	f2 := gofakeit.New(seed + 3)
	return &sharepb.MetaMultiDescription{
		Labels: []string{f.NounAbstract(), f2.NounConcrete()},
		Tags:   []string{f2.NounAbstract(), f.NounConcrete()},
	}
}
