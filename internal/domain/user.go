package domain

type User struct {
	ID        uint64
	Email     string
	Password  string
	Favorites []Asset
}
