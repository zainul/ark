package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/zainul/ark/xrandom"
)

const (
	RequestID  = "request_id"
	XRoundTrip = "x-roundtrip"
)

type middleware struct {
	ServiceName string
	Exporter    Exporter
}

type Option func(m *middleware)

func WithOptionExporter(expo Exporter) Option {
	return func(m *middleware) {
		m.Exporter = expo
	}
}

func New(serviceName string, options ...Option) *middleware {

	m := middleware{ServiceName: serviceName}

	for _, opt := range options {
		opt(&m)
	}

	return &m
}

// ContentType is middleware content type
func (m *middleware) ContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get(RequestID) == "" {
			rd, _ := xrandom.GenerateRandomString(20)
			rqID := m.ServiceName + rd
			c.Request.Header.Set(RequestID, rqID)
		}

		if c.Request.Header.Get(XRoundTrip) == "" {
			c.Request.Header.Set(XRoundTrip, m.ServiceName)
		}

		ctx := context.WithValue(c.Request.Context(), XRoundTrip, c.Request.Header.Get(XRoundTrip))
		xRound := c.Request.Header.Get(XRoundTrip)
		ctx = context.WithValue(ctx, XRoundTrip, "->"+xRound)

		c.Writer.Header().Add("Content-Type", "application/json")
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
