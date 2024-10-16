package app

import (
	"net/http"

	"github.com/platipy-io/d2s/internal/telemetry"
)

func Index(w http.ResponseWriter, r *http.Request) {
	_, span := telemetry.Provider.Tracer("").Start(r.Context(), "index")
	defer span.End()
	component := BaseTplt(IndexTplt("John", nil))
	component.Render(r.Context(), w)
}
