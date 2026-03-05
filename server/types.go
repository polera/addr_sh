package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type HandlerContainer struct {
	Log zerolog.Logger

	Hits prometheus.Counter
}
