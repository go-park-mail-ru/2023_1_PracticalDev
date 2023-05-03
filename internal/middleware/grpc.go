package middleware

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCMetricsMiddleware struct {
	mt metrics.PrometheusMetrics
}

func NewGRPCMetricsMiddleware(mt metrics.PrometheusMetrics) *GRPCMetricsMiddleware {
	return &GRPCMetricsMiddleware{mt: mt}
}

func (m *GRPCMetricsMiddleware) MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	begin := time.Now()

	reply, err := handler(ctx, req)
	errStatus, _ := status.FromError(err)
	code := errStatus.Code()

	m.mt.ExecutionTime().
		WithLabelValues(code.String(), info.FullMethod, info.FullMethod).
		Observe(time.Since(begin).Seconds())

	if code == codes.OK {
		m.mt.SuccessHits().WithLabelValues(codes.OK.String(), "", "").Inc()
	} else {
		m.mt.ErrorsHits().WithLabelValues(code.String(), "", "").Inc()
	}

	m.mt.TotalHits().Inc()

	return reply, err
}
