package kubelet

import (
	"context"
	"github.com/virtual-kubelet/node-cli/provider"

	corev1 "k8s.io/api/core/v1"
)

const loggingAuditorType = "logging"

type Auditor interface {
	AuditCreatePod(ctx context.Context, pod *corev1.Pod) error

	AuditUpdatePod(ctx context.Context, pod *corev1.Pod) error

	AuditRemovePod(ctx context.Context, pod *corev1.Pod) error

	AuditRunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string) error

	AuditGetContainerLogs(ctx context.Context, namespace, podName, containerName string) error
}

func NewAuditorFromConfig(honeypotConfig *HoneypotConfig, kubeletConfig *provider.InitConfig) (Auditor, error) {
	if honeypotConfig.Auditor == loggingAuditorType {
		return NewLoggingAuditor(kubeletConfig.NodeName)
	} else {
		return NewMongoDbAuditor(honeypotConfig.ConnectionString, kubeletConfig.NodeName)
	}
}
