package config

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

const (
	confJson = `{"secrets":[{
		"vaultPath":"secret/data/acme/tls_certificates/*.acme.org",
		"secretName": "wildcard-acme-org-tls",
		"namespaces":["default", "johndoe", "janedoe"]}]
	}`
	noConfJson = `{"secrets":[{
		"vualtPath":"secret/data/acme/tls_certificates/*.acme.org",
		"secretNmae": "wildcard-acme-org-tls",
		"namespaces11":["default", "johndoe", "janedoe"]}]
	}`
	badJson  = `BadFormat`
	confYaml = `
secrets:
- vaultPath: secret/data/acme/tls_certificates/*.acme.org
  secretName: "wildcard-acme-org-tls"
  namespaces:
  - "default"
  - "johndoe"
  - "janedoe"`
	noConfYaml = `---
secrets:
- vualtPath: "secret/data/acme/tls_certificates/*.acme.org"
  secretNmae: "wildcard-acme-org-tls"
  namespaces101:
  - "default"
  - "johndoe-dev"
  - "janedoe-dev"`
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
				if conf.Secrets[0].SecretName != "wildcard-acme-org-tls" {
					t.Errorf("Secret name: expected '%s', decoded '%s' from\n'%s'",
						"wildcard-acme-org-tls",
						conf.Secrets[0].SecretName,
						tst)
				}
				t.Logf("%s", conf.Secrets[0].Namespaces)

				if conf.Secrets[0].VaultPath != "secret/data/acme/tls_certificates/*.acme.org" {
					t.Errorf("Vault path: expected '%s', decoded '%s'",
						"secret/data/acme/tls_certificates/*.acme.org",
						conf.Secrets[0].VaultPath)
				}
			} else {
				t.Error(err.Error())
			}
		//noConfJson
		case 1:
			if err == nil {
				if conf.Secrets[0].SecretName == "wildcard-acme-org-tls" &&
					conf.Secrets[0].VaultPath == "secret/data/acme/tls_certificates/*.acme.org" {
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
				if conf.Secrets[0].SecretName != "wildcard-acme-org-tls" {
					t.Errorf("Secret name: expected '%s', decoded '%s' from\n'%s'",
						"wildcard-acme-org-tls",
						conf.Secrets[0].SecretName,
						tst)
				}
				if conf.Secrets[0].VaultPath != "secret/data/acme/tls_certificates/*.acme.org" {
					t.Errorf("Vault path: expected '%s', decoded '%s'",
						"secret/data/path/to/wildcard_certificates/*.acme.org",
						conf.Secrets[0].VaultPath)
				}
			} else {
				t.Errorf("%s\n%s", err.Error(), tst)
			}
		//noConfJson
		case 1:
			if err == nil {
				if conf.Secrets[0].SecretName == "wildcard-tls" &&
					conf.Secrets[0].VaultPath == "ssecret/data/acme/tls_certificates/*.acme.org" {
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
