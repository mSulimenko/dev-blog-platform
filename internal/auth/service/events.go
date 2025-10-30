package service

import "context"

type EventDispatcher interface {
	UserRegistered(ctx context.Context, email, token, username string) error
}
