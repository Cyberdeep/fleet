package service

import (
	"context"
	"encoding/json"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/kolide/fleet/server/kolide"
)

func (mw loggingMiddleware) EnrollAgent(ctx context.Context, enrollSecret string, hostIdentifier string) (string, error) {
	var (
		nodeKey string
		err     error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "EnrollAgent",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	nodeKey, err = mw.Service.EnrollAgent(ctx, enrollSecret, hostIdentifier)
	return nodeKey, err
}

func (mw loggingMiddleware) AuthenticateHost(ctx context.Context, nodeKey string) (*kolide.Host, error) {
	var (
		host *kolide.Host
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "AuthenticateHost",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	host, err = mw.Service.AuthenticateHost(ctx, nodeKey)
	return host, err
}

func (mw loggingMiddleware) GetClientConfig(ctx context.Context) (*kolide.OsqueryConfig, error) {
	var (
		config *kolide.OsqueryConfig
		err    error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetClientConfig",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	config, err = mw.Service.GetClientConfig(ctx)
	return config, err
}

func (mw loggingMiddleware) GetDistributedQueries(ctx context.Context) (map[string]string, uint, error) {
	var (
		queries    map[string]string
		err        error
		accelerate uint
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetDistributedQueries",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	queries, accelerate, err = mw.Service.GetDistributedQueries(ctx)
	return queries, accelerate, err
}

func (mw loggingMiddleware) SubmitDistributedQueryResults(ctx context.Context, results kolide.OsqueryDistributedQueryResults, statuses map[string]kolide.OsqueryStatus) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SubmitDistributedQueryResults",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.SubmitDistributedQueryResults(ctx, results, statuses)
	return err
}

func (mw loggingMiddleware) SubmitStatusLogs(ctx context.Context, logs []json.RawMessage) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SubmitStatusLogs",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.SubmitStatusLogs(ctx, logs)
	return err
}

func (mw loggingMiddleware) SubmitResultLogs(ctx context.Context, logs []json.RawMessage) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SubmitResultLogs",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.SubmitResultLogs(ctx, logs)
	return err
}
