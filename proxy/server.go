package proxy

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer(port int, origin string) {
	r := mux.NewRouter()

	// Stats route
	r.HandleFunc("/cache/stats", StatsHandler).Methods("GET")

	//Clear cashe route
	r.HandleFunc("/cache/clear", ClearCacheHandler).Methods("POST")

	// Others
	r.PathPrefix("/").Handler(ProxyHandler(origin))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Listening on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
