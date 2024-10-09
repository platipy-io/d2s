package lorem

import (
	"net/http"

	"github.com/IxDay/templ-exp/app"
)

func Index(w http.ResponseWriter, r *http.Request) {
	component := IndexTplt()
	if _, ok := r.Header["Hx-Request"]; !ok {
		component = app.BaseTplt(app.IndexTplt("John", component))
	}
	component.Render(r.Context(), w)
}
