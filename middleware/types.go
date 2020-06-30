package middleware

type Exporter interface {
	PushCount(Endpoint string)
	Push(Metric string, value map[string]interface{})
}
