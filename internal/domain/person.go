package domain

type Person struct {
}

type Role string

const (
	ADMIN    Role = "admin"
	USER     Role = "user"
	BUSINESS Role = "business"
)
