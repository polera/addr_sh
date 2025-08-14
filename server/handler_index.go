package server

import (
	"encoding/json"
	"net/http"

	addr "github.com/polera/addr_sh/pkg"
)

func (hc *HandlerContainer) Index(w http.ResponseWriter, r *http.Request) {
	hc.Hits.Inc()
	remoteAddr := addr.GetRemoteHost(r)
	tools := map[string]string{
		"/cidr/v4":   "[GET] Returns CIDR data for a given CIDR formatted address (/cidr/v4/192.168.0.0/24).  Can also be used to find the number of addresses for a CIDR - /cidr/v4/16.",
		"/headers":   "[GET|POST] Returns request headers",
		"/hostnames": "[GET] Performs reverse lookup for remote address (/hostnames) or specified address (/hostnames/8.8.8.8)",
		"/ip":        "[GET] Returns IPv4 address",
	}

	addr := addr.Addr{
		AboutRoute: "/about",
		IP:         remoteAddr,
		Tools:      tools,
	}
	jsonRes, _ := json.Marshal(addr)
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonRes)
}
