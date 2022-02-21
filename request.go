package request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	queryTagName  = "rq" // request query
	headerTagName = "rh" // request header
)

func exractTag(
	input interface{},
) (extracted map[string]map[string]string, err error) {
	extracted = map[string]map[string]string{
		queryTagName:  make(map[string]string),
		headerTagName: make(map[string]string),
	}

	st := reflect.TypeOf(input)
	sv := reflect.ValueOf(input)

	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
		st = st.Elem()
	}

	if sv.Kind() != reflect.Struct {
		err = fmt.Errorf("input should be a struct or a struct pointer")
		return
	}

	for i := 0; i < st.NumField(); i++ {
		v := sv.Field(i)
		t := st.Field(i)
		for _, tagField := range []string{queryTagName, headerTagName} {
			attr, ok := t.Tag.Lookup(tagField)
			if !ok {
				continue
			}
			ks := strings.Split(attr, ",")
			if len(ks) == 0 {
				continue
			}
			omitempty := false
			for i, k := range ks {
				k = strings.TrimSpace(k)
				ks[i] = k
				if i > 0 && k == "omitempty" {
					omitempty = true
				}
			}

			if v.IsZero() && omitempty {
				continue
			}

			v := fmt.Sprintf("%v", v)
			extracted[tagField][ks[0]] = v
		}
	}
	return
}

// RequestProvider embeds a native http.Client pointer
type RequestProvider struct {
	*http.Client
}

// NewRequestProvider : RequestProvider constructor
func NewRequestProvider(client *http.Client) (*RequestProvider, error) {
	if client == nil {
		return nil, fmt.Errorf("http client is nil")
	}
	return &RequestProvider{
		Client: client,
	}, nil
}

// NewRequestWithContext wraps http.NewRequestWithContext with query and header extacted from struct tags
func NewRequestWithContext(
	ctx context.Context,
	method string,
	baseUrl string,
	meta interface{},
	body io.Reader,
) (*http.Request, error) {
	exracted, err := exractTag(meta)
	if err != nil {
		return nil, err
	}
	uri, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	q := make(url.Values)
	for k, v := range exracted[queryTagName] {
		q.Set(k, v)
	}
	uri.RawQuery = q.Encode()
	targetUrl := uri.String()

	req, err := http.NewRequestWithContext(ctx, method, targetUrl, body)
	if err != nil {
		return nil, err
	}
	for k, v := range exracted[headerTagName] {
		req.Header.Set(k, v)
	}
	return req, nil
}
