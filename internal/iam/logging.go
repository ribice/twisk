package iam

import (
	"context"
	"time"

	"github.com/ribice/twisk/model"

	"github.com/ribice/twisk/rpc/iam"
)

var svcName = "iam"

// NewLoggingService creates new Template logging service
func NewLoggingService(svc iam.IAM, logger twisk.Logger) *LoggingService {
	return &LoggingService{
		IAM:    svc,
		logger: logger,
	}
}

// LoggingService represents iam logging service
type LoggingService struct {
	iam.IAM
	logger twisk.Logger
}

// Auth logging
func (ls *LoggingService) Auth(ctx context.Context, req *iam.AuthReq) (resp *iam.AuthResp, err error) {
	defer func(begin time.Time) {
		req.Password = "xxx-redacted-xxx"
		ls.logger.Log(
			ctx,
			svcName, "Auth request", err,
			map[string]interface{}{
				"took": time.Since(begin),
				"req":  req,
				"resp": resp,
			},
		)
	}(time.Now())
	return ls.IAM.Auth(ctx, req)
}

// Refresh token logging
func (ls *LoggingService) Refresh(ctx context.Context, req *iam.RefreshReq) (resp *iam.RefreshResp, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			ctx,
			svcName, "Refresh request", err,
			map[string]interface{}{
				"took": time.Since(begin),
				"req":  req,
				"resp": resp,
			},
		)
	}(time.Now())
	return ls.IAM.Refresh(ctx, req)
}
