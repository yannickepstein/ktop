package pods

import (
	v1 "k8s.io/api/core/v1"
	resourcehelper "k8s.io/kubectl/pkg/util/resource"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type PodStat struct {
	Node          string
	Namespace     string
	Name          string
	CPURequest    int64
	CPULimit      int64
	CPUUsage      int64
	MemoryRequest int64
	MemoryLimit   int64
	MemoryUsage   int64
}

func extract(podList *v1.PodList, podMetricsList *v1beta1.PodMetricsList) []PodStat {
	podStats := make([]PodStat, 0)
	for _, podMetrics := range podMetricsList.Items {
		podStats = append(podStats, transform(podList, podMetrics))
	}
	return podStats
}

func transform(podList *v1.PodList, podMetrics v1beta1.PodMetrics) PodStat {
	var matchingPod v1.Pod
	for _, pod := range podList.Items {
		if pod.Namespace == podMetrics.Namespace && pod.Name == podMetrics.Name {
			matchingPod = pod
			break
		}
	}

	reqs, limits := resourcehelper.PodRequestsAndLimits(&matchingPod)
	cpuRequest := reqs.Cpu().MilliValue()
	cpuLimit := limits.Cpu().MilliValue()
	memoryRequest := reqs.Memory().Value() / (1024 * 1024)
	memoryLimit := limits.Memory().Value() / (1024 * 1024)

	var cpuUsage int64 = 0
	var memoryUsage int64 = 0
	for _, container := range podMetrics.Containers {
		cpuUsage += container.Usage.Cpu().ToDec().MilliValue()
		memoryUsage += container.Usage.Memory().MilliValue()
	}

	return PodStat{
		Node:          matchingPod.Spec.NodeName,
		Namespace:     matchingPod.Namespace,
		Name:          matchingPod.Name,
		CPURequest:    cpuRequest,
		CPULimit:      cpuLimit,
		CPUUsage:      cpuUsage,
		MemoryRequest: memoryRequest,
		MemoryLimit:   memoryLimit,
		MemoryUsage:   memoryUsage,
	}
}
