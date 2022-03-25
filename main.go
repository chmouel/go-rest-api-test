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
	router, err := reflector.NewRouter()
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Handler:      router,
		Addr:         srvAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Serving on " + srvAddr)
	log.Fatal(srv.ListenAndServe())
}
