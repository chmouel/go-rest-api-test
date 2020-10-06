package reflector

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var YamlSingle = `---
headers:
  method: GET
  path: /hello/moto.txt
response:
  status: 200
  content-type: text/json
  output: 'hello-moto'
---`

var YamlDoubles = `---
headers:
  method: GET
  path: /1.txt
response:
  status: 200
  content-type: text/json
  output: 'This is one'
---
headers:
  method: GET
  path: /2.txt
response:
  status: 200
  content-type: text/json
  output: 'This is two'
`

var YamlWithFile = `---
headers:
  method: GET
  path: /testme.txt
response:
  status: 200
  content-type: text/json
  file: %s
---`

func newRecord(method, path string) *httptest.ResponseRecorder {
	router, _ := NewRouter()

	req, _ := http.NewRequest(method, path, nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	return res
}

func TestMain(m *testing.M) {
	oldConfig := os.Getenv("CONFIG")
	defer func() {
		fmt.Printf("RESTORING")
		os.Setenv("CONFIG", oldConfig)
	}()

	code := m.Run()
	os.Exit(code)
}

func TestNoConfig(t *testing.T) {
	os.Setenv("CONFIG", "")

	_, err := NewRouter()
	if err == nil {
		t.Error("NewRouter should have bimmed out without CONFIG")
	}
}

func TestNewRouterSingle(t *testing.T) {
	os.Setenv("CONFIG", YamlSingle)

	res := newRecord("GET", "/hello/moto.txt")
	if res.Code != 200 {
		t.Errorf("GET /hello/moto.txt didn't come back with error code 200: %d", res.Code)
	}
	if res.Body.String() != "hello-moto" {
		t.Errorf("GET /hello/moto.txt didn't output 'hello-moto': %s", res.Body.String())
	}

	res = newRecord("GET", "/inexistentasia")
	if res.Code != 404 {
		t.Errorf("GET /inexistentia should have 404: %d", res.Code)
	}

	res = newRecord("PATCH", "/hello/moto.txt")
	if res.Code != 404 {
		t.Errorf("PATCH /inexistentia should have 404: %d", res.Code)
	}

}

func TestNewRouterDoubles(t *testing.T) {
	os.Setenv("CONFIG", YamlDoubles)

	res := newRecord("GET", "/1.txt")
	if res.Code != 200 {
		t.Errorf("GET /hello/moto.txt didn't come back with error code 200")
	}
	if res.Body.String() != "This is one" {
		t.Errorf("GET /1.txt didn't output 'hello-moto': %s", res.Body.String())
	}

	res = newRecord("GET", "/2.txt")
	if res.Code != 200 {
		t.Errorf("GET /2.txt didn't come back with error code 200: %d", res.Code)
	}
	if res.Body.String() != "This is two" {
		t.Errorf("GET /2.txt didn't output 'hello-moto': %s", res.Body.String())
	}

}

func TestResponseFile(t *testing.T) {
	var thestring = "hello moto"
	file, err := ioutil.TempFile("", "go-test.*.json")
	if err != nil {
		t.Error(err)
	}
	_ = ioutil.WriteFile(file.Name(), []byte(thestring), 0644)
	defer os.Remove(file.Name())

	os.Setenv("CONFIG", fmt.Sprintf(YamlWithFile, file.Name()))

	res := newRecord("GET", "/testme.txt")
	if res.Code != 200 {
		t.Errorf("GET /testme.txt didn't come back with error code 200: %d", res.Code)
	}
	if res.Body.String() != thestring {
		t.Errorf("GET /hello/moto.txt didn't output '%s': \"%s\"", thestring, res.Body.String())
	}

}
