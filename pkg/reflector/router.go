package reflector

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
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

func handler(w http.ResponseWriter, r *http.Request, fixture *Fixture) {
	// vars := mux.Vars(r)
	var output []byte
	var err error

	if fixture.Response.File != "" {
		output, err = ioutil.ReadFile(fixture.Response.File)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
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

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	var yamlString string

	if yamlString = os.Getenv("CONFIG"); yamlString == "" {
		log.Fatal("cannot get configuration from environment variable `CONFIG`")
	}
	reader := bytes.NewReader([]byte(yamlString))
	decoder := yaml.NewDecoder(reader)

	for {
		var yamlFixture Fixture
		if decoder.Decode(&yamlFixture) != nil {
			break
		}
		log.Printf("Adding: %s", yamlFixture.Headers.Path)
		router.HandleFunc(yamlFixture.Headers.Path, func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, &yamlFixture)
		}).Methods(yamlFixture.Headers.Method)
	}

	router.Use(loggingMiddleware)
	router.NotFoundHandler = router.NewRoute().HandlerFunc(http.NotFound).GetHandler()
	return router
}
