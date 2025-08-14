package server

import (
	addr "github.com/polera/addr_sh/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type HandlerContainer struct {
	Log zerolog.Logger

	Hits prometheus.Counter

	LookupCache addr.Cache
}
