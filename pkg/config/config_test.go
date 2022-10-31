package config

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

const (
	confJson = `{"secrets":[{
		"vaultPath":"secret/data/path/to/wildcard_certificate",
		"secretName": "wildcard-tls",
		"namespaces":["default", "johndoe", "janedoe"]}]
	}`
	noConfJson = `{"secrets":[{
		"vualtPath":"secret/data/path/to/wildcard_certificate",
		"secretNmae": "wildcard-tls",
		"namespaces11":["default", "johndoe", "janedoe"]}]
	}`
	badJson  = `BadFormat`
	confYaml = `
secrets:
- vaultPath: secret/data/path/to/wildcard_certificate
  secretName: "wildcard-tls"
  namespaces:
  - "default"
  - "johndoe"
  - "janedoe"`
	noConfYaml = `---
secrets:
- vualtPath: "secret/data/path/to/wildcard_certificate"
  secretNmae: "wildcard-tls"
  namespaces101:
  - "default"
  - "johndoe"
  - "janedoe"`
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
						"wildcard-tls",
						conf.Secrets[0].SecretName,
						tst)
				}
				t.Logf("%s", conf.Secrets[0].Namespaces)

				if conf.Secrets[0].VaultPath != "secret/data/path/to/wildcard_certificate" {
					t.Errorf("Vault path: expected '%s', decoded '%s'",
						"secret/data/path/to/wildcard_certificate",
						conf.Secrets[0].VaultPath)
				}
			} else {
				t.Error(err.Error())
			}
		//noConfJson
		case 1:
			if err == nil {
				if conf.Secrets[0].SecretName == "wildcard-tls" &&
					conf.Secrets[0].VaultPath == "secret/data/path/to/wildcard_certificate" {
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
				if conf.Secrets[0].SecretName != "wildcard-tls" {
					t.Errorf("Secret name: expected '%s', decoded '%s' from\n'%s'",
						"wildcard-tls",
						conf.Secrets[0].SecretName,
						tst)
				}
				if conf.Secrets[0].VaultPath != "secret/data/path/to/wildcard_certificate" {
					t.Errorf("Vault path: expected '%s', decoded '%s'",
						"secret/data/path/to/wildcard_certificate",
						conf.Secrets[0].VaultPath)
				}
			} else {
				t.Errorf("%s\n%s", err.Error(), tst)
			}
		//noConfJson
		case 1:
			if err == nil {
				if conf.Secrets[0].SecretName == "wildcard-tls" &&
					conf.Secrets[0].VaultPath == "ssecret/data/path/to/wildcard_certificate" {
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
