package main

import (
	"context"
	"os"
	"strings"

	"cert-fetcher/pkg/kubeSecreter"
	"cert-fetcher/pkg/tlsSecret"
	"cert-fetcher/pkg/tlsSecret/vaultApi"
	"github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
)

const (
	vaultTimeoutSeconds = 10
)

func main() {
	pconfig := configure()
	ctx := context.Background()

	vsvc, err := vaultApi.NewWithKubernetes(ctx, pconfig.vault.addr, pconfig.vault.authpath, pconfig.vault.role, pconfig.vault.jwt, vaultTimeoutSeconds)
	if err != nil {
		logrus.Error(err.Error())
		vsvc, err = vaultApi.NewWithToken(pconfig.vault.addr, pconfig.vault.token, vaultTimeoutSeconds)
		if err != nil {
			logrus.Error(err.Error())
			os.Exit(-1)
		}
	}
	logrus.Printf("Vault configured")
	data := getNamespace(pconfig.kubernetes)
	if data != nil {
		for _, s := range data.Items {
			for key, value := range s.GetAnnotations() {
				if strings.Contains(key, "itlabs.io/certfetcher.") {
					keyArray := strings.Split(key, ".")
					secretName := keyArray[len(keyArray)-1]
					if data := getTLSData(ctx, value, vsvc); data != nil {
						updateWithNamespaceKubeSecrets(ctx, pconfig.kubernetes, secretName, *data, s.Name)
					}
				}
			}
		}
	}
}

func getNamespace(kubectl *kubeSecreter.KubeClient) *v1.NamespaceList {
	data, err := kubeSecreter.ReadNamespaces(kubectl.GetClientSet())
	if err != nil {
		logrus.Errorf("Ошибка получения namespaces: %s", err.Error())
		return nil
	}
	return data
}

func getTLSData(ctx context.Context, vaultPath string, vsvc *vaultApi.VaultService) *tlsSecret.TlsSecret {
	data, err := tlsSecret.ReadTLSSecret(ctx, vaultPath, vsvc)
	if err != nil {
		logrus.Errorf("Ошибка чтения секрета: %s", err.Error())
		return nil
	}
	return data
}

func updateWithNamespaceKubeSecrets(ctx context.Context, kubectl *kubeSecreter.KubeClient, secretName string, secret tlsSecret.TlsSecret, ns string) {
	logrus.Infof("Обновляются данные в кластере")
	nsl, err := kubeSecreter.ReadSecrets(kubectl.GetClientSet(), secretName, ns)
	if err != nil {
		logrus.Errorf("Can not get secrets for namespace %s: %s", ns, err.Error())
		return
	}
	if len(nsl.Items) == 0 {
		logrus.Errorf("Not found secret %s in namespace %s", secretName, ns)
		err = kubeSecreter.CreateNewSecret(kubectl.GetClientSet(), secretName, ns)
		if err != nil {
			logrus.Errorf("Can not create secret %s in namespace %s: %s", secretName, ns, err)
		} else {
			logrus.Infof("Not found secret %s in namespace %s", secretName, ns)
		}
	}

	sl, err := kubeSecreter.ReadSecrets(kubectl.GetClientSet(), secretName, ns)
	if err != nil {
		logrus.Errorf("Can not get list of objects: %s", err.Error())
		return
	}

	logrus.Infof("Get %d objets", len(sl.Items))
	for _, s := range sl.Items {
		logrus.Infof("Update object %s/%s", ns, s.Name)
		if err := kubeSecreter.UpdateSecret(&s, []byte(secret.FullChain), []byte(secret.Pkey)); err != nil {
			logrus.Errorf("Object data %s/%s does not update: %s", ns, s.Name, err.Error())
			continue
		}
		_, err := kubectl.GetClientSet().CoreV1().Secrets(ns).Update(&s)
		if err != nil {
			logrus.Errorf("Can not save object %s/%s: %s", ns, s.Name, err.Error())
		} else {
			logrus.Infof("Object saved: %s/%s", s.GetNamespace(), s.GetName())
		}
	}
}
