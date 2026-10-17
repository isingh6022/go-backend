package types

import (
	"net/http"
)

type httpMethod struct {
	value string
}

func (h httpMethod) String() string {
	return h.value
}

var (
	HttpGet     = httpMethod{value: http.MethodGet}
	HttpPost    = httpMethod{value: http.MethodPost}
	HttpPut     = httpMethod{value: http.MethodPut}
	HttpDelete  = httpMethod{value: http.MethodDelete}
	HttpPatch   = httpMethod{value: http.MethodPatch}
	HttpHead    = httpMethod{value: http.MethodHead}
	HttpOptions = httpMethod{value: http.MethodOptions}
)

type Middleware func(http.Handler) http.Handler
