package navi

import "net/http"

// Mapping holds underlying [ method -> handler ] mapping, and composes new handlers
type Mapping interface {
	// Compose composes the new handler
	Compose() func(w http.ResponseWriter, r *http.Request)
	// Authorized composes a new handler with the given authorization method
	Authorized(authorization func(
		w http.ResponseWriter, r *http.Request) int) func(w http.ResponseWriter, r *http.Request)
	fill(m mapping)
}

type mapping map[string]func(w http.ResponseWriter, r *http.Request)

// Methods holds methods waiting for a handler
type Methods []string

// Many creates a [Mapping] out of the given set of [Mapping]s
func Many(mappings ...Mapping) Mapping {
	m := mapping{}
	for _, other := range mappings {
		other.fill(m)
	}
	return m
}

// Just creates a [Mapping] for some specific methods
func Just(handler func(w http.ResponseWriter,
	r *http.Request), methods ...string) Mapping {
	return Methods(methods).employ(handler)
}

func (m mapping) Compose() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if handle, ok := m[r.Method]; ok {
			handle(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m mapping) Authorized(authorization func(
	w http.ResponseWriter, r *http.Request) int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if status := authorization(w, r); status != 0 {
			w.WriteHeader(status)
			return
		}

		if handle, ok := m[r.Method]; ok {
			handle(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m mapping) fill(target mapping) {
	for k, v := range m {
		target[k] = v
	}
}

func (methods Methods) employ(
	handler func(w http.ResponseWriter, r *http.Request)) Mapping {
	m := mapping{}
	for _, method := range methods {
		m[method] = handler
	}
	return m
}
