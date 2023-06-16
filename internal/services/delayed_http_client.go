package services

import (
	"context"
	"net/http"
	"time"
)

func NewDelayedHTTPClient(ctx context.Context, delayPerReq time.Duration, httpClient HTTPClient) *DelayedHTTPClient {
	d := DelayedHTTPClient{
		ticker:     time.NewTicker(delayPerReq),
		httpClient: httpClient,
	}

	go func() {
		<-ctx.Done()
		d.ticker.Stop()
	}()

	return &d
}

// DelayedHTTPClient с простой реализацией задержек при отправке запросов во внешний мир. Поскольку сайт displayspecifications
// банит клиента, если он совершает слишком много запросов, необходимо глобально в рамках приложения ограничивать нагрузку
// на внешний ресурс. Инжектим инстанс этого клиента во все сервисы, после чего каждый последующий запрос будет выполнен
// только тогда, когда сработает тикер. Если из нескольких горутин вызвать DelayedHTTPClient.Do(), то http-запрос
// уйдет в сеть только у одной горутины, остальные будут ждать следующего тика.
type DelayedHTTPClient struct {
	ticker     *time.Ticker
	httpClient HTTPClient
}

func (d *DelayedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	<-d.ticker.C
	return d.httpClient.Do(req)
}
