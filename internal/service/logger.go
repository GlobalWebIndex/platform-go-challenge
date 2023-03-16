package service

import (
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
)

type LoggerService interface {
	LogUserActivity(userActivity *dto.UserActivity)
}

type loggerService struct {
	dbHandler repository.DBHandler
}

func NewloggerService(dbHandler repository.DBHandler) LoggerService {
	return &loggerService{dbHandler}
}

// SendMessage implements NotifyService
func (l *loggerService) LogUserActivity(userActivity *dto.UserActivity) {
	l.dbHandler.NewloggerService().LogUserActivity(userActivity)
}
