package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chmouel/go-rest-api-test/pkg/reflector"
)

const srvAddr = "0.0.0.0:8080"

func main() {
	router := reflector.NewRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         srvAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Serving on " + srvAddr + "....\n")
	log.Fatal(srv.ListenAndServe())
}
