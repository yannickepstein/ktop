package kapi

import (
	"context"
	"os"
	"path/filepath"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	// use for vendor specific authentication
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type client struct {
	k8Clientset *kubernetes.Clientset
	metrics     *metricsv.Clientset
}

func (c client) Pods(ctx context.Context, namespace string) ([]coreV1.Pod, error) {
	podList, err := c.k8Clientset.CoreV1().Pods(namespace).List(ctx, metaV1.ListOptions{})
	if err != nil {
		return []coreV1.Pod{}, err
	}
	return podList.Items, nil
}

func (c client) PodMetrices(ctx context.Context, namespace string) ([]v1beta1.PodMetrics, error) {
	metricsList, err := c.metrics.MetricsV1beta1().PodMetricses(namespace).List(ctx, metaV1.ListOptions{})
	if err != nil {
		return []v1beta1.PodMetrics{}, err
	}
	return metricsList.Items, nil
}

func NewClient() (*client, error) {
	kubeconfig, err := newKubeconfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	metricsClientset, err := metricsv.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	return &client{
		k8Clientset: clientset,
		metrics:     metricsClientset,
	}, nil
}

func newKubeconfig() (*rest.Config, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
