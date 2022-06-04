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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list := make(chan string)
	go func() {
		opts := metav1.ListOptions{}
		nodeList, err := api.Nodes().List(ctx, opts)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		list <- listNodes(nodeList)
	}()
	go func() {
		opts := metav1.ListOptions{}
		// TODO: namespace as arg
		podList, err := api.Pods("notifications").List(ctx, opts)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		list <- listPods(podList)
	}()

	for {
		select {
		case out := <-list:
			fmt.Fprint(os.Stdout, out)
		case <-ctx.Done():
			return
		}
	}
}

func listNodes(nodeList *v1.NodeList) string {
	nodes := ""
	for _, node := range nodeList.Items {
		nodes = nodes + fmt.Sprintln(node.Name, node.Status.NodeInfo.Architecture)
	}
	return nodes
}

func listPods(podList *v1.PodList) string {
	pods := ""
	pods = pods + fmt.Sprintln("NAME", "HOST PID")
	for _, pod := range podList.Items {
		pods = pods + fmt.Sprintln(pod.Name, pod.Spec.HostPID)
	}
	return pods
}
