package request

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"strings"
	"testing"
)

const host = "defaultHost"

func TestNew(t *testing.T) {
	r, err := New(host)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expectedUrl := "http://" + host
	if r.URL.String() != expectedUrl {
		t.Errorf("final url does not match: expected %s, result: %s", expectedUrl, r.URL.String())
		t.FailNow()
	}
	expectedLen := 0
	if len(r.Header) > expectedLen {
		t.Errorf("final headers len does not match: expected %d, result: %d", expectedLen, len(r.Header))
		t.FailNow()
	}
	expectedMethod := string(MethodGet)
	if r.Method != expectedMethod {
		t.Errorf("final method does not match: expected %s, result: %s", expectedMethod, r.Method)
		t.FailNow()
	}
}

func TestNewMethod(t *testing.T) {
	r, err := New(host, WithMethod(MethodPost))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expectedMethod := string(MethodPost)
	if r.Method != expectedMethod {
		t.Errorf("final method does not match: expected %s, result: %s", expectedMethod, r.Method)
		t.FailNow()
	}
}

func TestNewPath(t *testing.T) {
	path := "/newpath"
	r, err := New(host, WithPath(path))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expectedUrl := "http://" + host + path
	if r.URL.String() != expectedUrl {
		t.Errorf("final url does not match: expected %s, result: %s", expectedUrl, r.URL.String())
		t.FailNow()
	}
}

func TestNewProtocol(t *testing.T) {
	protocol := "https"
	r, err := New(host, WithProtocol(protocol))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expectedUrl := protocol + "://" + host
	if r.URL.String() != expectedUrl {
		t.Errorf("final url does not match: expected %s, result: %s", expectedUrl, r.URL.String())
		t.FailNow()
	}
}

func TestNewCtx(t *testing.T) {
	ctx := context.Background()
	r, err := New(host, WithContext(ctx))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if r.Context() != ctx {
		t.Errorf("final contexet does not match: expected %s, result: %s", ctx, r.Context())
		t.FailNow()
	}
}

func TestNewHeaders(t *testing.T) {
	header := "Myheader"
	headerV := "myHeaderValue"
	headerV2 := "myHeaderValue2"
	headerInt := "Myheaderint"
	headerIntV := "myHeaderIntValue"
	r, err := New(host, WithHeaders(map[string][]interface{}{
		header:    {headerV, headerV2},
		headerInt: {headerIntV},
	}))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expectedLen := 2
	if len(r.Header) > expectedLen {
		t.Errorf("final headers len does not match: expected %d, result: %d", expectedLen, len(r.Header))
		t.FailNow()
	}

	if r.Header[header][0] != headerV {
		t.Errorf("final header does not match: expected %s, result: %s", headerV, r.Header[header][0])
		t.FailNow()
	}

	if r.Header[header][1] != headerV2 {
		t.Errorf("final header does not match: expected %s, result: %s", headerV2, r.Header[header][1])
		t.FailNow()
	}

	if r.Header[headerInt][0] != headerIntV {
		t.Errorf("final header does not match: expected %s, result: %s", headerIntV, r.Header[headerInt][0])
		t.FailNow()
	}
}

func TestNewHeader(t *testing.T) {
	header := "Myheader"
	headerV := "myHeaderValue"
	headerV2 := "myHeaderValue2"
	headerInt := "Myheaderint"
	headerIntV := "myHeaderIntValue"
	r, err := New(host,
		WithHeader(header, headerV),
		WithHeader(header, headerV2),
		WithHeader(headerInt, headerIntV),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expectedLen := 2
	if len(r.Header) > expectedLen {
		t.Errorf("final headers len does not match: expected %d, result: %d", expectedLen, len(r.Header))
		t.FailNow()
	}

	if r.Header[header][0] != headerV {
		t.Errorf("final header does not match: expected %s, result: %s", headerV, r.Header[header][0])
		t.FailNow()
	}

	if r.Header[header][1] != headerV2 {
		t.Errorf("final header does not match: expected %s, result: %s", headerV2, r.Header[header][1])
		t.FailNow()
	}

	if r.Header[headerInt][0] != headerIntV {
		t.Errorf("final header does not match: expected %s, result: %s", headerIntV, r.Header[headerInt][0])
		t.FailNow()
	}
}

func TestNewQueries(t *testing.T) {
	query := "myQuery"
	queryV := "queryValue"
	queryV2 := "queryValue2"
	queryInt := "myQueryInt"
	queryIntV := "myQueryIntValue"
	r, err := New(host,
		WithQueries(map[string][]interface{}{
			query:    {queryV, queryV2},
			queryInt: {queryIntV},
		}),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	exp1 := "myQuery=queryValue"
	if !strings.Contains(r.URL.String(), exp1) {
		t.Errorf("final url does not has query: expected %s, result: %s", exp1, r.URL.String())
		t.FailNow()
	}
	exp2 := "myQuery=queryValue2"
	if !strings.Contains(r.URL.String(), exp2) {
		t.Errorf("final url does not has query: expected %s, result: %s", exp2, r.URL.String())
		t.FailNow()
	}
	exp3 := "myQueryInt=myQueryIntValue"
	if !strings.Contains(r.URL.String(), exp3) {
		t.Errorf("final url does not has query: expected %s, result: %s", exp3, r.URL.String())
		t.FailNow()
	}
}

func TestNewQuery(t *testing.T) {
	query := "myQuery"
	queryV := "queryValue"
	queryV2 := "queryValue2"
	queryInt := "myQueryInt"
	queryIntV := "myQueryIntValue"
	r, err := New(host,
		WithQuery(query, queryV),
		WithQuery(query, queryV2),
		WithQuery(queryInt, queryIntV),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	exp1 := "myQuery=queryValue"
	if !strings.Contains(r.URL.String(), exp1) {
		t.Errorf("final url does not has query: expected %s, result: %s", exp1, r.URL.String())
		t.FailNow()
	}
	exp2 := "myQuery=queryValue2"
	if !strings.Contains(r.URL.String(), exp2) {
		t.Errorf("final url does not has query: expected %s, result: %s", exp2, r.URL.String())
		t.FailNow()
	}
	exp3 := "myQueryInt=myQueryIntValue"
	if !strings.Contains(r.URL.String(), exp3) {
		t.Errorf("final url does not has query: expected %s, result: %s", exp3, r.URL.String())
		t.FailNow()
	}
}

func TestNewBody(t *testing.T) {
	body := "myBody"
	buffer := bytes.NewBufferString(body)
	r, err := New(host,
		WithBody(buffer),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if body != string(all) {
		t.Errorf("final body does not match: expected %s, result: %s", body, string(all))
		t.FailNow()
	}
}

func TestNewString(t *testing.T) {
	body := "myBody"
	r, err := New(host,
		WithString(body),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if body != string(all) {
		t.Errorf("final body does not match: expected %s, result: %s", body, string(all))
		t.FailNow()
	}
}

func TestNewJson(t *testing.T) {
	body := struct {
		Field string `json:"field"`
	}{Field: "myField"}

	r, err := New(host,
		WithJson(body),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	marshal, _ := json.Marshal(body)

	if string(marshal) != string(all) {
		t.Errorf("final body does not match: expected %s, result: %s", string(marshal), string(all))
		t.FailNow()
	}

	if r.Header[headerContentType][0] != "application/json" {
		t.Errorf("final header does not match: expected %s, result: %s", "application/json", r.Header[headerContentType][0])
		t.FailNow()
	}
}

func TestNewXml(t *testing.T) {
	body := struct {
		XMLName xml.Name `xml:"obj"`
		Field   string   `xml:"field"`
	}{Field: "myField"}

	r, err := New(host,
		WithXml(body),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	marshal, _ := xml.Marshal(body)

	if string(marshal) != string(all) {
		t.Errorf("final body does not match: expected %s, result: %s", string(marshal), string(all))
		t.FailNow()
	}

	if r.Header[headerContentType][0] != "application/xml" {
		t.Errorf("final header does not match: expected %s, result: %s", "application/xml", r.Header[headerContentType][0])
		t.FailNow()
	}
}

func TestNewJsonError(t *testing.T) {
	_, err := New(host,
		WithJson(make(chan int, 1)),
	)

	if err == nil {
		t.Error("it supposed to return an error")
		t.FailNow()
	}
}

func TestNewXmlError(t *testing.T) {
	body := struct {
		Field string `xml:"field"`
	}{Field: "myField"}

	_, err := New(host,
		WithXml(body),
	)

	if err == nil {
		t.Error("it supposed to return an error")
		t.FailNow()
	}
}

func TestNewRequestError(t *testing.T) {
	_, err := New("",
		WithProtocol(""),
	)

	if err == nil {
		t.Error("it supposed to return an error")
		t.FailNow()
	}
}

func TestNewRequestCtxError(t *testing.T) {
	_, err := New("",
		WithContext(context.Background()),
		WithProtocol(""),
	)

	if err == nil {
		t.Error("it supposed to return an error")
		t.FailNow()
	}
}
