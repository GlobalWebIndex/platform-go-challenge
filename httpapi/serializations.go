package httpapi

import "platform-go-challenge/domain"

func fromUserDomainToUserJson(user domain.User) UserJson {
	return UserJson{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}
}
