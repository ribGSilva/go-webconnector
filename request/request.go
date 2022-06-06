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

// Builder carries all the data necessary to execute a http request
type Builder struct {
	// Context for the Builder
	Context context.Context
	// method is the http GET, POST...
	Method string
	// Host is the host of the Builder
	// Example:
	// 		http://my.host.com
	Host string
	// Path is the path for the Builder
	// Example:
	//		/my/path
	//		/:myParam
	Path string
	// Params has the params to bind in the path
	Params map[string]string
	// Headers has the headers of the Builder
	Headers http.Header
	// Queries has the queries of the Builder
	Queries map[string][]string
	// Encoder has the encoder for the Body
	Encoder EncoderFunc
	// Body has the body for the Builder
	Body any
}

//EncoderFunc encodes the Body
type EncoderFunc func(any) ([]byte, error)

// Option add optional values to the Builder
type Option func(*Builder)

// New creates a new *http.Request
// Example:
//		func buildReq(ctx context.Context, id string, body interface{}) {
//			req, err := New("http://my.host.com",
//				Context(ctx),
//				Method(MethodPatch), // by default is GET
//				Path("/path/:id"),
//				Param("id", id),
//				Query("myQuery", "someValue"),
//				Header("Authorization", "myauth"),
//				JSON(body),
//			)
//		}
func New(host string, options ...Option) (*http.Request, error) {
	return NewBuilder(host, options...).Build()
}

// NewBuilder a new Builder
// Example:
//		func reqBuilder(ctx context.Context, id string, body interface{}) {
//			builder := NewBuilder("http://my.host.com",
//				Context(ctx),
//				Method(MethodPatch), // by default is GET
//				Path("/path/:id"),
//				Param("id", id),
//				Query("myQuery", "someValue"),
//				Header("Authorization", "myauth"),
//				Body(body),
//			)
//		}
func NewBuilder(host string, options ...Option) *Builder {
	r := Builder{
		Context: context.Background(),
		Method:  http.MethodGet,
		Host:    host,
		Params:  make(map[string]string),
		Headers: make(http.Header),
		Queries: make(map[string][]string),
		Encoder: json.Marshal,
	}
	for _, o := range options {
		o(&r)
	}

	return &r
}

func (r *Builder) Build() (*http.Request, error) {
	q := ""

	for k, v := range r.Queries {

		for _, qv := range v {
			if len(q) == 0 {
				q = "?"
			} else {
				q = q + "&"
			}

			q = q + k + "=" + qv
		}
	}

	p := r.Path
	for k, v := range r.Params {
		p = strings.ReplaceAll(p, ":"+k, v)
	}

	url := fmt.Sprintf("%s%s%s", r.Host, p, q)

	var body io.Reader
	if r.Body != nil {
		b, err := r.Encoder(r.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(r.Context, r.Method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header = r.Headers

	return req, nil
}

// Context specify the context for the Builder
func Context(ctx context.Context) Option {
	return func(r *Builder) {
		r.Context = ctx
	}
}

// Method specify the http method for the Builder
func Method(method string) Option {
	return func(r *Builder) {
		r.Method = method
	}
}

// Path sets the path
// To set path params, use :{value}
// Example:
// 			...
// 			Path("/:userId/address/:addId")
//			Param("userId", "123")
//			Param("addId", "2")
// 			...
func Path(path string) Option {
	return func(r *Builder) {
		r.Path = path
	}
}

// Param adds a param bind
func Param(key string, value interface{}) Option {
	return func(r *Builder) {
		r.Params[key] = fmt.Sprint(value)
	}
}

// Params sets the params
func Params(params map[string]interface{}) Option {
	return func(r *Builder) {
		for k, v := range params {
			r.Params[k] = fmt.Sprint(v)
		}
	}
}

// Header adds to the header a value
// The header name will always be first letter Upper
// Example:
// 			...
// 			WithHeader("authoRIZATION", "someHASH")
// 			WithHeader("content-tyPE", "someContent")
// 			...
//     this will end up as a header:
//			Authorization: someHASH
//			Content-Type:  someContent

func Header(key string, value interface{}) Option {
	return func(r *Builder) {
		r.Headers.Add(key, fmt.Sprint(value))
	}
}

// Headers sets the headers
func Headers(headers http.Header) Option {
	return func(r *Builder) {
		r.Headers = headers
	}
}

// Query adds query param to the Builder
func Query(key string, value interface{}) Option {
	return func(r *Builder) {
		if _, ok := r.Queries[key]; ok {
			r.Queries[key] = append(r.Queries[key], fmt.Sprint(value))
		} else {
			r.Queries[key] = []string{fmt.Sprint(value)}
		}
	}
}

// Queries sets the query params
func Queries(queries map[string][]interface{}) Option {
	return func(r *Builder) {
		for k, v := range queries {
			for _, qv := range v {
				if _, ok := r.Queries[k]; ok {
					r.Queries[k] = append(r.Queries[k], fmt.Sprint(qv))
				} else {
					r.Queries[k] = []string{fmt.Sprint(qv)}
				}
			}
		}
	}
}

// Encoder sets the encoder
func Encoder(f EncoderFunc) Option {
	return func(r *Builder) {
		r.Encoder = f
	}
}

// Body sets the body
func Body(body any) Option {
	return func(r *Builder) {
		r.Body = body
	}
}

// String sets the body as a string
func String(body string) Option {
	return func(r *Builder) {
		r.Body = bytes.NewBufferString(body)
		r.Encoder = func(any) ([]byte, error) {
			return []byte(body), nil
		}
	}
}

// JSON sets the body as a json
// This method already sets the Content-Type header as application/json
func JSON(body interface{}) Option {
	return func(r *Builder) {
		r.Body = body
		r.Encoder = json.Marshal
		r.Headers.Add("Content-Type", "application/json")
	}
}

// XML sets the body as a xml
// This method already sets the Content-Type header as application/xml
func XML(body interface{}) Option {
	return func(r *Builder) {
		r.Body = body
		r.Encoder = xml.Marshal
		r.Headers.Add("Content-Type", "application/xml")
	}
}
