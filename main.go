package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

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
		output, err = ioutil.ReadFile(filepath.Join("fixtures", fixture.Response.File))
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

func main() {
	router := mux.NewRouter()

	if _, err := os.Stat("fixtures/"); os.IsNotExist(err) {
		log.Fatal("Fixtures repository doesn't exist")
	}

	fixtures, err := filepath.Glob("fixtures/*.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if len(fixtures) == 0 {
		log.Fatal("No tests has been found in fixtures/ directory")
	}

	for _, fixture := range fixtures {
		fp, err := ioutil.ReadFile(fixture)
		if err != nil {
			log.Fatal(err)
		}
		var yamlFixture Fixture
		err = yaml.Unmarshal(fp, &yamlFixture)
		if err != nil {
			log.Fatal(err)
		}
		router.HandleFunc(yamlFixture.Headers.Path, func(w http.ResponseWriter, r *http.Request) {
			Handler(w, r, &yamlFixture)
		}).Methods(yamlFixture.Headers.Method)
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Serving on 127.0.0.1:8080....\n")
	log.Fatal(srv.ListenAndServe())
}
