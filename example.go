package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OsagieDG/mlog/service/middleware"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello Everyone!")
	})

	router.HandleFunc("GET /welcome", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "You are welcome to use logStack!")
	})

	mlog := middleware.MLog(
		middleware.LogRequest,
		middleware.LogResponse,
		middleware.RecoverPanic,
	)

	listenAddr := ":6862"
	log.Printf("Server is listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, mlog(router)); err != nil {
		log.Fatal("HTTP server error:", err)
	}
}
