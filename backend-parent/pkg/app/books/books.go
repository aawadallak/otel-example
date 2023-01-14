package books

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	props := otel.GetTextMapPropagator()
	props.Inject(r.Context(), propagation.HeaderCarrier(w.Header()))

	http.Redirect(w, r, "", http.StatusSeeOther)
}
