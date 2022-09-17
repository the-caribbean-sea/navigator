# Introduction

This is a `simple` and (*of course as it is simple, hence*) `lightweight` approach to gracefully compose http server routes.

> Call me router-ish~~mael~~

# Example
```go
package main

import (
	"net/http"

	navi "github.com/the-caribbean-sea/navigator"
)

func authorize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Token", "I-Am-The-Token")
}

func authorized(w http.ResponseWriter, r *http.Request) int {
	if r.Header.Get("Token") != "I-Am-The-Token" {
		return http.StatusUnauthorized
	}
	return 0
}

func list(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[\"hello\", \"world\"]"))
}

func gen1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("1"))
}

func gen2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("2"))
}

func privileged(w http.ResponseWriter, r *http.Request) int {
	if r.Header.Get("Token") != "I-Am-The-Privileged" {
		return http.StatusForbidden
	}
	return 0
}

func admin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am the admin"))
}

func main() {
	navi.NewRoute(http.HandleFunc).
		Just("/auth", authorize, http.MethodPost).
		// following handlers will have the [authorized] method to check before handling
		MustAuthorized(authorized).
		Just("/user/list", list, http.MethodGet, http.MethodPost).
		Many("/number/gen",
			navi.Just(gen1, http.MethodGet),
			navi.Just(gen2, http.MethodPost)).
		// following handlers will have the [privileged] method to check before handling
		MustAuthorized(privileged).
		Just("/admin", admin, http.MethodPost)

	http.ListenAndServe("0.0.0.0:80", nil)
}
```
