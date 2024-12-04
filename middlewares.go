package shortener

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostShorten(ctx context.Context, p Shorten) (shorten ShortenResp, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostShorten", "id", p, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostShorten(ctx, p)
}

func (mw loggingMiddleware) Resolve(ctx context.Context, id string) (url string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Resolve", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Resolve(ctx, id)
}
