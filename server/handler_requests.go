package server

import (
	"encoding/json"
	"net"
	"net/http"

	addr "github.com/polera/addr_sh/pkg"
)

func (hc *HandlerContainer) RemoteAddress(w http.ResponseWriter, r *http.Request) {
	hc.Hits.Inc()
	hc.Log.Debug().Msg("called remoteAddress")
	remoteAddr := addr.GetRemoteHost(r)
	addr := make(map[string]string)
	addr["ip"] = remoteAddr.String()
	jsonRes, _ := json.Marshal(addr)
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonRes)
}

func (hc *HandlerContainer) RequestHeaders(w http.ResponseWriter, r *http.Request) {
	hc.Hits.Inc()
	hc.Log.Debug().Msg("called requestHeaders")
	jsonRes, _ := json.Marshal(r.Header)
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonRes)
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
	jsonRes, _ := json.Marshal(ptr)
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonRes)
}
