# GO WebConnector

This lib is an easy and re-utilizable way to make http requests

## Request

To make requests, follow the example:

```go
func buildReq(ctx context.Context, id string, body interface{}) {
    req, err := New("my.host.com",
        WithContext(ctx),
        WithMethod(MethodPatch), // by default is GET
        WithProtocol("https"), // by default is http
        WithPath("/path/:id"),
        WithParam("id", id),
        WithQuery("myQuery", "someValue"),
        WithHeader("Authorization", "myauth"),
        WithJson(body),
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

	responder, err := NewResponder(
		ForStatus(404), // Does nothing
		ForJson(200, &response),
		ForDefault(func(response Response) error {
			return errors.New("response: not mapped status")
		}),
	)
	if err != nil {
		return err
	}

	return responder.Respond(resp)
}
```

## Connector

Centralize the requests in one connector:

```go
func execRequests() error {
	getAllPath := "/"
	getPath := "/:id"
	postPath := "/"

	c, err := New("my.host.com",
		http.DefaultClient,
		WithGeneral(request.WithHeader("My-Header", "myHeaderKey")), // apply in all requests
		WithPath(getPath), // get
		WithPath(postPath, request.WithMethod(request.MethodPost)), // post
	)
	if err != nil {
		return err
	}

	//get
	getBody := struct {
		Name string `json:"name"`
	}{}
	responder, err := response.NewResponder(response.ForJson(200, getBody))
	if err != nil {
		return err
	}
	err = c.DoBuild(getPath, &responder, request.WithParam("id", "123"))
	if err != nil {
		return err
	}
	fmt.Printf("%+v \n", getBody)

	//post
	postBody := struct {
		Name string `json:"name"`
	}{Name: "my name"}
	responder, err = response.NewResponder(
		response.ForStatus(201),
		response.ForDefault(func(r response.Response) error {
			return errors.New("error creating")
		}),
	)
	if err != nil {
		return err
	}
	err = c.DoBuild(getPath, &responder, request.WithJson(postBody))
	if err != nil {
		return err
	}
	fmt.Println("created")

	//get all
	getAll := make([]struct {
		Name string `json:"name"`
	}, 0)

	responder, err = response.NewResponder(response.ForJson(200, getAll))
	if err != nil {
		return err
	}
	err = c.DoBuild(getAllPath, &responder, request.WithQuery("page", "3"))
	if err != nil {
		return err
	}
	fmt.Printf("%+v \n", getAll)

	return nil
}
```

Developer:
Gabriel Ribeiro Silva
https://www.linkedin.com/in/gabriel-ribeiro-silva/