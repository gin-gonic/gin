package concurrent

import "context"

type Executor interface {
	Go(handler func(ctx context.Context))
}