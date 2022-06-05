package kapi

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	// use for vendor specific authentication
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type Client struct {
	Metrics   *metricsv.Clientset
	Resources *kubernetes.Clientset
}

func NewClient() (*Client, error) {
	kubeconfig, err := newKubeconfig()
	if err != nil {
		return nil, err
	}
	metrics, err := metricsv.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	return &Client{
		Metrics:   metrics,
		Resources: clientset,
	}, nil
}

func newKubeconfig() (*rest.Config, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
