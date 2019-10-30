package kubelet

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

type Auditor interface {
	CreatePod(ctx context.Context, pod *corev1.Pod) error

	UpdatePod(ctx context.Context, pod *corev1.Pod) error
}
