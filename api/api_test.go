package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gathering-gg/parser/config"
	"github.com/stretchr/testify/assert"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server

	// the server URL
	baseURL string
)

// setup sets up a test HTTP server along with a client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	config.Root = server.URL
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestApiNewRequest(t *testing.T) {
}

func TestApiDo(t *testing.T) {
	a := assert.New(t)
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a.Equal("POST", r.Method)
		fmt.Fprint(w, `{}`)
	})

	data := map[string]string{}
	req, err := Upload("/", data)
	a.Nil(err)

	var resp interface{}
	res, err := Do(req, &resp)
	a.Nil(err)
	a.NotNil(res)
}
