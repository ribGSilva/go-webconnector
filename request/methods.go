package request

import "net/http"

type httpMethod string

const (
	MethodPost    = httpMethod(http.MethodPost)
	MethodGet     = httpMethod(http.MethodGet)
	MethodPatch   = httpMethod(http.MethodPatch)
	MethodPut     = httpMethod(http.MethodPut)
	MethodDelete  = httpMethod(http.MethodDelete)
	MethodHead    = httpMethod(http.MethodHead)
	MethodConnect = httpMethod(http.MethodConnect)
	MethodOptions = httpMethod(http.MethodOptions)
	MethodTrace   = httpMethod(http.MethodTrace)
)
