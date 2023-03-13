package kubeSecreter

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
  "context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type base64bytes []byte
type secretBytes []byte
type hashString string
type hashBytes []byte

const tlsHashAnnotation = "itlabs.io.tls.hash.sha1"
const tlsSecretType = "kubernetes.io/tls"
const tlsDataKey = "tls.key"
const tlsDataCert = "tls.crt"

func HashCertPlusKey(cert, key base64bytes) hashBytes {
	hasher := sha1.New()
	hasher.Write(cert)
	hasher.Write(key)
	h := hasher.Sum(nil)
	return h[:]
}

func ReadNamespaces(k *kubernetes.Clientset) (*v1.NamespaceList, error) {
	labelOptions := metav1.ListOptions{
		LabelSelector: "itlabs.io/certfetcher=true",
	}
	return k.CoreV1().Namespaces().List(context.TODO(), labelOptions)
}

func ReadSecrets(k *kubernetes.Clientset, secretName, namespace string) (*v1.SecretList, error) {
	return k.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.name=" + secretName})
}

func CreateNewSecret(k *kubernetes.Clientset, secretName, namespace string) error {
	//TODO need test
	s := v1.Secret{}
	s.SetName(secretName)
	s.SetNamespace(namespace)
	s.Type = tlsSecretType
	s.Data = map[string][]byte{
		tlsDataCert: encodeBase64([]byte("testcert")),
		tlsDataKey:  encodeBase64([]byte("testkey")),
	}
	labelOptions := metav1.CreateOptions{DryRun: []string{"Hello", "Theodore"}}
	_, err := k.CoreV1().Secrets(namespace).Create(context.TODO(), &s, labelOptions)
	return err
}

func UpdateSecret(secret *v1.Secret, cert, key secretBytes) error {
	hashdata := HashCertPlusKey(encodeBase64(cert), encodeBase64(key))
	if isSecretNeedToBeUpdated(secret, hashdata) {
		updateSecretAnnotationHash(secret, cert, key)
		_, err := updateSecretData(secret, cert, key)
		return err
	}
	return fmt.Errorf("no need any updates")
}

func updateSecretData(secret *v1.Secret, cert, key secretBytes) (bool, error) {
	if secret.Type != tlsSecretType {
		return false, fmt.Errorf("wrong secret type, need type TLS")
	}
	data := make(map[string][]byte)
	data[tlsDataCert] = cert
	data[tlsDataKey] = key
	secret.Data = data
	return true, nil
}

func updateSecretAnnotationHash(secret *v1.Secret, cert, key secretBytes) bool {
	cb64 := encodeBase64(cert)
	kb64 := encodeBase64(key)
	annv := HashCertPlusKey(cb64, kb64)
	changed := updateAnnotation(secret, hashString(hex.EncodeToString(annv)))
	return changed
}

func updateAnnotation(secret *v1.Secret, value hashString) bool {
	var changed bool
	if secret.Annotations == nil {
		secret.Annotations = make(map[string]string)
	}
	if hashString(secret.Annotations[tlsHashAnnotation]) != value {
		changed = true
		secret.Annotations[tlsHashAnnotation] = string(value)
	}
	return changed
}

func isSecretNeedToBeUpdated(secret *v1.Secret, newDataHash hashBytes) bool {
	if !isValidHashAnnotation(secret) {
		return true
	}
	if !isHashesEquals(secret, newDataHash) {
		return true
	}
	return false
}

func isHashesEquals(secret *v1.Secret, newDataHash hashBytes) bool {
	if secret.Annotations == nil {
		return false
	}
	annotationValue, ok := secret.Annotations[tlsHashAnnotation]
	if !ok {
		return false
	}
	hs := hashString(hex.EncodeToString(newDataHash))
	return hashString(annotationValue) == hs
}

func isValidHashAnnotation(secret *v1.Secret) bool {
	hashData := hex.EncodeToString(hashSecretData(secret))
	hashAnnotation := secret.Annotations[tlsHashAnnotation]
	return hashData == hashAnnotation
}

func hashSecretData(secret *v1.Secret) []byte {
	cert, key := getSecretData(secret)
	if secret.Data != nil {
		return HashCertPlusKey(encodeBase64(cert), encodeBase64(key))
	}
	return []byte{}
}

func encodeBase64(src secretBytes) (dest base64bytes) {
	dest = make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dest, src)
	return
}

func decodeBase64(src base64bytes) (dest secretBytes, err error) {
	dest = make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dest, src)
	if err != nil {
		return dest[:], err
	}
	return dest[:n], nil
}

func getSecretData(secret *v1.Secret) (cert, key secretBytes) {
	cert = secret.Data[tlsDataCert]
	key = secret.Data[tlsDataKey]
	return cert, key
}
