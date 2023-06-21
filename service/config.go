package service

type CoreName string

const (
	NameUser      CoreName = CoreName("u_user")
	NameAsset     CoreName = CoreName("a_asset")
	NameFavourite CoreName = CoreName("uaf_favourite")
	NameOpinion   CoreName = CoreName("uao_opinion")
)

func CoreNames() []CoreName {
	return []CoreName{
		NameAsset,
		NameUser,
		NameFavourite,
		NameOpinion,
	}
}

func Edges() []CoreName {
	return []CoreName{
		NameFavourite,
		NameOpinion,
	}
}

func (x CoreName) String() string {
	return string(x)
}

func (x CoreName) Valid() bool {
	for _, v := range CoreNames() {
		if x == v {
			return true
		}
	}

	return false
}

func (x CoreName) IsEdge() bool {
	for _, v := range Edges() {
		if x == v {
			return true
		}
	}

	return false
}
