package tlsSecret

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

var testResponse = `{
	"request_id": "7d2858cf-250c-16c3-e567-4682bf0ec0aa",
	"lease_id": "",
	"lease_duration": 0,
	"renewable": false,
	"data": {
	  "data": {
		"cert": "cert-data",
		"fullchain": "cert1\ncert2\ncert3",
		"key": "key-data"
	  },
	  "metadata": {
		"created_time": "2020-10-06T11:15:40.824000028Z",
		"deletion_time": "",
		"destroyed": false,
		"version": 1
	  }
	},
	"warnings": null
  }`

type testGetter func(ctx context.Context, path string) (io.ReadCloser, error)

func (t testGetter) GetWithContext(ctx context.Context, path string) (io.ReadCloser, error) {
	return t(ctx, path)
}

func TestReadTlsSecret(t *testing.T) {
	ctx := context.Background()
	s, err := ReadTLSSecret(ctx, "/secret", testGetter(func(ctx context.Context, path string) (io.ReadCloser, error) {
		return ioutil.NopCloser(strings.NewReader(testResponse)), nil
	}))
	if err != nil {
		t.Error(err.Error())
	}
	if s.Cert == "" {
		t.Errorf("No value")
	}
	fmt.Println(s.Cert)
}
