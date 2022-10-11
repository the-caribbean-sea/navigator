package navi

import "net/http"

// Router defines a routes builder
type Router interface {
	MustAuthorized(checker func(
		http.ResponseWriter, *http.Request) (int, *http.Request)) Router
	Just(pattern string, handler func(
		http.ResponseWriter, *http.Request), methods ...string) Router
	Many(pattern string, mappings ...Mapping) Router
}

// router implements the [Router] interface
type router struct {
	checker  func(http.ResponseWriter, *http.Request) (int, *http.Request)
	register func(string, func(http.ResponseWriter, *http.Request))
}

// Navigate creates a new instance of an implementation of the [Router] interface
func Navigate(register func(string, func(http.ResponseWriter, *http.Request))) Router {
	return &router{register: register}
}

// MustAuthorized sets the authorization method for the following [Just] or [Many] registrations
func (r *router) MustAuthorized(checker func(
	http.ResponseWriter, *http.Request) (int, *http.Request)) Router {
	r.checker = checker
	return r
}

// Just registers just one handler mapping to a pattern
func (r *router) Just(pattern string,
	handler func(http.ResponseWriter, *http.Request), methods ...string) Router {
	r.register(pattern, r.compose(Just(handler, methods...)))
	return r
}

// Many registers many handlers mapping to a pattern, the handlers are varied according to different methods
func (r *router) Many(pattern string, mappings ...Mapping) Router {
	r.register(pattern, r.compose(Many(mappings...)))
	return r
}

func (r *router) compose(mapping Mapping) func(http.ResponseWriter, *http.Request) {
	if r.checker == nil {
		return mapping.Compose()
	}
	return mapping.Authorized(r.checker)
}
