package token

type RefreshRepository interface {
	Save(token string, record RefreshRecord)
	Get(token string) (RefreshRecord, bool)
	Delete(token string)
}
