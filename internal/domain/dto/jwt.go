package dto

type JWTPayload struct {
	UserID string `json:"userID" binding:"required"`
}

func (j *JWTPayload) GetUserID() string {
	return j.UserID
}
