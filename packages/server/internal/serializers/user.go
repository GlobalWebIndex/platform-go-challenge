package serializers

type SignUpUser struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	GivenName string `json:"givenName"`
	Surname   string `json:"surname"`
	Password  string `json:"password"`
}

type CurrentUserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	GivenName string `json:"givenName"`
	Surname   string `json:"surname"`
}
