package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"cert-fetcher/pkg/config"
	"cert-fetcher/pkg/kubeSecreter"
	"github.com/sirupsen/logrus"
)

const (
	kubeSAPath           = "/var/run/secrets/kubernetes.io/serviceaccount/"
	defaultVaultAddr     = "http://127.0.0.1:8200"
	defaultVaultPath     = "kubernetes"
	defaultVaultRole     = "default"
	defaultKubeConfig    = ".kube/config"
	defaultSecretsConfig = "secrets.yaml"
)

func defaultKubeConfigPath() string {
	return fmt.Sprintf("%s/%s", os.Getenv("HOME"), defaultKubeConfig)
}

type programConfig struct {
	vault      vaultConfig
	kubernetes *kubeSecreter.KubeClient
	kubeConfig string
	configYaml string
	secrets    []config.SecretsConfig
}

type vaultConfig struct {
	addr        string
	token       string
	authpath    string
	jwt         string
	role        string
	tlsNoVerify bool
}

func configure() programConfig {
	// TODO need update
	flags, pc := programConfig{}, programConfig{}
	flag.StringVar(&flags.vault.addr, "vault-addr", defaultVaultAddr, "Vault server address")
	flag.StringVar(&flags.vault.authpath, "vault-path", defaultVaultPath, "Vault authmethod path if not token")
	flag.StringVar(&flags.vault.token, "vault-token", "", "Vault client token")
	flag.StringVar(&flags.vault.jwt, "vault-jwt", "", "Vault jwt token")
	flag.StringVar(&flags.vault.role, "vault-role", defaultVaultRole, "vault role")
	flag.StringVar(&flags.kubeConfig, "kubeconfig", defaultKubeConfigPath(), "path to kubeconfig")
	flag.BoolVar(&flags.vault.tlsNoVerify, "tlsnoverify", false, "don't check certificate authority")
	flag.Parse()
	if jwt, err := readInClusterJWT(); err != nil || flags.vault.jwt != "" {
		pc.vault.jwt = flags.vault.jwt
	} else {
		pc.vault.jwt = string(jwt)
	}
	pc.vault.addr = readEnvButPreferFlagString("VAULT_ADDR", flags.vault.addr, defaultVaultAddr)
	pc.vault.token = readEnvButPreferFlagString("VAULT_TOKEN", flags.vault.token, "")
	pc.vault.authpath = readEnvButPreferFlagString("VAULT_PATH", flags.vault.authpath, defaultVaultPath)
	pc.vault.role = readEnvButPreferFlagString("VAULT_ROLE", flags.vault.role, defaultVaultRole)
	pc.vault.tlsNoVerify = readEnvButPreferFlagBool("VAULT_SKIP_VERIFY", flags.vault.tlsNoVerify, false)
	pc.kubeConfig = flags.kubeConfig
	createKubeClient(&flags, &pc)
	return pc
}

func createKubeClient(flagsConfig, programConfig *programConfig) {
	if client, err := kubeSecreter.AuthByDefaultKubeconfig(); err == nil {
		programConfig.kubernetes = client
		return
	} else {
		logrus.Print(err.Error())
	}
	if client, err := kubeSecreter.AuthInCluster(); err == nil {
		programConfig.kubernetes = client
		return
	} else {
		logrus.Printf(err.Error())
	}
}

func readEnvButPreferFlagString(envName, flagValue, defaultValue string) string {
	val := os.Getenv(envName)
	if flagValue == "" || flagValue == defaultValue {
		if val != "" {
			return val
		}
	}
	return flagValue
}

func readEnvButPreferFlagBool(envName string, flagValue, defaultValue bool) bool {
	val := os.Getenv(envName)
	if flagValue == defaultValue {
		if val != "" {
			return stringToBool(val)
		}
	}
	return flagValue
}

func stringToBool(val string) bool {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "F") || strings.HasPrefix(val, "f") {
		return false
	}
	return true
}

func readInClusterJWT() ([]byte, error) {
	tokenFile := fmt.Sprintf("%s/token", kubeSAPath)
	return ioutil.ReadFile(tokenFile)
}
