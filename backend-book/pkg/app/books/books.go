package books

import (
	"encoding/json"
	"net/http"
)

type Book struct {
	Author string `json:"author"`
	Name   string `json:"name"`
	Year   int    `json:"year"`
}

var books = []Book{
	{
		Author: "Book example",
		Name:   "Author example",
		Year:   2022,
	},
	{
		Author: "Book example",
		Name:   "Author example",
		Year:   2022,
	},
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(books)
	if err != nil {
		panic(err)
	}

	if _, err = w.Write(res); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}
