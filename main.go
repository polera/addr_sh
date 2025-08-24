package main

import (
	"errors"
	"fmt"
	"github.com/polera/addr_sh/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
	"github.com/tristanfisher/patchpanel"
	"net/http"
	"os"
	"time"

	"github.com/polera/addr_sh/server"

	"github.com/polera/addr_sh/pkg"
)

func main() {
	// grab a values/configuration file path from our environment using patchpanel
	valuesFile := patchpanel.GetFileEnvOrPath(patchpanel.ENV_CONFIG_FILE, patchpanel.FLAG_CONFIG_FILE)

	conf, err := config.ParseConfig(valuesFile, config.Config{})
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "error parsing configuration: %s\n", err.Error())
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	hc := server.HandlerContainer{
		Log: logger,
		Hits: promauto.NewCounter(prometheus.CounterOpts{
			Name: "requests_total",
			Help: "The total number of requests",
		}),
		LookupCache: addr.Cache{
			Values: make(map[string]*addr.IPv4CIDR),
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", hc.Index)
	mux.HandleFunc("/about", hc.GetAbout)
	mux.HandleFunc("/cidr/v4/{address}/{cidr}", hc.Getv4cidrDetail)
	mux.HandleFunc("/cidr/v4/{cidr}", hc.Getv4cidrData)
	mux.HandleFunc("/headers", hc.RequestHeaders)
	mux.HandleFunc("/hostnames/{address}", hc.RequestHost)
	mux.HandleFunc("/hostnames", hc.RequestHost)
	mux.HandleFunc("/ip", hc.RemoteAddress)

	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
		Addr:         conf.ListenPort,
	}

	if conf.EnableTLS {
		go func() {
			logger.Info().Msg(fmt.Sprintf("TLS Enabled: %v", conf.EnableTLS))
			httpsServer := &http.Server{
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				IdleTimeout:  120 * time.Second,
				Handler:      mux,
				Addr:         conf.TLSListenPort,
			}
			logger.Info().Str("tlsListenPort", conf.TLSListenPort).Msg("starting TLS server")
			err := httpsServer.ListenAndServeTLS("fullchain.pem", "privkey.pem")
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal().Err(err).Msg("error starting TLS server")
			}
		}()
	}

	logger.Info().Str("listenPort", conf.ListenPort).Msg("starting HTTP server")
	err = httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error().Err(err).Msg("error starting HTTP server")
	}
	logger.Info().Msg("Shutting down HTTP server")
}
