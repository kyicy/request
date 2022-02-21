[![Go Report Card](https://goreportcard.com/badge/kyicy/request)](https://goreportcard.com/report/kyicy/request)
[![License](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://github.com/kyicy/request/blob/master/LICENSE)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/kyicy/request?tab=doc)



# request
A wrapper for net/http client


## Example Customize url query with struct tag `rq`
```go
func TestRequestQuery(t *testing.T) {
	client := &http.Client{}
	ctx := context.Background()
	rp, err := NewRequestProvider(client)
	require.NoError(t, err)
	r, err := NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://jsonplaceholder.typicode.com/posts",
		struct {
			UserId string `rq:"userId"`
		}{
			UserId: "1",
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, "https://jsonplaceholder.typicode.com/posts?userId=1", r.URL.String())
	res, err := rp.Do(r)
	require.NoError(t, err)
	res.Body.Close()

}
```

## Example Customize header with struct tag `rh`
```go
func TestRequestHeader(t *testing.T) {
	client := &http.Client{}
	ctx := context.Background()
	rp, err := NewRequestProvider(client)
	require.NoError(t, err)

	bs, err := json.Marshal(&struct {
		Title  string `json:"title"`
		Body   string `json:"bar"`
		UserId int    `json:"userId"`
	}{"foo", "bar", 1})
	require.NoError(t, err)

	r, err := NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://jsonplaceholder.typicode.com/posts",
		struct {
			ContentType string `rh:"Content-Type"`
		}{
			ContentType: "application/json; charset=UTF-8",
		},
		bytes.NewReader(bs),
	)
	require.NoError(t, err)
	require.Equal(t, "application/json; charset=UTF-8", r.Header.Get("Content-Type"))
	res, err := rp.Do(r)
	require.NoError(t, err)
	res.Body.Close()
}
```