package helpers

import (
	"gwiEng/pkg/logs"
	"net/http"

	"go.uber.org/zap"
)

type MsgResponse struct {
	Msg string `json:"msg"`
}

func RespondAndLogError(action, collection, resourceID, requestID string, err error, code int, w http.ResponseWriter, logger *zap.Logger) error {
	http.Error(w, err.Error(), code)
	logs.LogErrorAction(logs.ErrLogMsg{
		Action:     action,
		Collection: collection,
		ResourceID: resourceID,
		RequestID:  requestID,
		Error:      err,
	}, logger)
	return nil
}
