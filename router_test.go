package lever

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type state struct {
	called     bool
	middleware bool
}

func (s *state) get(w http.ResponseWriter, req *http.Request) {
	s.called = true
}

func TestGenerics(t *testing.T) {
	s := &state{}
	get := func(params []string) (*state, http.HandlerFunc) {
		return s, s.get
	}
	pub := Middlewares[*state]{
		func(s *state, w http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
			s.middleware = true
			return w, req
		},
	}
	router := Router[*state]{
		pub.Get("/", get),
	}

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !s.called {
		t.Fatal("the endpoint was not called correctly")
	}
	if !s.middleware {
		t.Fatal("the middleware was not called correctly")
	}
}
