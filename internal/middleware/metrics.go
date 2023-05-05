package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
)

type HttpMetricsMiddleware struct {
	mt metrics.PrometheusMetrics
}

func NewHttpMetricsMiddleware(mt metrics.PrometheusMetrics) *HttpMetricsMiddleware {
	return &HttpMetricsMiddleware{
		mt: mt,
	}
}

func (m *HttpMetricsMiddleware) MetricsMiddleware(handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error, log *zap.Logger) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		begin := time.Now()
		err := handler(w, r, p)

		var httpCode int

		if err != nil {
			errCause := errors.Cause(err)

			httpCode, _ = pkgErrors.GetHTTPCodeByError(errCause)
		} else {
			httpCode = 200
		}

		m.mt.ExecutionTime().
			WithLabelValues(strconv.Itoa(httpCode), r.URL.String(), r.Method).
			Observe(float64(time.Since(begin).Milliseconds()))

		m.mt.TotalHits().Inc()

		if 200 <= httpCode && httpCode <= 399 {
			m.mt.SuccessHits().
				WithLabelValues(strconv.Itoa(httpCode), r.URL.String(), r.Method).Inc()
		} else {
			m.mt.ErrorsHits().
				WithLabelValues(strconv.Itoa(httpCode), r.URL.String(), r.Method).Inc()
		}

		return err
	}
}
