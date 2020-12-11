package http

import (
	netHttp "net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type (
	HTTPMethodType string

	HandlerContext interface {
		Bind(interface{}) error
		JSON(int, interface{}) error
		String(int, string) error
		Get(string) interface{}
		Param(string) string
	}

	HandlerFunc func(HandlerContext) error

	handlerContext struct {
		ctx echo.Context
	}

	HTTP interface {
		Listen(string) error
		ListenMetrics(string, string)
		AddRoute(HTTPMethodType, string, HandlerFunc)
	}

	http struct {
		provider *echo.Echo
	}
)

const (
	GET    HTTPMethodType = "GET"
	PUT                   = "PUT"
	POST                  = "POST"
	DELETE                = "DELETE"
)

func (h handlerContext) JSON(statusCode int, i interface{}) error {
	return h.ctx.JSON(statusCode, i)
}

func (h handlerContext) Bind(i interface{}) error {
	return h.ctx.Bind(&i)
}

func (h handlerContext) String(code int, message string) error {
	return h.ctx.String(code, message)
}

func (h handlerContext) Get(key string) interface{} {
	return h.ctx.Get(key)
}

func (h handlerContext) Param(key string) string {
	return h.ctx.Param(key)
}

func CreateHTTP(MaxRequestsPerIPPerSecond int) HTTP {
	e := echo.New()
	ConfigRateLimitMiddleware(e, rate.Every(time.Duration(MaxRequestsPerIPPerSecond)*time.Second))

	return &http{
		provider: e,
	}
}

func ConfigRateLimitMiddleware(e *echo.Echo, MaxRequestsPerIPPerSecond rate.Limit) {
	limiter := NewIPRateLimiter(MaxRequestsPerIPPerSecond, 1)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			l := limiter.GetLimiter(ctx.Request().RemoteAddr)
			if !l.Allow() {
				return ctx.String(429, "too many requests")
			}
			return next(ctx)
		}
	})
}

func (h *http) ListenMetrics(route string, routeAddress string) {
	h.provider.Use(PrometheusMiddleware())
	go func() {
		netHttp.Handle(route, promhttp.Handler())
		logrus.Error(netHttp.ListenAndServe(routeAddress, nil))
	}()
}

func (h *http) Listen(address string) error {
	return h.provider.Start(address)
}

func (h *http) AddRoute(method HTTPMethodType, path string, handler HandlerFunc) {
	var wrappedHandler = HandlerWrapper(handler)

	switch method {
	case "GET":
		h.provider.GET(path, wrappedHandler)
		break
	case "PUT":
		h.provider.PUT(path, wrappedHandler)
		break
	case "POST":
		h.provider.POST(path, wrappedHandler)
		break
	case "DELETE":
		h.provider.DELETE(path, wrappedHandler)
		break
	}
}

func HandlerWrapper(handler HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return handler(handlerContext{
			ctx: ctx,
		})
	}
}
