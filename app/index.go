package app

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	component := BaseTplt(IndexTplt("John", nil))
	component.Render(r.Context(), w)
}
