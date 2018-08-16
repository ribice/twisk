package user

import (
	"context"
	"time"

	"github.com/ribice/twisk/model"

	"github.com/ribice/twisk/rpc/user"
)

var svcName = "user"

// NewLoggingService creates new Template logging service
func NewLoggingService(svc user.User, logger twisk.Logger) *LoggingService {
	return &LoggingService{
		User:   svc,
		logger: logger,
	}
}

// LoggingService represents iam logging service
type LoggingService struct {
	user.User
	logger twisk.Logger
}

// Create user logging
func (ls *LoggingService) Create(ctx context.Context, req *user.CreateReq) (resp *user.Resp, err error) {
	defer func(begin time.Time) {
		req.Password = "xxx-redacted-xxx"
		ls.logger.Log(
			ctx,
			svcName, "Create ticket request", err,
			map[string]interface{}{
				"took": time.Since(begin),
				"req":  req,
				"resp": resp,
			},
		)
	}(time.Now())
	return ls.User.Create(ctx, req)
}

// View user logging
func (ls *LoggingService) View(ctx context.Context, req *user.IDReq) (resp *user.Resp, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			ctx,
			svcName, "View user request", err,
			map[string]interface{}{
				"took": time.Since(begin),
				"req":  req,
				"resp": resp,
			},
		)
	}(time.Now())
	return ls.User.View(ctx, req)
}

// List user logging
func (ls *LoggingService) List(ctx context.Context, req *user.ListReq) (resp *user.ListResp, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			ctx,
			svcName, "List users request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
				"resp": resp,
			},
		)
	}(time.Now())
	return ls.User.List(ctx, req)
}

// Update user logging
func (ls *LoggingService) Update(ctx context.Context, req *user.UpdateReq) (resp *user.Resp, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			ctx,
			svcName, "Update user request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.User.Update(ctx, req)
}

// Delete user logging
func (ls *LoggingService) Delete(ctx context.Context, req *user.IDReq) (resp *user.MessageResp, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			ctx,
			svcName, "Delete user request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
				"resp": resp,
			},
		)
	}(time.Now())
	return ls.User.Delete(ctx, req)
}
