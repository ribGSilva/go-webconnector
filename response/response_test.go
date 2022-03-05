package response

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewResponder(t *testing.T) {
	_, err := NewResponder()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestNewResponderEmpty(t *testing.T) {
	r, err := NewResponder()
	_ = r.Respond(&http.Response{StatusCode: 200})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestNewResponderFor(t *testing.T) {
	var ok bool
	r, err := NewResponder(For(200, func(response Response) error {
		ok = true
		return nil
	}))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_ = r.Respond(&http.Response{StatusCode: 200})
	if !ok {
		t.Error("error using default handler")
		t.FailNow()
	}
}

func TestNewResponderForDefault(t *testing.T) {
	var ok bool
	r, err := NewResponder(ForDefault(func(response Response) error {
		ok = true
		return nil
	}))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_ = r.Respond(&http.Response{StatusCode: 200})
	if !ok {
		t.Error("error using default handler")
		t.FailNow()
	}
}

func TestNewResponderForStatus(t *testing.T) {
	var ok bool
	r, err := NewResponder(ForDefault(func(response Response) error {
		ok = true
		return nil
	}), ForStatus(200))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_ = r.Respond(&http.Response{StatusCode: 200})
	if ok {
		t.Error("error using status handler")
		t.FailNow()
	}
}

func TestNewResponderForString(t *testing.T) {
	var resp string
	r, err := NewResponder(ForString(200, &resp))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_ = r.Respond(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString("name field"))})
	if resp != "name field" {
		t.Error("error using string responder")
		t.FailNow()
	}
}

func TestNewResponderForJson(t *testing.T) {
	resp := struct {
		Name string `json:"name"`
	}{Name: ""}
	r, err := NewResponder(ForJson(200, &resp))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	marshal, _ := json.Marshal(struct {
		Name string `json:"name"`
	}{Name: "name field"})
	_ = r.Respond(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader(marshal))})
	if resp.Name != "name field" {
		t.Error("error using json responder")
		t.FailNow()
	}
}

func TestNewResponderForXml(t *testing.T) {
	resp := struct {
		XMLName xml.Name `xml:"obj"`
		Name    string   `xml:"name"`
	}{Name: ""}
	r, err := NewResponder(ForXml(200, &resp))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	marshal, _ := xml.Marshal(struct {
		XMLName xml.Name `xml:"obj"`
		Name    string   `xml:"name"`
	}{Name: "name field"})
	_ = r.Respond(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader(marshal))})
	if resp.Name != "name field" {
		t.Error("error using xml responder")
		t.FailNow()
	}
}

func TestNewResponderNilBody(t *testing.T) {
	var ok bool
	r, err := NewResponder(ForDefault(func(response Response) error {
		ok = true
		return nil
	}))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_ = r.Respond(nil)
	if ok {
		t.Error("error using nil response")
		t.FailNow()
	}
}

func TestNewResponderOptionErr(t *testing.T) {
	_, err := NewResponder(func(responder *Responder) error {
		return errors.New("mocked error")
	})
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

func TestNewResponderForStringError(t *testing.T) {
	resp := ""
	r, err := NewResponder(ForString(200, &resp))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	errReq := r.Respond(&http.Response{StatusCode: 200, Body: mockedErrorReadCloser{}})
	if errReq == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

func TestNewResponderForJsonError(t *testing.T) {
	resp := struct {
		Name string `json:"name"`
	}{Name: ""}
	r, err := NewResponder(ForJson(200, &resp))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	errReq := r.Respond(&http.Response{StatusCode: 200, Body: mockedErrorReadCloser{}})
	if errReq == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

func TestNewResponderForXmlError(t *testing.T) {
	resp := struct {
		Name string `xml:"name"`
	}{Name: ""}
	r, err := NewResponder(ForXml(200, &resp))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	errReq := r.Respond(&http.Response{StatusCode: 200, Body: mockedErrorReadCloser{}})
	if errReq == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

type mockedErrorReadCloser struct {
}

func (m mockedErrorReadCloser) Read(_ []byte) (n int, err error) {
	return 0, errors.New("expected error")
}

func (m mockedErrorReadCloser) Close() error {
	return errors.New("expected error")
}
