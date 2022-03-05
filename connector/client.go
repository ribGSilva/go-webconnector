package connector

import (
	"errors"
	"github.com/ribGSilva/go-webconnector/request"
	"github.com/ribGSilva/go-webconnector/response"
	"net/http"
)

// WebClient is an interface that is able to performs http requests
// the http.Client can be used there
type WebClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Connector contains data to perform the connections to a client
type Connector struct {
	// host to be set in all requests
	host string
	// generalOption has the options to apply to all endpoints
	generalOption []request.Option
	// pathOptions contains the options to each endpoint mapping
	pathOptions map[string][]request.Option
	// webClient contains the client to perform the http request
	webClient WebClient
}

// New creates a new Connector
func New(host string, client WebClient, options ...Option) (Connector, error) {
	c := Connector{
		host:          host,
		generalOption: make([]request.Option, 0),
		pathOptions:   make(map[string][]request.Option),
		webClient:     client,
	}

	for _, o := range options {
		if err := o(&c); err != nil {
			return Connector{}, err
		}
	}

	return c, nil
}

// Option add optional values to the Connector
type Option func(*Connector) error

// WithGeneral adds a general option to the requests
func WithGeneral(o ...request.Option) Option {
	return func(c *Connector) error {
		c.generalOption = append(c.generalOption, o...)
		return nil
	}
}

// WithPath sets a path to the Connector
func WithPath(path string, o ...request.Option) Option {
	return func(c *Connector) error {
		c.pathOptions[path] = o
		return nil
	}
}

// WithPaths sets the paths to the Connector
func WithPaths(po map[string][]request.Option) Option {
	return func(c *Connector) error {
		c.pathOptions = po
		return nil
	}
}

// DoBuild builds the request accordingly to the options and executes it
// the options are applied in the order: general -> pathDefaults -> custom
func (c Connector) DoBuild(path string, responder response.Responder, options ...request.Option) error {
	pathDefaultOption, ok := c.pathOptions[path]
	if !ok {
		return errors.New("connector: unmapped path '" + path + "'")
	}

	reqOptions := make([]request.Option, 0)
	reqOptions = append(reqOptions, c.generalOption...)
	reqOptions = append(reqOptions, pathDefaultOption...)
	reqOptions = append(reqOptions, options...)
	req, err := request.New(c.host, reqOptions...)
	if err != nil {
		return err
	}

	return c.Do(req, responder)
}

// Do should execute the request and triggers the responder
func (c Connector) Do(request *http.Request, responder response.Responder) error {
	if res, err := c.webClient.Do(request); err != nil {
		return err
	} else {
		return responder.Respond(res)
	}
}
