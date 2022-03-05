package request

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

const (
	headerContentType = "Content-Type"
)

// request carries all the data necessary to execute a request
type request struct {
	ctx      context.Context
	method   httpMethod
	protocol string
	host     string
	path     string
	headers  map[string][]string
	queries  map[string][]string
	body     io.Reader
}

// New creates a new request
// By default is a MethodGet request
// By default the protocol is http
func New(host string, options ...Option) (*http.Request, error) {
	r := request{
		method:   MethodGet,
		host:     host,
		protocol: "http",
		headers:  make(map[string][]string),
		queries:  make(map[string][]string),
	}
	for _, o := range options {
		if err := o(&r); err != nil {
			return nil, err
		}
	}

	return build(r)
}

func build(r request) (*http.Request, error) {
	q := ""

	for k, v := range r.queries {

		for _, qv := range v {
			if len(q) == 0 {
				q = "?"
			} else {
				q = q + "&"
			}

			q = q + k + "=" + qv
		}
	}

	url := fmt.Sprintf("%s://%s%s%s", r.protocol, r.host, r.path, q)

	req := new(http.Request)
	if r.ctx != nil {
		var err error
		if req, err = http.NewRequestWithContext(r.ctx, string(r.method), url, r.body); err != nil {
			return nil, err
		}
	} else {
		var err error
		if req, err = http.NewRequest(string(r.method), url, r.body); err != nil {
			return nil, err
		}
	}

	for k, v := range r.headers {
		for _, hv := range v {
			req.Header.Add(k, hv)
		}
	}

	return req, nil
}

// Option add optional values to the request
type Option func(*request) error

// WithContext specify the context for the request
func WithContext(ctx context.Context) Option {
	return func(r *request) error {
		r.ctx = ctx
		return nil
	}
}

// WithProtocol specify the protocol for the request
func WithProtocol(protocol string) Option {
	return func(r *request) error {
		r.protocol = protocol
		return nil
	}
}

// WithMethod specify the http method for the request
func WithMethod(method httpMethod) Option {
	return func(r *request) error {
		r.method = method
		return nil
	}
}

// WithPath sets the path
func WithPath(path string) Option {
	return func(r *request) error {
		r.path = path
		return nil
	}
}

// WithHeader adds to the header a value
func WithHeader(key string, value interface{}) Option {
	return func(r *request) error {
		if _, ok := r.headers[key]; ok {
			r.headers[key] = append(r.headers[key], fmt.Sprint(value))
		} else {
			r.headers[key] = []string{fmt.Sprint(value)}
		}
		return nil
	}
}

// WithHeaders sets the headers
func WithHeaders(headers map[string][]interface{}) Option {
	return func(r *request) error {
		for k, v := range headers {
			for _, hv := range v {
				if _, ok := r.headers[k]; ok {
					r.headers[k] = append(r.headers[k], fmt.Sprint(hv))
				} else {
					r.headers[k] = []string{fmt.Sprint(hv)}
				}
			}
		}
		return nil
	}
}

// WithQuery adds query param to the request
func WithQuery(key string, value interface{}) Option {
	return func(r *request) error {
		if _, ok := r.queries[key]; ok {
			r.queries[key] = append(r.queries[key], fmt.Sprint(value))
		} else {
			r.queries[key] = []string{fmt.Sprint(value)}
		}
		return nil
	}
}

// WithQueries sets the query params
func WithQueries(queries map[string][]interface{}) Option {
	return func(r *request) error {
		for k, v := range queries {
			for _, qv := range v {
				if _, ok := r.queries[k]; ok {
					r.queries[k] = append(r.queries[k], fmt.Sprint(qv))
				} else {
					r.queries[k] = []string{fmt.Sprint(qv)}
				}
			}
		}
		return nil
	}
}

// WithBody sets the body
func WithBody(body io.Reader) Option {
	return func(r *request) error {
		r.body = body
		return nil
	}
}

// WithString sets the body as a string
func WithString(body string) Option {
	return func(r *request) error {
		r.body = bytes.NewBufferString(body)
		return nil
	}
}

// WithJson sets the body as a json
func WithJson(body interface{}) Option {
	return func(r *request) error {
		if b, err := json.Marshal(body); err != nil {
			return err
		} else {
			r.headers[headerContentType] = []string{"application/json"}
			r.body = bytes.NewBuffer(b)
		}
		return nil
	}
}

// WithXml sets the body as a xml
func WithXml(body interface{}) Option {
	return func(r *request) error {
		if b, err := xml.Marshal(body); err != nil {
			return err
		} else {
			r.headers[headerContentType] = []string{"application/xml"}
			r.body = bytes.NewBuffer(b)
		}
		return nil
	}
}
