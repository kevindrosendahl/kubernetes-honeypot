package kubelet

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

type Auditor interface {
	AuditCreatePod(ctx context.Context, pod *corev1.Pod) error

	AuditUpdatePod(ctx context.Context, pod *corev1.Pod) error

	AuditRemovePod(ctx context.Context, pod *corev1.Pod) error

	AuditRunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string) error

	AuditGetContainerLogs(ctx context.Context, namespace, podName, containerName string) error
}
