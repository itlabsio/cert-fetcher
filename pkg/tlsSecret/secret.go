package tlsSecret

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

  "cert-fetcher/pkg/tlsSecret/vaultApi"
)

type TlsSecret struct {
	Cert      string `json:"cert,omitempty"`
	FullChain string `json:"fullchain,omitempty"`
	Pkey      string `json:"key,omitempty"`
}

func newTlsSecret(s *vaultApi.VaultKVV2Data) (*TlsSecret, error) {
	tlss := new(TlsSecret)
	err := json.Unmarshal(s.Data, &tlss)
	return tlss, err
}

type getter interface {
	GetWithContext(ctx context.Context, path string) (io.ReadCloser, error)
}

func ReadTLSSecret(ctx context.Context, path string, v getter) (secret *TlsSecret, err error) {
	r, err := v.GetWithContext(ctx, path)
	defer func() {
		if err == nil {
			cErr := r.Close()
			err = cErr
		}
	}()
	if err != nil {
		return nil, err
	}
	dataDecoder := json.NewDecoder(r)
	data := new(vaultApi.RawResponse)
	if err = dataDecoder.Decode(data); err != nil {
		return nil, err
	}
	if data.Data != nil && data.Data.Data != nil {
		return newTlsSecret(data.Data)
	}
	return nil, fmt.Errorf("not found data in request")
}
