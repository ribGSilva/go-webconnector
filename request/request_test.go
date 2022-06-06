package request

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const host = "http://defaultHost"

func TestNew(t *testing.T) {
	r, err := New(host)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if r.URL.String() != host {
		t.Errorf("final url does not match: expected %s, result: %s", host, r.URL.String())
		t.FailNow()
	}
	expectedLen := 0
	if len(r.Header) > expectedLen {
		t.Errorf("final headers len does not match: expected %d, result: %d", expectedLen, len(r.Header))
		t.FailNow()
	}
	expectedMethod := http.MethodGet
	if r.Method != expectedMethod {
		t.Errorf("final method does not match: expected %s, result: %s", expectedMethod, r.Method)
		t.FailNow()
	}
}

func TestNewMethod(t *testing.T) {
	r, err := New(host, Method(http.MethodPost))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if r.Method != http.MethodPost {
		t.Errorf("final method does not match: expected %s, result: %s", http.MethodPost, r.Method)
		t.FailNow()
	}
}

func TestNewPath(t *testing.T) {
	path := "/newpath"
	r, err := New(host, Path(path))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expectedUrl := host + path
	if r.URL.String() != expectedUrl {
		t.Errorf("final url does not match: expected %s, result: %s", expectedUrl, r.URL.String())
		t.FailNow()
	}
}

func TestNewCtx(t *testing.T) {
	ctx := context.Background()
	r, err := New(host, Context(ctx))
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
	r, err := New(host, Headers(http.Header{
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
		Header(header, headerV),
		Header(header, headerV2),
		Header(headerInt, headerIntV),
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
		Queries(map[string][]interface{}{
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
		Query(query, queryV),
		Query(query, queryV2),
		Query(queryInt, queryIntV),
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

func TestNewParam(t *testing.T) {
	param := "user"
	paramV := "userValue"
	paramInt := "myQueryInt"
	paramIntV := "myQueryIntValue"
	r, err := New(host,
		Path("/:"+param+"/:"+param+"/:"+paramInt),
		Param(param, paramV),
		Param(paramInt, paramIntV),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expected := "/" + paramV + "/" + paramV + "/" + paramIntV
	if !strings.Contains(r.URL.String(), expected) {
		t.Errorf("final url does not has params: expected %s, result: %s", expected, r.URL.String())
		t.FailNow()
	}
}

func TestNewParams(t *testing.T) {
	param := "user"
	paramV := "userValue"
	paramInt := "myQueryInt"
	paramIntV := "myQueryIntValue"
	r, err := New(host,
		Path("/:"+param+"/:"+param+"/:"+paramInt),
		Params(map[string]interface{}{
			param:    paramV,
			paramInt: paramIntV,
		}),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expected := "/" + paramV + "/" + paramV + "/" + paramIntV
	if !strings.Contains(r.URL.String(), expected) {
		t.Errorf("final url does not has params: expected %s, result: %s", expected, r.URL.String())
		t.FailNow()
	}
}

func TestNewBody(t *testing.T) {
	body := "myBody"
	r, err := New(host,
		Body(body),
		Encoder(func(any) ([]byte, error) {
			return []byte(body), nil
		}),
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
		String(body),
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
		JSON(body),
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

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("final header does not match: expected %s, result: %s", "application/json", r.Header.Get("Content-Type"))
		t.FailNow()
	}
}

func TestNewXml(t *testing.T) {
	body := struct {
		XMLName xml.Name `xml:"obj"`
		Field   string   `xml:"field"`
	}{Field: "myField"}

	r, err := New(host,
		XML(body),
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

	if r.Header.Get("Content-Type") != "application/xml" {
		t.Errorf("final header does not match: expected %s, result: %s", "application/xml", r.Header.Get("Content-Type"))
		t.FailNow()
	}
}

func TestNewJsonError(t *testing.T) {
	_, err := New(host,
		JSON(make(chan int, 1)),
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
		XML(body),
	)

	if err == nil {
		t.Error("it supposed to return an error")
		t.FailNow()
	}
}

func TestNewRequestError(t *testing.T) {
	_, err := New("", Context(nil))

	if err == nil {
		t.Error("it supposed to return an error")
		t.FailNow()
	}
}
