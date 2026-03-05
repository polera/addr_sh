package server

import (
	"fmt"
	"net/http"
	"strconv"

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

func (hc *HandlerContainer) SplitV4CIDR(w http.ResponseWriter, r *http.Request) {
	logger := hc.Log.With().Str("handler", "SplitV4CIDR").Logger()
	logger.Debug().Msg("called")

	network := r.PathValue("network")
	prefix := r.PathValue("prefix")
	countStr := r.PathValue("count")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "count must be an integer")
		return
	}

	cidr := fmt.Sprintf("%s/%s", network, prefix)
	result, err := addr.SplitCIDR(cidr, count)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
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
