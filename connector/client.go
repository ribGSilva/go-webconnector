package connector

import (
	"go-webconnector/response"
	"net/http"
)

type WebClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Connector struct {
	webClint WebClient
}

func (c *Connector) Do(request *http.Request, responder response.Responder) error {
	if res, err := c.webClint.Do(request); err != nil {
		return err
	} else {
		return responder.Respond(res)
	}
}
