package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	// use for vendor specific authentication
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	api := clientset.CoreV1()
	metricsClientset, err := metricsv.NewForConfig(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	chList := make(chan string)
	chErr := make(chan error)
	go func() {
		opts := metav1.ListOptions{}
		nodeList, err := api.Nodes().List(ctx, opts)
		if err != nil {
			chErr <- err
			return
		}
		nodeMetricsList, err := metricsClientset.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
		if err != nil {
			chErr <- err
			return
		}
		chList <- listNodes(nodeList, nodeMetricsList)
	}()
	go func() {
		opts := metav1.ListOptions{}
		namespace := "notifications"
		// TODO: namespace as arg
		podList, err := api.Pods(namespace).List(ctx, opts)
		if err != nil {
			chErr <- err
			return
		}
		podMetricsList, err := metricsClientset.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			chErr <- err
			return
		}
		chList <- listPods(podList, podMetricsList)
	}()

	for {
		select {
		case out := <-chList:
			fmt.Fprint(os.Stdout, out)
		case err := <-chErr:
			fmt.Fprintln(os.Stderr, err)
		case <-ctx.Done():
			return
		}
	}
}

func listNodes(nodeList *v1.NodeList, nodeMetricsList *v1beta1.NodeMetricsList) string {
	nodes := ""
	for _, node := range nodeList.Items {
		nodes = nodes + fmt.Sprintln(node.Name, node.Status.NodeInfo.Architecture)
	}
	return nodes
}

func listPods(podList *v1.PodList, podMetricsList *v1beta1.PodMetricsList) string {
	pods := ""
	for _, pod := range podList.Items {
		pods += fmt.Sprintf("%s %s", pod.Name, pod.Spec.NodeName)
	}
	for _, podMetric := range podMetricsList.Items {
		podString := podMetric.GetName()
		for _, container := range podMetric.Containers {
			podString = podString + "\n\t" + fmt.Sprintf("%s: %vm %vMi", container.Name, container.Usage.Cpu().MilliValue(), container.Usage.Memory().Value()/(1024*1024))
		}
		pods = pods + podString + "\n"
	}
	return pods
}
