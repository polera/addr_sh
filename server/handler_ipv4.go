package server

import (
	"encoding/json"
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
		w.WriteHeader(406)
		return
	}

	jsonRes, _ := json.Marshal(v4cidrData)
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonRes)

}

func (hc *HandlerContainer) Getv4cidrDetail(w http.ResponseWriter, r *http.Request) {

	logger := hc.Log.With().Str("handler", "Getv4cidrDetail").Logger()
	logger.Debug().Msg("called")
	passedAddress := r.PathValue("address")
	cidr := r.PathValue("cidr")
	cidrAsInt, err := strconv.Atoi(cidr)
	if err != nil || cidrAsInt > 32 {
		w.WriteHeader(406)
		return
	}

	var v4cidr addr.IPv4CIDR
	hc.LookupCache.RLock()
	cacheKey := fmt.Sprintf("%s_%s", passedAddress, cidr)
	cacheValue, hit := hc.LookupCache.Values[cacheKey]
	if hit {
		v4cidr = *cacheValue
	}
	hc.LookupCache.RUnlock()

	if cacheValue == nil {
		v4cidr, err = addr.GetIpv4CIDR(passedAddress, cidr)

		if err != nil {
			w.WriteHeader(406)
			return
		}
		hc.LookupCache.Lock()
		hc.LookupCache.Values[cacheKey] = &v4cidr
		hc.LookupCache.Unlock()

	}

	jsonRes, _ := json.Marshal(v4cidr)
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonRes)
}
