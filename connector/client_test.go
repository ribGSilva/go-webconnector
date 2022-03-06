package connector

import (
	"errors"
	"github.com/ribGSilva/go-webconnector/request"
	"net/http"
	"testing"
)

const host = "defaultHost"

func TestNew(t *testing.T) {
	_, err := New(host, &mockWebClient{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestNewPath(t *testing.T) {
	reqGet := "/get-endpoint"
	c, err := New(host, &mockWebClient{
		expectedUrl:    "http://" + host + reqGet,
		expectedMethod: "GET",
	},
		WithPath(reqGet))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = c.DoBuild(reqGet, &mockResponder{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestNewPaths(t *testing.T) {
	reqGet := "/get-endpoint"
	c, err := New(host, &mockWebClient{
		expectedUrl:    "http://" + host + reqGet,
		expectedMethod: "GET",
	},
		WithPaths(make(map[string][]request.Option)))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = c.DoBuild(reqGet, &mockResponder{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestNewGeneralPath(t *testing.T) {
	reqGet := "/get-endpoint"
	c, err := New(host, &mockWebClient{
		expectedUrl:    "https://" + host + reqGet,
		expectedMethod: "GET",
	},
		WithGeneral(request.WithProtocol("https")),
		WithPath(reqGet))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = c.DoBuild(reqGet, &mockResponder{})
}

func TestNewGeneralPathCustom(t *testing.T) {
	reqGet := "/get-endpoint"
	c, err := New(host, &mockWebClient{
		expectedUrl:    "http://" + host + reqGet + "?myQuery=queryValue",
		expectedMethod: "GET",
	},
		WithPath(reqGet))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = c.DoBuild(reqGet, &mockResponder{}, request.WithQuery("myQuery", "queryValue"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestNewErr(t *testing.T) {
	_, err := New(host, &mockWebClient{}, func(c *Connector) error {
		return errors.New("mocked error")
	})
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

func TestNewErrBuild(t *testing.T) {
	reqGet := "/get-endpoint"
	c, err := New(host, &mockWebClient{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = c.DoBuild(reqGet, &mockResponder{}, func(r *request.Builder) error {
		return errors.New("mocked error")
	})
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

func TestNewErrDo(t *testing.T) {
	reqGet := "/get-endpoint"
	c, err := New(host, &mockWebClient{err: errors.New("mocked error")})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = c.DoBuild(reqGet, &mockResponder{})
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

type mockWebClient struct {
	expectedUrl    string
	expectedMethod string
	resp           *http.Response
	err            error
}

func (m *mockWebClient) Do(req *http.Request) (*http.Response, error) {
	if m.expectedUrl != "" && req.URL.String() != m.expectedUrl {
		return nil, errors.New("unmatching url")
	}
	if m.expectedMethod != "" && req.Method != m.expectedMethod {
		return nil, errors.New("unmatching method")
	}
	return m.resp, m.err
}

type mockResponder struct {
	err error
}

func (m *mockResponder) Respond(*http.Response) error {
	return m.err
}
