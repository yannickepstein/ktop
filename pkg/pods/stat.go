package pods

import (
	"context"
	"fmt"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"github.com/alexeyco/simpletable"
)

type pod = coreV1.Pod
type podMetrics = v1beta1.PodMetrics

type client interface {
	Pods(ctx context.Context, namespace string) ([]pod, error)
	PodMetrices(ctx context.Context, namespace string) ([]podMetrics, error)
}

type newTable func(ctx context.Context, namespace string) (string, error)

func Stats(c client) newTable {
	return func(ctx context.Context, namespace string) (string, error) {
		pods, err := c.Pods(ctx, namespace)
		if err != nil {
			return "", err
		}
		podMetrices, err := c.PodMetrices(ctx, namespace)
		if err != nil {
			return "", err
		}
		return renderTable(extract(pods, podMetrices)), err
	}
}

func renderTable(podStats []podStat) string {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: "Node"},
			{Align: simpletable.AlignLeft, Text: "Namespace"},
			{Align: simpletable.AlignLeft, Text: "Pod"},
			{Align: simpletable.AlignLeft, Text: "CPU (Request)"},
			{Align: simpletable.AlignLeft, Text: "CPU (Limit)"},
			{Align: simpletable.AlignLeft, Text: "CPU (Usage)"},
			{Align: simpletable.AlignLeft, Text: "Memory (Request)"},
			{Align: simpletable.AlignLeft, Text: "Memory (Limit)"},
			{Align: simpletable.AlignLeft, Text: "Memory (Usage)"},
		},
	}
	for _, stat := range podStats {
		row := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: stat.Node},
			{Align: simpletable.AlignLeft, Text: stat.Namespace},
			{Align: simpletable.AlignLeft, Text: stat.Name},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%dm", stat.CPURequest)},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%dm", stat.CPULimit)},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%dm", stat.CPUUsage)},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%dm", stat.MemoryRequest)},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%dm", stat.MemoryLimit)},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%dm", stat.MemoryUsage)},
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}
	table.SetStyle(simpletable.StyleDefault)
	return table.String()
}
