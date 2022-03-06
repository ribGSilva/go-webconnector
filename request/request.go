// Builder package brings facilities to build http.Request.
// it brings re-utilizable codes with options with the most usual necessities

package request

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	headerContentType = "Content-Type"
)

// Builder carries all the data necessary to execute a http request
type Builder struct {
	// ctx context for the Builder
	ctx context.Context
	// method is the http GET, POST...
	method httpMethod
	// protocol is the protocol for the Builder
	// Example:
	// 		http
	// 		https
	protocol string
	// host is the host of the Builder
	// Example:
	// 		my.host.com
	host string
	// path is the path for the Builder
	// Example:
	//		/my/path
	//		/:myParam
	path string
	// params has the params to bind in the path
	params map[string]string
	// headers has the headers of the Builder
	headers map[string][]string
	// queries has the queries of the Builder
	queries map[string][]string
	// body has the body for the Builder
	body io.Reader
}

// New creates a new Builder
// By default is a MethodGet Builder
// By default the protocol is http
// Example:
//		func buildReq(ctx context.Context, id string, body interface{}) {
//			req, err := New("my.host.com",
//				WithContext(ctx),
//				WithMethod(MethodPatch), // by default is GET
//				WithProtocol("https"), // by default is http
//				WithPath("/path/:id"),
//				WithParam("id", id),
//				WithQuery("myQuery", "someValue"),
//				WithHeader("Authorization", "myauth"),
//				WithJson(body),
//			)
//		}
func New(host string, options ...Option) (*http.Request, error) {
	r := Builder{
		method:   MethodGet,
		host:     host,
		protocol: "http",
		params:   make(map[string]string),
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

func build(r Builder) (*http.Request, error) {
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

	p := r.path
	for k, v := range r.params {
		p = strings.ReplaceAll(p, ":"+k, v)
	}

	url := fmt.Sprintf("%s://%s%s%s", r.protocol, r.host, p, q)

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

// Option add optional values to the Builder
type Option func(*Builder) error

// WithContext specify the context for the Builder
func WithContext(ctx context.Context) Option {
	return func(r *Builder) error {
		r.ctx = ctx
		return nil
	}
}

// WithProtocol specify the protocol for the Builder
func WithProtocol(protocol string) Option {
	return func(r *Builder) error {
		r.protocol = protocol
		return nil
	}
}

// WithMethod specify the http method for the Builder
func WithMethod(method httpMethod) Option {
	return func(r *Builder) error {
		r.method = method
		return nil
	}
}

// WithPath sets the path
// To set path params, use :{value}
// Example:
// 			...
// 			WithPath("/:userId/address/:addId")
//			WithParam("userId", "123")
//			WithParam("addId", "2")
// 			...
func WithPath(path string) Option {
	return func(r *Builder) error {
		r.path = path
		return nil
	}
}

// WithParam adds a param bind
func WithParam(key string, value interface{}) Option {
	return func(r *Builder) error {
		r.params[key] = fmt.Sprint(value)
		return nil
	}
}

// WithParams sets the params
func WithParams(params map[string]interface{}) Option {
	return func(r *Builder) error {
		for k, v := range params {
			r.params[k] = fmt.Sprint(v)
		}
		return nil
	}
}

// WithHeader adds to the header a value
// The header name will always be first letter Upper
// Example:
// 			...
// 			WithHeader("authoRIZATION", "someHASH")
// 			WithHeader("content-tyPE", "someContent")
// 			...
//     this will end up as a header:
//			Authorization: someHASH
//			Content-Type:  someContent

func WithHeader(key string, value interface{}) Option {
	return func(r *Builder) error {
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
	return func(r *Builder) error {
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

// WithQuery adds query param to the Builder
func WithQuery(key string, value interface{}) Option {
	return func(r *Builder) error {
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
	return func(r *Builder) error {
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
	return func(r *Builder) error {
		r.body = body
		return nil
	}
}

// WithString sets the body as a string
func WithString(body string) Option {
	return func(r *Builder) error {
		r.body = bytes.NewBufferString(body)
		return nil
	}
}

// WithJson sets the body as a json
// This method already sets the Content-Type header as application/json
func WithJson(body interface{}) Option {
	return func(r *Builder) error {
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
// This method already sets the Content-Type header as application/xml
func WithXml(body interface{}) Option {
	return func(r *Builder) error {
		if b, err := xml.Marshal(body); err != nil {
			return err
		} else {
			r.headers[headerContentType] = []string{"application/xml"}
			r.body = bytes.NewBuffer(b)
		}
		return nil
	}
}
