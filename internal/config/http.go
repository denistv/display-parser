package config

import "time"

type HTTP struct {
	// Задержка между HTTP-запросами в сервис
	// Если не использовать ограничений, сервис забанит вас на какое-то время.
	Timeout         time.Duration
	DelayPerRequest time.Duration
}
