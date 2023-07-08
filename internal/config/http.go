package config

import "time"

type HTTP struct {
	Timeout time.Duration

	// Задержка между HTTP-запросами в сервис
	// Если не использовать ограничений, сервис забанит вас на какое-то время.
	DelayPerRequest time.Duration
}
