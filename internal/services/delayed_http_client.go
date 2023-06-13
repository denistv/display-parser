package services

import (
	"context"
	"net/http"
	"time"
)

func NewDelayedHTTPClient(ctx context.Context, delay time.Duration, httpClient HTTPClient) *DelayedHTTPClient {
	d := DelayedHTTPClient{
		ticker:     time.NewTicker(delay),
		httpClient: httpClient,
	}

	go func() {
		<-ctx.Done()
		d.ticker.Stop()
	}()

	return &d
}

// Обертка над тикером с урезанной функциональностью, чтобы не передавать сырой тикер в сервисы системы.
// Поскольку мы обращаемся во внешний сервис, нам нужен тикер-синглтон, который будет отстреливать в момент,
// когда множеству сервисов нашей системы можно отправить запрос
// Так же необходимо запретить возможность остановки тикера внутри сервисов.
type DelayedHTTPClient struct {
	ticker     *time.Ticker
	httpClient HTTPClient
}

func (d *DelayedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	<-d.ticker.C
	return d.httpClient.Do(req)
}
