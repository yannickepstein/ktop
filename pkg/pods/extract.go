package pods

import (
	resourcehelper "k8s.io/kubectl/pkg/util/resource"
)

type podStat struct {
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

func extract(pods []pod, podMetrices []podMetrics) []podStat {
	podStats := make([]podStat, 0)
	for _, podMetrics := range podMetrices {
		podStats = append(podStats, transform(pods, podMetrics))
	}
	return podStats
}

func transform(pods []pod, metrics podMetrics) podStat {
	var matchingPod pod
	for _, it := range pods {
		if it.Namespace == metrics.Namespace && it.Name == metrics.Name {
			matchingPod = it
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
	for _, container := range metrics.Containers {
		cpuUsage += container.Usage.Cpu().ToDec().MilliValue()
		memoryUsage += container.Usage.Memory().MilliValue()
	}

	return podStat{
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
