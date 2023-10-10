package lever

import (
	"net/http"
	"regexp"
	"strings"
)

// Middleware allows wrapping an http handler with functionality before or after
// the request. such as authorization or logging
type Middleware[T any] func(T, http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request)

// Middlewares is a slice of Middleware that allow the generation of route handlers
type Middlewares[T any] []Middleware[T]

// Collapse collapses a slice of Middlewares into a single Middleware, this is
// useful for extending Middlewares into other Middleware slices
func (mid Middlewares[T]) Collapse() Middleware[T] {
	current := mid[:]
	return func(s T, w http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
		for i := len(current) - 1; i >= 0; i-- {
			mid := current[i]
			w, req = mid(s, w, req)
			if w == nil || req == nil {
				return nil, nil
			}
		}
		return w, req
	}
}

// Post generates and returns a Route for use in a Router.
func (m Middlewares[T]) Post(pattern string, handler HandlerFactory[T]) Route[T] {
	return newRoute("POST", pattern, handler, m)
}

// Get generates and returns a Route for use in a Router.
func (m Middlewares[T]) Get(pattern string, handler HandlerFactory[T]) Route[T] {
	return newRoute("GET", pattern, handler, m)
}

// Put generates and returns a Route for use in a Router.
func (m Middlewares[T]) Put(pattern string, handler HandlerFactory[T]) Route[T] {
	return newRoute("PUT", pattern, handler, m)
}

// Del generates and returns a Route for use in a Router.
func (m Middlewares[T]) Del(pattern string, handler HandlerFactory[T]) Route[T] {
	return newRoute("DELETE", pattern, handler, m)
}

// All generates and returns a Route for use in a Router.
func (m Middlewares[T]) All(pattern string, handler HandlerFactory[T]) Route[T] {
	return newRoute("*", pattern, handler, m)
}

// HandlerFactory is a function that accepts matched parameters and returns a state
// and a function for handling the request
type HandlerFactory[T any] func([]string) (T, http.HandlerFunc)

// Route holds all the configuration needed to match url paths to middleware and handlers
type Route[T any] struct {
	method  string
	regex   *regexp.Regexp
	handler HandlerFactory[T]
	mid     Middleware[T]
}

// Router is a convience type for dealing with collections of Routes
type Router[T any] []Route[T]

func newRoute[T any](method, pattern string, handler HandlerFactory[T], mid Middlewares[T]) Route[T] {
	return Route[T]{method, regexp.MustCompile("^" + pattern + "$"), handler, mid.Collapse()}
}

// ServeHTTP allows the Rotuer to be added to an HTTP server and route requests to handlers
func (router Router[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range router {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			continue
		}
		if r.Method != route.method && route.method != "*" {
			allow = append(allow, route.method)
			continue
		}
		state, handler := route.handler(matches[1:])
		w, r = route.mid(state, w, r)
		if w != nil && r != nil {
			handler(w, r)
		}
		return
	}
	if len(allow) == 0 {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Allow", strings.Join(allow, ", "))
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}
