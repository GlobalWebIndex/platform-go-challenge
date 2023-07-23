package domain

type User struct {
	ID        uint64
	Email     string
	Password  string
	Favorites map[uint64]bool
}
