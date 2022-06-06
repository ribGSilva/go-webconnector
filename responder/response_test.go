package responder

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewResponderEmpty(t *testing.T) {
	r := New()
	res := r.Respond(&http.Response{StatusCode: 200})
	if res.Err != ErrNoResponseHandler {
		t.Error(res.Err)
		t.FailNow()
	}
}

func TestNewResponderFor(t *testing.T) {
	r := New(For(200, func(response io.ReadCloser) (any, error) {
		return true, nil
	}))
	res := r.Respond(&http.Response{StatusCode: 200})
	ok := res.Body.(bool)
	if !ok {
		t.Error("error using default handler")
		t.FailNow()
	}
}

func TestNewResponderForDefault(t *testing.T) {
	r := New(Default(func(response io.ReadCloser) (any, error) {
		return true, nil
	}))
	res := r.Respond(&http.Response{StatusCode: 200})
	ok := res.Body.(bool)
	if !ok {
		t.Error("error using default handler")
		t.FailNow()
	}
}

func TestNewResponderForStatus(t *testing.T) {
	r := New(Default(func(response io.ReadCloser) (any, error) {
		return true, nil
	}), Status(200))

	res := r.Respond(&http.Response{StatusCode: 200})
	if res.Err != nil || res.Body != nil {
		t.Error("error using status handler")
		t.FailNow()
	}
}

func TestNewResponderForString(t *testing.T) {
	r := New(String(200))
	res := r.Respond(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString("name field"))})
	if res.Body != "name field" {
		t.Error("error using string responder")
		t.FailNow()
	}
}

type nameStr struct {
	XMLName xml.Name `xml:"obj"`
	Name    string   `json:"name" xml:"name"`
}

func TestNewResponderForJson(t *testing.T) {
	r := New(For(200, func(body io.ReadCloser) (any, error) {
		var b nameStr
		err := json.NewDecoder(body).Decode(&b)
		if err != nil {
			return nil, err
		}
		return b, nil
	}))
	marshal, _ := json.Marshal(nameStr{Name: "name field"})
	res := r.Respond(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader(marshal))})
	resp := res.Body.(nameStr)
	if resp.Name != "name field" {
		t.Error("error using json responder")
		t.FailNow()
	}
}

func TestNewResponderForXml(t *testing.T) {
	r := New(For(200, func(response io.ReadCloser) (any, error) {
		var b nameStr
		err := xml.NewDecoder(response).Decode(&b)
		if err != nil {
			return nil, err
		}
		return b, nil
	}))
	marshal, _ := xml.Marshal(struct {
		XMLName xml.Name `xml:"obj"`
		Name    string   `xml:"name"`
	}{Name: "name field"})
	res := r.Respond(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader(marshal))})
	resp := res.Body.(nameStr)
	if resp.Name != "name field" {
		t.Error("error using xml responder")
		t.FailNow()
	}
}

func TestNewResponderNilBody(t *testing.T) {
	r := New(Default(func(response io.ReadCloser) (any, error) {
		return true, nil
	}))
	res := r.Respond(nil)
	if res.Err != ErrNoHttpResponse {
		t.Error("error using nil responder")
		t.FailNow()
	}
}

func TestNewResponderForStringError(t *testing.T) {
	r := New(String(200))
	res := r.Respond(&http.Response{StatusCode: 200, Body: mockedErrorReadCloser{}})
	if res.Err == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

func TestNewResponderForJsonError(t *testing.T) {
	r := New(String(200))
	res := r.Respond(&http.Response{StatusCode: 200, Body: mockedErrorReadCloser{}})
	if res.Err == nil {
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
