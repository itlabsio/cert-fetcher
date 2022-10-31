package kubeSecreter

import (
	"encoding/hex"
	"fmt"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const testKey = "test certificate key"
const testCert = "test certificate"

var kubeClient *KubeClient

func getUpdatedCertAndKey() (cert, key secretBytes) {
	return secretBytes([]byte("updated certificate")), secretBytes([]byte("updated key"))
}

func getTestNsLabels() map[string]string {
	return map[string]string{"deleteme": "true"}
}

func getTestNamespaces() []string {
	return []string{"ns1", "ns2", "ns3", "ns4"}
}

func getTestSecret() secret {
	return secret{}
}

func TestBase64Encoding(t *testing.T) {
	testStr := "Hello!"
	desired := "SGVsbG8h"
	result := encodeBase64([]byte(testStr))
	if string(result) != desired {
		t.Errorf("Кодирование неверно: ожидалось %s, пришло %s", desired, result)
	}
}

func TestBase64Decoding(t *testing.T) {
	desired := "Hello!"
	testStr := "SGVsbG8h"
	result, _ := decodeBase64(base64bytes([]byte(testStr)))
	if string(result) != desired {
		t.Errorf("Кодирование неверно: ожидалось %s, пришло %s", desired, result)
	}

}

func setNsLabels(client *KubeClient, ns string, labels map[string]string) {
	nsObj, _ := client.clientSet.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})

	nsObj.SetLabels(labels)
	client.clientSet.CoreV1().Namespaces().Update(nsObj)
}

func createNS(client *KubeClient) {
	for _, nsn := range getTestNamespaces() {
		ns := v1.Namespace{}
		ns.Name = nsn
		ns.SetLabels(map[string]string{"deleteme": "true"})
		_, err := client.clientSet.CoreV1().Namespaces().Create(&ns)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				setNsLabels(client, nsn, map[string]string{"deleteme": "true"})
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}
func deleteAll(client *KubeClient) {
	for _, nsn := range getTestNamespaces() {
		client.clientSet.CoreV1().Namespaces().Delete(nsn, &metav1.DeleteOptions{})
	}
}

func createSecret(client *KubeClient) {
	s := v1.Secret{}
	for _, nsn := range getTestNamespaces() {
		s.SetName("test-secret")
		s.SetNamespace(nsn)
		s.Type = tlsSecretType
		s.Data = map[string][]byte{
			tlsDataCert: encodeBase64([]byte("testcert")),
			tlsDataKey:  encodeBase64([]byte("testkey")),
		}
		_, err := client.clientSet.CoreV1().Secrets(s.GetNamespace()).Create(&s)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func createSecretWithValidData() *v1.Secret {
	s := new(v1.Secret)
	cert, key := secretBytes([]byte(testCert)), secretBytes([]byte(testKey))
	hashdata := HashCertPlusKey(encodeBase64(cert), encodeBase64(key))
	hashdatastring := hex.EncodeToString(hashdata)
	s.SetAnnotations(map[string]string{
		tlsHashAnnotation: hashdatastring,
	})
	s.Type = tlsSecretType
	s.Data = map[string][]byte{
		tlsDataCert: cert,
		tlsDataKey:  key,
	}
	return s
}

func TestIsHashesEquals(t *testing.T) {
	s := createSecretWithValidData()
	cert, key := encodeBase64(secretBytes([]byte(testCert))), encodeBase64(secretBytes([]byte(testKey)))
	hashdata := HashCertPlusKey(cert, key)
	if !isHashesEquals(s, hashdata) {
		t.Errorf("Ожидалось %s, по факту %s", s.Annotations[tlsHashAnnotation], string(hashdata))
	}
}

func TestAnnotationValidation(t *testing.T) {
	s := createSecretWithValidData()
	//cert, key := encodeBase64(secretBytes([]byte(testCert))), encodeBase64(secretBytes([]byte(testKey)))
	if !isValidHashAnnotation(s) {
		t.Error("Неверная хэш-сумма в аннотации")
	}
}

func TestIsSecretNeedToBeUpdated(t *testing.T) {
	s := createSecretWithValidData()
	cert, key := encodeBase64([]byte(testCert)), encodeBase64([]byte(testKey))
	hashdata := HashCertPlusKey(cert, key)
	if isSecretNeedToBeUpdated(s, hashdata) {
		t.Error("Данные секрета и хэш совпадают, обновление не требуется")
	}
	s1 := s.DeepCopy()
	s1.Data[tlsDataCert] = []byte("Неведомая хуйня")
	if !isSecretNeedToBeUpdated(s1, hashdata) {
		t.Error("Данные секрета изменены, требуется обновление")
	}
	s2 := s.DeepCopy()
	s2.Annotations[tlsHashAnnotation] = "Неведомая хуйня"
	if !isSecretNeedToBeUpdated(s2, hashdata) {
		t.Error("Хэш-сумма в аннотации изменена, требуется обновление")
	}
}

func TestUpdateSecretData(t *testing.T) {
	secret := createSecretWithValidData()
	secret.Type = "NoTLS"
	up, err := updateSecretData(secret, []byte(testCert), []byte(testKey))
	if err == nil {
		t.Error("Неправильный тип, обновление должно завершаться с ошибкой")
	}
	if err != nil && up {
		t.Error(err.Error())
	}
	dk := secret.Data[tlsDataKey]
	if string(dk) != testKey {
		t.Errorf("Ожидалось %s, получено %s", testKey, dk)
	}
	dc := secret.Data[tlsDataCert]
	if string(dc) != string(testCert) {
		t.Errorf("Ожидалось %s, получено %s", testCert, dc)
	}
}

func TestUpdateSecret(t *testing.T) {
	secret := createSecretWithValidData()
	upcrt, upkey := getUpdatedCertAndKey()
	err := UpdateSecret(secret, upcrt, upkey)
	if err != nil {
		t.Errorf("Ошибка обновления: %s", err.Error())
	}
	err = UpdateSecret(secret, upcrt, upkey)
	if err == nil {
		t.Errorf("Попытка повторного обновления должна завершаться с ошибкой")
	}
}

func TestSecretUpdate(t *testing.T) {
	client, err := AuthByDefaultKubeconfig()
	if err != nil {
		t.Skip()
	}
	if client.kubeContext != "minikube" {
		t.Skip()
	}
	createNS(client)
	createSecret(client)
	slst, err := ReadSecrets(client.clientSet, "test-secret", "")
	if err != nil {
		t.Error(err.Error())
	}
	for _, s := range slst.Items {
		updateSecretAnnotationHash(&s, []byte(testCert), []byte(testKey))
		up, err := updateSecretData(&s, []byte(testCert), []byte(testKey))
		if err != nil {
			t.Error(err.Error())
		}
		if !up {
			t.Error("Данные секрета не обновлены")
		}
		s1, err := client.clientSet.CoreV1().Secrets(s.Namespace).Update(&s)
		if err != nil {
			t.Error(err.Error())
		}
		c1, _ := s1.Data[tlsDataCert]
		if string(c1) != testCert {
			t.Errorf("Неверные данные секрета. Ожидалось %s, по факту %s", testCert, string(c1))
		}
	}
	t.Cleanup(func() {
		deleteAll(client)
	})
}
