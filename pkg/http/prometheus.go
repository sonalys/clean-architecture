package http

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_requests",
			Help: "Total number of requests processed by the api",
		},
		[]string{"url", "method", "status_code", "size", "duration"},
	)
	requestsDurationHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "requests_duration",
			Help:    "Histogram of app's requests duration",
			Buckets: prometheus.LinearBuckets(10, 100, 15),
		},
	)

	requestsPayloadHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "requests_payload",
			Help:    "Histogram of app's requests payload size, in bytes",
			Buckets: prometheus.LinearBuckets(1000, 10000, 15),
		},
	)
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	prometheus.MustRegister(requestsCounterVec)
	prometheus.MustRegister(requestsDurationHistogram)
	prometheus.MustRegister(requestsPayloadHistogram)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			t1 := time.Now()
			defer func() {
				duration := float64(time.Since(t1) / time.Millisecond)
				req := c.Request()
				res := c.Response()
				url := c.Path()
				statusCode := strconv.Itoa(res.Status)
				method := req.Method
				size := float64(res.Size)
				sizeLabel := fmt.Sprintf("%f", size)
				durationLabel := fmt.Sprintf("%f", duration)
				requestsCounterVec.WithLabelValues(url, method, statusCode, sizeLabel, durationLabel).Inc()
				requestsDurationHistogram.Observe(duration)
				requestsPayloadHistogram.Observe(size)
			}()
			return next(c)
		}
	}
}
