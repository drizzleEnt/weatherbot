package routes

import (
	"fmt"
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("webhook")
		greatings := "webhook hello"
		w.Write([]byte(greatings))
	})
	return r
}
