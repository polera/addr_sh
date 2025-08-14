package server

import (
	"encoding/json"
	"net/http"

	addr "github.com/polera/addr_sh/pkg"
)

func (hc *HandlerContainer) GetAbout(w http.ResponseWriter, r *http.Request) {
	hc.Hits.Inc()
	logger := hc.Log.With().Str("handler", "GetAbout").Logger()
	logger.Debug().Msg("called about")

	about := &addr.About{
		Text:   "JSON API for various HTTP/networking related tools.",
		Email:  "james@uncryptic.com",
		GitHub: "https://github.com/polera",
	}

	jsonRes, _ := json.Marshal(about)
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonRes)
}
