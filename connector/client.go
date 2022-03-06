// connector is a package to centralize the request configs
// it can re-utilize the common configs of the requests

package connector

import (
	"github.com/ribGSilva/go-webconnector/request"
	"net/http"
)

// WebClient is an interface that is able to performs http requests
// the http.Client can be used there
type WebClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Responder is an interface that is capable of handle a http.Response
// made to be used with the response package
type Responder interface {
	Respond(*http.Response) error
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
// Example:
// 		func execRequests() error {
//			getAllPath := "/"
//			getPath := "/:id"
//			postPath := "/"
//
//			c, err := New("my.host.com",
//				http.DefaultClient,
//				WithGeneral(request.WithHeader("My-Header", "myHeaderKey")), // apply in all requests
//				WithPath(getPath), // get
//				WithPath(postPath, request.WithMethod(request.MethodPost)), // post
//			)
//			if err != nil {
//				return err
//			}
//
//			//get
//			getBody := struct {
//				Name string `json:"name"`
//			}{}
//			responder, err := response.NewResponder(response.ForJson(200, getBody))
//			if err != nil {
//				return err
//			}
//			err = c.DoBuild(getPath, &responder, request.WithParam("id", "123"))
//			if err != nil {
//				return err
//			}
//			fmt.Printf("%+v \n", getBody)
//
//			//post
//			postBody := struct {
//				Name string `json:"name"`
//			}{Name: "my name"}
//			responder, err = response.NewResponder(
//				response.ForStatus(201),
//				response.ForDefault(func(r response.Response) error {
//					return errors.New("error creating")
//				}),
//			)
//			if err != nil {
//				return err
//			}
//			err = c.DoBuild(getPath, &responder, request.WithJson(postBody))
//			if err != nil {
//				return err
//			}
//			fmt.Println("created")
//
//			//get all
//			getAll := make([]struct {
//				Name string `json:"name"`
//			}, 0)
//
//			responder, err = response.NewResponder(response.ForJson(200, getAll))
//			if err != nil {
//				return err
//			}
//			err = c.DoBuild(getAllPath, &responder, request.WithQuery("page", "3"))
//			if err != nil {
//				return err
//			}
//			fmt.Printf("%+v \n", getAll)
//
//			return nil
//		}
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
func (c Connector) DoBuild(path string, responder Responder, options ...request.Option) error {

	reqOptions := []request.Option{request.WithPath(path)}
	reqOptions = append(reqOptions, c.generalOption...)

	pathDefaultOption, ok := c.pathOptions[path]
	if ok {
		reqOptions = append(reqOptions, pathDefaultOption...)
	}

	reqOptions = append(reqOptions, options...)

	req, err := request.New(c.host, reqOptions...)
	if err != nil {
		return err
	}

	return c.Do(req, responder)
}

// Do should execute the request and triggers the responder
func (c Connector) Do(request *http.Request, responder Responder) error {
	if res, err := c.webClient.Do(request); err != nil {
		return err
	} else {
		return responder.Respond(res)
	}
}
