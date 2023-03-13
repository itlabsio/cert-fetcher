package kubeSecreter

import (
	"testing"
  "context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type secret struct{}

func TestConnect(t *testing.T) {
	k, err := AuthByDefaultKubeconfig()
	if err != nil {
		t.Skip()
	}
	_, err = k.clientSet.CoreV1().Namespaces().Get(context.TODO(), "default", metav1.GetOptions{})
	if err != nil {
		t.Fail()
	}
}
