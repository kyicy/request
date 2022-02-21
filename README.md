[![Go Report Card][goreportcard-img]][goreportcard]
[![License][license-img]][license]
[![GoDoc][doc-img]][doc]
[![Build Status][ci-img]][ci]
[![Coverage Status][cov-img]][cov]



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

[goreportcard-img]: https://goreportcard.com/badge/kyicy/request
[goreportcard]: https://goreportcard.com/report/kyicy/request
[license-img]: https://img.shields.io/badge/License-AGPL_v3-blue.svg
[license]: https://github.com/kyicy/request/blob/master/LICENSE
[doc-img]: https://pkg.go.dev/badge/github.com/kyicy/request
[doc]: https://pkg.go.dev/github.com/kyicy/request?tab=doc
[ci-img]: https://github.com/kyicy/request/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/kyicy/request/actions/workflows/go.yml
[cov-img]: https://codecov.io/gh/kyicy/request/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/kyicy/request