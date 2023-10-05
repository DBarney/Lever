package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type middleware func(*state, http.HandlerFunc) http.HandlerFunc
type middlewares []middleware

var NotFound = http.NotFound

func (mid middlewares) Collapse() middleware {
	reversed := middlewares{}
	for i := len(mid) - 1; i >= 0; i-- {
		reversed = append(reversed, mid[i])
	}
	return func(s *state, next http.HandlerFunc) http.HandlerFunc {
		for _, mid := range reversed {
			next = mid(s, next)
		}
		return next
	}
}

type handlerFactory func([]string) (*state, http.HandlerFunc)
type route struct {
	method  string
	regex   *regexp.Regexp
	handler handlerFactory
	mid     middleware
}
type Router []route

func abort(matches []string) (*state, http.HandlerFunc) {
	return nil, nil
}

func (m middlewares) post(pattern string, handler handlerFactory) route {
	return newRoute("POST", pattern, handler, m)
}
func (m middlewares) get(pattern string, handler handlerFactory) route {
	return newRoute("GET", pattern, handler, m)
}
func (m middlewares) put(pattern string, handler handlerFactory) route {
	return newRoute("PUT", pattern, handler, m)
}
func (m middlewares) del(pattern string, handler handlerFactory) route {
	return newRoute("DELETE", pattern, handler, m)
}
func (m middlewares) all(pattern string, handler handlerFactory) route {
	return newRoute("*", pattern, handler, m)
}

func newRoute(method, pattern string, handler handlerFactory, mid middlewares) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler, mid.Collapse()}
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range router {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			continue
		}
		// html forms can't do DELETE or PUT requests. :P
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			m := r.PostFormValue("_method")
			if m != "" {
				r.Method = m
				delete(r.PostForm, "_method")
			}
		}
		if r.Method != route.method && route.method != "*" {
			allow = append(allow, route.method)
			continue
		}
		state, handler := route.handler(matches[1:])
		route.mid(state, handler)(w, r)
		return
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		fmt.Println(r.URL.Path, allow, r.Method)
		return
	}
	NotFound(w, r)
}
