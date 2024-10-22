package lorem

import (
	"net/http"

	"github.com/platipy-io/d2s/app"
)

func Index(w http.ResponseWriter, r *http.Request) {
	component := IndexTplt()
	if _, ok := r.Header["Hx-Request"]; !ok {
		component = app.BaseTplt(app.IndexTplt(component))
	}
	component.Render(r.Context(), w)
}
