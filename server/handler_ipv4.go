package server

import (
	"net/http"

	addr "github.com/polera/addr_sh/pkg"
)

func (hc *HandlerContainer) Getv4cidrData(w http.ResponseWriter, r *http.Request) {
	logger := hc.Log.With().Str("handler", "Getv4cidrData").Logger()
	logger.Debug().Msg("called")

	cidr := r.PathValue("cidr")
	v4cidrData, err := addr.CalculateV4CIDR(cidr)

	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid CIDR value")
		return
	}

	writeJSON(w, http.StatusOK, v4cidrData)
}

func (hc *HandlerContainer) Getv4cidrDetail(w http.ResponseWriter, r *http.Request) {
	logger := hc.Log.With().Str("handler", "Getv4cidrDetail").Logger()
	logger.Debug().Msg("called")
	passedAddress := r.PathValue("address")
	cidr := r.PathValue("cidr")

	v4cidr, err := addr.GetIpv4CIDR(passedAddress, cidr)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid CIDR or address")
		return
	}

	writeJSON(w, http.StatusOK, v4cidr)
}
