package gwi

type User struct {
	Id         string  `json:"id"`
	Favourites []Asset `json:"favs"`
}
