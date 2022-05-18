package pool

import (
	"context"

	"github.com/amplify-security/probe"
	"github.com/rs/zerolog"
)

type (
	Pool struct {
		log      *zerolog.Logger
		ctx      context.Context
		childCtx context.Context
		cancel   context.CancelFunc
		work     chan probe.Runner
	}
)

func NewPool(ctx context.Context, log *zerolog.Logger) *Pool {
	ctxLogger := log.With().Str("source", "probe.Pool").Logger()
	childCtx, cancel := context.WithCancel(ctx)
	return &Pool{
		log:    &ctxLogger,
		ctx:    childCtx,
		cancel: cancel,
	}
}
