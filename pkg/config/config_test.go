package config

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

const (
	confJson = `{"secrets":[{
		"vaultPath":"secret/data/itlabs/tls_certificates/*.itlabs.io",
		"secretName": "wildcard-itlabs-io-tls",
		"namespaces":["default", "ennergiia-dev", "sdvor-dev"]}]
	}`
	noConfJson = `{"secrets":[{
		"vualtPath":"secret/data/itlabs/tls_certificates/*.itlabs.io",
		"secretNmae": "wildcard-itlabs-io-tls",
		"namespaces11":["default", "ennergiia-dev", "sdvor-dev"]}]
	}`
	badJson  = `BadFormat`
	confYaml = `
secrets:
- vaultPath: secret/data/itlabs/tls_certificates/*.itlabs.io
  secretName: "wildcard-itlabs-io-tls"
  namespaces:
  - "default"
  - "ennergiia-dev"
  - "sdvor-dev"`
	noConfYaml = `---
secrets:
- vualtPath: "secret/data/itlabs/tls_certificates/*.itlabs.io"
  secretNmae: "wildcard-itlabs-io-tls"
  namespaces101:
  - "default"
  - "ennergiia-dev"
  - "sdvor-dev"`
	badYaml = `Bad yaml`
)

func TestConfigJsonDecoder(t *testing.T) {
	var test = []string{confJson, noConfJson, badJson}
	for i, tst := range test {
		conf, err := ReadConfig(json.NewDecoder(strings.NewReader(tst)))
		switch i {
		//validConfig
		case 0:
			if err == nil {
				if conf.Secrets[0].SecretName != "wildcard-itlabs-io-tls" {
					t.Errorf("Secret name: expected '%s', decoded '%s' from\n'%s'",
						"wildcard-itlabs-io-tls",
						conf.Secrets[0].SecretName,
						tst)
				}
				t.Logf("%s", conf.Secrets[0].Namespaces)

				if conf.Secrets[0].VaultPath != "secret/data/itlabs/tls_certificates/*.itlabs.io" {
					t.Errorf("Vault path: expected '%s', decoded '%s'",
						"secret/data/itlabs/tls_certificates/*.itlabs.io",
						conf.Secrets[0].VaultPath)
				}
			} else {
				t.Error(err.Error())
			}
		//noConfJson
		case 1:
			if err == nil {
				if conf.Secrets[0].SecretName == "wildcard-itlabs-io-tls" &&
					conf.Secrets[0].VaultPath == "secret/data/itlabs/tls_certificates/*.itlabs.io" {
					t.Error("Secret name")
				}
			} else {
				t.Error(err.Error())
			}
		//badJson
		case 2:
			if err == nil {
				t.Fail()
			}
		}
	}
}
func TestConfigDecoder(t *testing.T) {
	var tests = []string{confYaml, noConfYaml, badYaml}
	for i, tst := range tests {
		conf, err := ReadConfig(yaml.NewDecoder(strings.NewReader(tst)))
		switch i {
		//validConfig
		case 0:
			if err == nil {
				if conf.Secrets[0].SecretName != "wildcard-itlabs-io-tls" {
					t.Errorf("Secret name: expected '%s', decoded '%s' from\n'%s'",
						"wildcard-itlabs-io-tls",
						conf.Secrets[0].SecretName,
						tst)
				}
				if conf.Secrets[0].VaultPath != "secret/data/itlabs/tls_certificates/*.itlabs.io" {
					t.Errorf("Vault path: expected '%s', decoded '%s'",
						"secret/data/itlabs/tls_certificates/*.itlabs.io",
						conf.Secrets[0].VaultPath)
				}
			} else {
				t.Errorf("%s\n%s", err.Error(), tst)
			}
		//noConfJson
		case 1:
			if err == nil {
				if conf.Secrets[0].SecretName == "wildcard-itlabs-io-tls" &&
					conf.Secrets[0].VaultPath == "secret/data/itlabs/tls_certificates/*.itlabs.io" {
					t.Error("Secret name")
				}
			} else {
				t.Errorf("%s\n%s", err.Error(), tst)
			}
		//badJson
		case 2:
			if err == nil {
				t.Fail()
			}
		}
	}
}
