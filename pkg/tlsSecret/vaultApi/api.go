package vaultApi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

type VaultAuthMethodNum int

const vaultRequestTimeoutSeconds = 5

type VaultService struct {
	client *api.Client
}

func NewWithKubernetes(ctx context.Context, address, authPath, role, jwt string, timeout int32) (*VaultService, error) {
	v := new(VaultService)
	config := api.DefaultConfig()

	tlsConf := api.TLSConfig{Insecure: true}
	if err := config.ConfigureTLS(&tlsConf); err != nil {
		return nil, err
	}
	if address != "" {
		config.Address = address
	}
	if timeout == 0 {
		config.Timeout = time.Second * time.Duration(vaultRequestTimeoutSeconds)
	} else {
		config.Timeout = time.Second * time.Duration(timeout)
	}
	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/auth/%s/login", authPath)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"role": role,
		"jwt":  jwt,
	})
	if err != nil {
		return nil, err
	}
	if secret.Auth.ClientToken == "" {
		return nil, fmt.Errorf("Not found token")
	}
	c.SetToken(secret.Auth.ClientToken)
	v.client = c
	return v, nil
}

func NewWithToken(address, token string, timeout int32) (*VaultService, error) {
	v := new(VaultService)
	config := api.DefaultConfig()
	//TODO make configure
	tlsConf := api.TLSConfig{Insecure: true}
	if err := config.ConfigureTLS(&tlsConf); err != nil {
		return nil, err
	}
	if address != "" {
		config.Address = address
	}
	if timeout == 0 {
		config.Timeout = time.Second * time.Duration(vaultRequestTimeoutSeconds)
	} else {
		config.Timeout = time.Second * time.Duration(timeout)
	}
	c, err := api.NewClient(config)
	if token != "" {
		c.SetToken(token)
	}
	if err != nil {
		return nil, err
	}
	v.client = c
	return v, nil
}

func (v *VaultService) Client() *api.Client {
	return v.client
}

func (v *VaultService) GetWithContext(ctx context.Context, path string) (io.ReadCloser, error) {
	params := requestParams{}
	request := v.requestPath("GET", path, params)
	response, err := v.client.RawRequestWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, fmt.Errorf("No data")
	}
	return response.Body, nil
}

func (v *VaultService) PostWithContext(ctx context.Context, path string, params requestParams) (io.ReadCloser, error) {
	request := v.requestPath("POST", path, params)
	response, err := v.client.RawRequestWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

type requestParams struct {
	params map[string][]string
	body   []byte
}

func (v *VaultService) requestPath(method, path string, params requestParams) *api.Request {
	path = strings.TrimLeft(path, "/")
	fullPath := fmt.Sprintf("/v1/%s", path)
	request := v.client.NewRequest(method, fullPath)
	switch method {
	case "GET":
		request.Params = url.Values(params.params)
	case "POST":
		request.BodyBytes = params.body
	default:

	}
	return request
}

type RawResponse struct {
	RequestId     string         `json:"request_id,omitempty"`
	LeaseId       string         `json:"lease_id,omitempty"`
	LeaseDuration int32          `json:"lease_duration,omitempty"`
	Renewable     bool           `json:"renewable,omitempty"`
	Warnings      []string       `json:"warnings"`
	Data          *VaultKVV2Data `json:"data,omitempty"`
}
type VaultKVV2Data struct {
	Metadata *ValueMeta      `json:"metadata"`
	Data     json.RawMessage `json:"data"`
}

type ValueMeta struct {
	CreatedTime  string `json:"created_time,omitempty"`
	DeletionTime string `json:"deletion_time,omitempty"`
	Destroyed    bool   `json:"destroyed,omitempty"`
	Version      uint16 `json:"version,omitempty"`
}
