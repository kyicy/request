package request

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractTag(t *testing.T) {
	_, err := exractTag(nil)
	require.Error(t, err)

	meta := struct {
		ContentType string `rh:"Content-Type"`
		Token       string `rh:"Token,omitempty"`
		Page        int    `rq:"page"`
		CurrentPage int    `rq:"current_page, omitempty"`
	}{"json", "", 1, 0}

	_, err = exractTag(meta)
	require.NoError(t, err)

	extracted, err := exractTag(&meta)
	require.NoError(t, err)
	header := extracted[headerTagName]
	require.Equal(t, "json", header["Content-Type"])
	require.NotContains(t, header, "Token")

	query := extracted[queryTagName]
	require.Equal(t, "1", query["page"])
	require.NotContains(t, query, "current_page")
}

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
