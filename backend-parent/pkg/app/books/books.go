package books

import (
	"io"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var backendURL = os.Getenv("URL_BACKEND_CHILD")

func GetBooks(w http.ResponseWriter, r *http.Request) {
	props := otel.GetTextMapPropagator()
	props.Inject(r.Context(), propagation.HeaderCarrier(r.Header))

	req, err := http.NewRequest(http.MethodGet, backendURL+"/books", nil)
	if err != nil {
		log.Println("http.NewRequest() returnd error:", err.Error())
		return
	}

	req.Header = r.Header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("client.Do() returnd error:", err.Error())
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
