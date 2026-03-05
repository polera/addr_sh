package server

import (
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

	writeJSON(w, http.StatusOK, about)
}
