package kubeSecreter

import (
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	clientSet   *kubernetes.Clientset
	kubeContext string
}

func (k KubeClient) GetClientSet() *kubernetes.Clientset {
	return k.clientSet
}

func (k KubeClient) getKubeContext() string {
	return k.kubeContext
}

func AuthByDefaultKubeconfig() (*KubeClient, error) {
	kubeconfig := clientcmd.NewDefaultPathOptions().GetDefaultFilename()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	r, _ := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	return &KubeClient{
		clientSet:   kubeclient,
		kubeContext: r.CurrentContext,
	}, nil
}

func AuthInCluster() (*KubeClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	kubeclient, err := kubernetes.NewForConfig(config)
	return &KubeClient{
		clientSet:   kubeclient,
		kubeContext: "",
	}, err
}
