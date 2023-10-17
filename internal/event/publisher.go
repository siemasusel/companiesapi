package event

import "context"

type Publisher interface {
	PushEvents(ctx context.Context, events ...Event) error
}
