package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

func Log(ctx context.Context, req *app.RequestContext) {
	start := time.Now()
	req.Next(ctx)
	end := time.Now()
	latency := end.Sub(start).Microseconds
	hlog.CtxTracef(ctx, "status=%d cost=%d method=%s full_path=%s client_ip=%s host=%s query=%s",
		req.Response.StatusCode(), latency,
		req.Request.Header.Method(), req.Request.URI().PathOriginal(), req.ClientIP(), req.Request.Host(),
		req.Request.QueryString())
}
