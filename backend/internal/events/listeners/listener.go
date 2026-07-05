package listeners

import (
	"context"
)

type Listener interface {
	Handle(ctx context.Context, params ...any) error
}
