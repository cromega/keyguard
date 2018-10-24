package powermux

import (
	"net/http"
	"strings"
)

// Redirect adds a redirect handler for ANY method for this route.
//
// Redirects use either http.StatusPermanentRedirect or http.StatusTemporaryRedirect as their code.
func (r *Route) Redirect(url string, permanent bool) *Route {
	var h http.Handler
	if permanent {
		h = http.RedirectHandler(url, http.StatusPermanentRedirect)
	} else {
		h = http.RedirectHandler(url, http.StatusTemporaryRedirect)
	}
	return r.Any(h)
}

type methodNotAllowedHandler []string

// ServeHTTP responses with a Method Not Allowed and includes an "Allow" header containing the
// valid methods for this route.
func (h methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Sets the Allow header
	w.Header().Add("Allow", strings.Join(h, ", "))
	w.WriteHeader(http.StatusMethodNotAllowed)
}

type defaultOptionsHandler struct {
	methods []string
}

// methodNotAllowed is called internally by Route to generate a 405 handler
func (r *Route) methodNotAllowed() http.Handler {

	// determine what methods ARE supported by this route
	methods := make([]string, 0, 8)

	for method := range r.handlers {
		if method != methodAny && method != notFound {
			methods = append(methods, method)
		}
	}

	// 405 only makes sense if some methods are allowed
	if len(methods) > 0 {
		return methodNotAllowedHandler(methods)
	}

	return nil
}
