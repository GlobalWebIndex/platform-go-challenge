package domain

import "context"

type App interface {
	Start(context.Context) error
	Stop()
}
