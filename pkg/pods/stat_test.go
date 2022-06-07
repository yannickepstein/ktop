package pods_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/yannickepstein/ktop/pkg/pods"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type pod = coreV1.Pod
type podMetrics = v1beta1.PodMetrics

type client struct{}

func (c client) Pods(ctx context.Context, namespace string) ([]pod, error) {
	return []pod{
		newPod(namespace),
	}, nil
}

func (c client) PodMetrices(ctx context.Context, namespace string) ([]podMetrics, error) {
	return []podMetrics{
		newPodMetrics(namespace),
	}, nil
}

func TestStats(t *testing.T) {
	t.Run("Table with one entry per Pod", func(t *testing.T) {
		namespace := "default"
		table, err := pods.Stats(client{})(context.Background(), namespace)

		if err != nil {
			t.Error("expected no error but got", err)
		}
		if table == "" {
			t.Error("expected non-empty table")
		}
		if !strings.Contains(table, namespace) {
			t.Error("expected namespace to be present in table")
		}
	})
}

func newPod(namespace string) pod {
	return pod{
		TypeMeta:   metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod", Namespace: namespace},
		Spec: coreV1.PodSpec{
			Containers: []coreV1.Container{
				{
					Name:                     "test-container",
					Image:                    "test-image",
					Command:                  []string{},
					Args:                     []string{},
					WorkingDir:               "",
					Ports:                    []coreV1.ContainerPort{},
					EnvFrom:                  []coreV1.EnvFromSource{},
					Env:                      []coreV1.EnvVar{},
					Resources:                coreV1.ResourceRequirements{},
					VolumeMounts:             []coreV1.VolumeMount{},
					VolumeDevices:            []coreV1.VolumeDevice{},
					LivenessProbe:            &coreV1.Probe{},
					ReadinessProbe:           &coreV1.Probe{},
					StartupProbe:             &coreV1.Probe{},
					Lifecycle:                &coreV1.Lifecycle{},
					TerminationMessagePath:   "",
					TerminationMessagePolicy: "",
					ImagePullPolicy:          "Always",
					SecurityContext:          &coreV1.SecurityContext{},
					Stdin:                    false,
					StdinOnce:                false,
					TTY:                      false,
				},
			},
		},
		Status: coreV1.PodStatus{},
	}
}

func newPodMetrics(namespace string) podMetrics {
	return podMetrics{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: namespace,
		},
		Timestamp: metav1.NewTime(time.Now()),
		Window:    metav1.Duration{Duration: 10 * time.Second},
		Containers: []v1beta1.ContainerMetrics{
			{
				Name:  "test-container",
				Usage: coreV1.ResourceList{},
			},
		},
	}
}
