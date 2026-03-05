package server

import (
	"net"
	"net/http"

	addr "github.com/polera/addr_sh/pkg"
)

func (hc *HandlerContainer) RemoteAddress(w http.ResponseWriter, r *http.Request) {
	hc.Hits.Inc()
	hc.Log.Debug().Msg("called remoteAddress")
	remoteAddr := addr.GetRemoteHost(r)
	ipMap := map[string]string{"ip": remoteAddr.String()}
	writeJSON(w, http.StatusOK, ipMap)
}

func (hc *HandlerContainer) RequestHeaders(w http.ResponseWriter, r *http.Request) {
	hc.Hits.Inc()
	hc.Log.Debug().Msg("called requestHeaders")
	writeJSON(w, http.StatusOK, r.Header)
}

func (hc *HandlerContainer) RequestHost(w http.ResponseWriter, r *http.Request) {
	hc.Hits.Inc()
	hc.Log.Debug().Msg("called requestHost")
	remoteAddr := addr.GetRemoteHost(r)
	passedAddress := r.PathValue("address")

	if passedAddress != "" {
		parsedPassedAddress := net.ParseIP(passedAddress)
		if parsedPassedAddress != nil {
			remoteAddr = &parsedPassedAddress
		}
	}
	ptr, err := net.LookupAddr(remoteAddr.String())
	if err != nil {
		ptr = []string{"Not found"}
	}
	writeJSON(w, http.StatusOK, ptr)
}
