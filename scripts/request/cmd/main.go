package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

var backendURL = os.Getenv("URL_BACKEND_PARENT")

func main() {
	req, err := http.NewRequest(http.MethodGet, backendURL+"/books", nil)
	if err != nil {
		log.Println("http.NewRequest() returnd error:", err.Error())
		return
	}

	for i := 0; ; i++ {
		f := func() {
			defer time.Sleep(time.Second)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("client.Do() returnd error:", err.Error())
				return
			}

			defer res.Body.Close()

			if res.StatusCode != 200 {
				log.Println("fail to request backend")
				return
			}

			log.Println("successfully requested to backend")
		}

		f()
	}
}
