# GO WebConnector

[![Go Reference](https://pkg.go.dev/badge/github.com/ribGSilva/go-webconnector.svg)](https://pkg.go.dev/github.com/ribGSilva/go-webconnector)

This lib is an easy and re-utilizable way to make http requests

## Install

To install just run:

```ssh
    go get github.com/ribGSilva/go-webconnector
```

## Request

To make requests, follow the example:

```go
func buildReq(ctx context.Context, id string, body interface{}) {
    req := rq.New("my.host.com",
	    rq.Context(ctx),
        rq.Method(MethodPatch), // by default is GET
        rq.Protocol("https"), // by default is http
        rq.Path("/path/:id"),
        rq.Param("id", id),
        rq.Query("myQuery", "someValue"),
        rq.Header("Authorization", "myauth"),
        rq.Json(body),
    )
}
```

## Response

To handle response, follow the example:

```go
func handleResponse(resp *http.Response) error {
	response := struct {
		Name string `json:"name"`
	}{}

	responder := rp.New(
	    rp.Status(http.StatusNotFound), // Does nothing
        rp.For(http.StatusOK, func(body io.ReadCloser) (any, error) {
            var b myStruct
            err := json.NewDecoder(body).Decode(&b)
            if err != nil {
                return nil, err
            }
            return b, nil
        }),
        rp.Default(func(responder io.ReadCloser) (any, error) {
            return nil, errors.New("responder: not mapped status")
        }),
	)

	return responder.Respond(resp)
}
```

Developer:
Gabriel Ribeiro Silva
https://www.linkedin.com/in/gabriel-ribeiro-silva/