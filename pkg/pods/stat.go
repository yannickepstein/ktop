package pods

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/yannickepstein/ktop/pkg/kapi"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Display(ctx context.Context, client *kapi.Client) error {
	namespaceList, err := client.Resources.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}
	namespaces := namespaceList.Items[0:10]
	podStats := make(chan []PodStat, len(namespaces))
	var wg sync.WaitGroup
	wg.Add(len(namespaces))
	for _, namespace := range namespaces {
		go func(nsName string) {
			defer wg.Done()
			podList, err := client.Resources.CoreV1().Pods(nsName).List(ctx, v1.ListOptions{})
			if err != nil {
				return
			}
			podMetricsList, err := client.Metrics.MetricsV1beta1().PodMetricses(nsName).List(ctx, v1.ListOptions{})
			if err != nil {
				return
			}
			podStats <- extract(podList, podMetricsList)
		}(namespace.Name)
	}
	wg.Wait()
	fmt.Sprintln("completed wait")
	for i := 0; i < len(namespaces); i++ {
		podStats := <-podStats
		render(podStats)
	}
	return nil
}

func render(podStats []PodStat) {
	for _, podStat := range podStats {
		columns := []string{
			podStat.Node,
			podStat.Namespace,
			podStat.Name,
			fmt.Sprintf("%dm", podStat.CPURequest),
			fmt.Sprintf("%dm", podStat.CPULimit),
			fmt.Sprintf("%dm", podStat.CPUUsage),
			fmt.Sprintf("%dm", podStat.MemoryRequest),
			fmt.Sprintf("%dm", podStat.MemoryLimit),
			fmt.Sprintf("%dm", podStat.MemoryUsage),
		}
		fmt.Println(strings.Join(columns[:], "\t "))
	}
}
