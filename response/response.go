package response

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Response holds data of the http response
type Response struct {
	// HttpResponse the original response
	HttpResponse *http.Response
}

// Responder holds data about which function it should respond for reach http status
type Responder struct {
	// responders has the map for the status:func handler
	responders map[int]Func
	// defResponder has the default func handler
	defResponder Func
}

// Func handles a response
type Func func(Response) error

// Respond handles how to proceed with a http.Response
// I search in its configuration and calls the specific function for the http status
// If not mapped, it will call a default responder function (if set)
// And if in some point has an error, the method will return the error
func (r *Responder) Respond(res *http.Response) error {
	if res == nil {
		return nil
	}

	response := Response{
		HttpResponse: res,
	}

	f, ok := r.responders[res.StatusCode]
	if ok {
		return f(response)
	} else if r.defResponder != nil {
		return r.defResponder(response)
	}
	return nil
}

// NewResponder creates a new Responder
func NewResponder(options ...Option) (Responder, error) {
	r := Responder{
		responders:   make(map[int]Func),
		defResponder: nil,
	}

	for _, o := range options {
		if err := o(&r); err != nil {
			return Responder{}, err
		}
	}

	return r, nil

}

// Option add optional values to the request
type Option func(*Responder) error

// For specify function to handle a specific status
func For(status int, f Func) Option {
	return func(r *Responder) error {
		r.responders[status] = f
		return nil
	}
}

// ForDefault specify function to handle non mapped status
func ForDefault(f Func) Option {
	return func(r *Responder) error {
		r.defResponder = f
		return nil
	}
}

// ForStatus specify that for that status, the application will do nothing
func ForStatus(status int) Option {
	return func(r *Responder) error {
		r.responders[status] = func(response Response) error {
			return nil
		}
		return nil
	}
}

// ForString specify function to handle a specific status returning a parsed string
func ForString(status int, resp *string) Option {
	return func(r *Responder) error {
		r.responders[status] = func(response Response) error {
			if data, err := ioutil.ReadAll(response.HttpResponse.Body); err != nil {
				return err
			} else {
				*resp = string(data)
				return nil
			}
		}
		return nil
	}
}

// ForJson specify function to handle a specific status returning a parsed json
func ForJson(status int, int interface{}) Option {
	return func(r *Responder) error {
		r.responders[status] = func(response Response) error {
			if data, err := ioutil.ReadAll(response.HttpResponse.Body); err != nil {
				return err
			} else {
				return json.Unmarshal(data, int)
			}
		}
		return nil
	}
}

// ForXml specify function to handle a specific status returning a parsed xml
func ForXml(status int, int interface{}) Option {
	return func(r *Responder) error {
		r.responders[status] = func(response Response) error {
			if data, err := ioutil.ReadAll(response.HttpResponse.Body); err != nil {
				return err
			} else {
				return xml.Unmarshal(data, int)
			}
		}
		return nil
	}
}
