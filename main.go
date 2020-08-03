package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

const srvAddr = "0.0.0.0:8080"

type Fixture struct {
	Headers struct {
		Method string
		Path   string
	}
	Response struct {
		Status      int
		File        string
		ContentType string `yaml:"content-type"`
		Output      string
	}
}

func Handler(w http.ResponseWriter, r *http.Request, fixture *Fixture) {
	// vars := mux.Vars(r)
	var output []byte
	var err error

	if fixture.Response.File != "" {
		output, err = ioutil.ReadFile(fixture.Response.File)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	if fixture.Response.Output != "" {
		output = []byte(fixture.Response.Output)
	}

	if fixture.Response.ContentType != "" {
		w.Header().Set("Content-Type", fixture.Response.ContentType)
	}

	w.WriteHeader(fixture.Response.Status)

	_, err = w.Write(output)
	if err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()
	var yamlString string

	if yamlString = os.Getenv("CONFIG"); yamlString == "" {
		log.Fatal("cannot get configuration from environment variable `CONFIG`")
	}
	var yamlFixture Fixture
	reader := bytes.NewReader([]byte(yamlString))
	decoder := yaml.NewDecoder(reader)
	for decoder.Decode(&yamlFixture) == nil {
		router.HandleFunc(yamlFixture.Headers.Path, func(w http.ResponseWriter, r *http.Request) {
			Handler(w, r, &yamlFixture)
		}).Methods(yamlFixture.Headers.Method)
	}

	router.Use(loggingMiddleware)

	srv := &http.Server{
		Handler:      router,
		Addr:         srvAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Serving on " + srvAddr + "....\n")
	log.Fatal(srv.ListenAndServe())
}
