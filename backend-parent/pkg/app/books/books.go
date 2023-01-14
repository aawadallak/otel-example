package books

import (
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	props := otel.GetTextMapPropagator()
	props.Inject(r.Context(), propagation.HeaderCarrier(w.Header()))

	log.Println("opaaa")

	http.Redirect(w, r, "http://backend-child:5001/books", http.StatusSeeOther)
}
