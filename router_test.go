package lever

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type state struct {
	called bool
}

func (s *state) get(w http.ResponseWriter, req *http.Request) {
	s.called = true
}

func TestGenerics(t *testing.T) {
	s := &state{}
	get := func(params []string) (*state, http.HandlerFunc) {
		return s, s.get
	}
	pub := Middlewares[*state]{}
	router := Router[*state]{
		pub.Get("/", get),
	}

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !s.called {
		t.Fatal("the endpoint was not called correctly")
	}
}
