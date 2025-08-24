package logs

import "go.uber.org/zap"

type LogMsg struct {
	Action string
	Collection string
	ResourceID string
	RequestID string
}

type ErrLogMsg struct {
	Action string
	Collection string
	ResourceID string
	RequestID string
	Error error
}

func LogAction(logMsg LogMsg, logger *zap.Logger) {
	logger.Info("Action Added", 
		zap.String("action", logMsg.Action),
		zap.String("collection", logMsg.Collection),
		zap.String("resourceID", logMsg.ResourceID),
		zap.String("requestID", logMsg.RequestID))
}

func LogErrorAction(errLogMsg ErrLogMsg, logger *zap.Logger) {
	logger.Error("Action Failed",
		zap.String("action", errLogMsg.Action),
		zap.String("collection", errLogMsg.Collection),
		zap.String("resourceID", errLogMsg.ResourceID),
		zap.String("requestID", errLogMsg.RequestID),
		zap.String("error", errLogMsg.Error.Error()))
}
