package app

import (
	"net/http"

	"github.com/IxDay/templ-exp/internal/telemetry"
)

func Index(w http.ResponseWriter, r *http.Request) {
	_, span := telemetry.Provider.Tracer("server").Start(r.Context(), "index")
	defer span.End()
	component := BaseTplt(IndexTplt("John", nil))
	component.Render(r.Context(), w)
}
