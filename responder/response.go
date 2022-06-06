// responder package brings facilities to handle responses og http.Response
// it brings parsed responses, accordingly with the status

package responder

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// Response holds data of the http responder
type Response struct {
	// Status has the http status of the request
	Status int
	// Body holds the body parsed
	Body any
	// HttpResponse the original responder
	HttpResponse *http.Response
	// Err holds errors if it had one
	Err error
}

// BodyParser handles the body parsing
type BodyParser func(io.ReadCloser) (any, error)

// Option add optional values to the request
type Option func(*Responder)

// Responder holds data about which function it should respond for reach http status
type Responder struct {
	// responders has the map for the status:func handler
	responders map[int]BodyParser
	// defResponder has the default func handler
	defResponder BodyParser
}

// New creates a new Responder
// Example:
// 		func handleResponse(resp *http.Response) error {
//			responder := NewResponder(
//				Status(http.StatusNotFound), // Does nothing
//				For(http.StatusOK, func(body io.ReadCloser) (any, error) {
//					var b myStruct
//					err := json.NewDecoder(body).Decode(&b)
//					if err != nil {
//						return nil, err
//					}
//					return b, nil
//				}),
//				Default(func(responder io.ReadCloser) (any, error) {
//					return nil, errors.New("responder: not mapped status")
//				}),
//			)
//
//			return responder.Respond(resp)
//		}
func New(options ...Option) *Responder {
	r := &Responder{
		responders:   make(map[int]BodyParser),
		defResponder: nil,
	}

	for _, o := range options {
		o(r)
	}

	return r

}

var ErrNoHttpResponse = errors.New("connector/responder: no http response")
var ErrNoResponseHandler = errors.New("connector/responder: no response handler")

// Respond handles how to proceed with a http.Response
// I search in its configuration and calls the specific function for the http status
// If not mapped, it will call a default responder function (if set)
// And if in some point has an error, the method will return the error
func (r *Responder) Respond(res *http.Response) *Response {
	if res == nil {
		return &Response{
			Err: ErrNoHttpResponse,
		}
	}

	f, ok := r.responders[res.StatusCode]
	if ok {
		b, err := f(res.Body)
		return &Response{
			Status:       res.StatusCode,
			HttpResponse: res,
			Body:         b,
			Err:          err,
		}
	} else if r.defResponder != nil {
		b, err := r.defResponder(res.Body)
		return &Response{
			Status:       res.StatusCode,
			HttpResponse: res,
			Body:         b,
			Err:          err,
		}
	}
	return &Response{
		Err: ErrNoResponseHandler,
	}
}

// For specify function to handle a specific status
func For(status int, f BodyParser) Option {
	return func(r *Responder) {
		r.responders[status] = f
	}
}

// Default specify function to handle non mapped status
func Default(f BodyParser) Option {
	return func(r *Responder) {
		r.defResponder = f
	}
}

// Status specify that for that status, the application will do nothing
func Status(status int) Option {
	return func(r *Responder) {
		r.responders[status] = func(io.ReadCloser) (any, error) {
			return nil, nil
		}
	}
}

// String specify function to handle a specific status returning a parsed string
func String(status int) Option {
	return func(r *Responder) {
		r.responders[status] = func(body io.ReadCloser) (any, error) {
			if data, err := ioutil.ReadAll(body); err != nil {
				return nil, err
			} else {
				return string(data), nil
			}
		}
	}
}
