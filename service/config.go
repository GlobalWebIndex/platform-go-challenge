package service

type CoreName string

const (
	NameAsset     CoreName = CoreName("asset")
	NameUser      CoreName = CoreName("user")
	NameFavourite CoreName = CoreName("favourite")
	NameOpinion   CoreName = CoreName("opinion")
)

func CoreNames() []CoreName {
	return []CoreName{
		NameAsset,
		NameUser,
		NameFavourite,
		NameOpinion,
	}
}
