package kubelet

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"time"
)

type LoggingAuditor struct {
	nodeName string
}

func NewLoggingAuditor(nodeName string) (*LoggingAuditor, error) {
	return &LoggingAuditor{
		nodeName: nodeName,
	}, nil
}

func (a *LoggingAuditor) AuditCreatePod(ctx context.Context, pod *corev1.Pod) error {
	event := mongoDbAuditorPodEvent{
		EventType: CreateEvent,
		Node:      a.nodeName,
		Pod:       pod,
		Timestamp: time.Now(),
	}
	fmt.Printf("%v+\n", event)
	return nil
}

func (a *LoggingAuditor) AuditUpdatePod(ctx context.Context, pod *corev1.Pod) error {
	event := mongoDbAuditorPodEvent{
		EventType: UpdateEvent,
		Node:      a.nodeName,
		Pod:       pod,
		Timestamp: time.Now(),
	}

	fmt.Printf("%v+\n", event)
	return nil
}

func (a *LoggingAuditor) AuditRemovePod(ctx context.Context, pod *corev1.Pod) error {
	event := mongoDbAuditorPodEvent{
		EventType: RemoveEvent,
		Node:      a.nodeName,
		Pod:       pod,
		Timestamp: time.Now(),
	}

	fmt.Printf("%v+\n", event)
	return nil
}

func (a *LoggingAuditor) AuditRunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string) error {
	request := mongoDbAuditorExecRequest{
		Node:          a.nodeName,
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		Cmd:           cmd,
		Timestamp:     time.Now(),
	}

	fmt.Printf("%v+\n", request)
	return nil
}

func (a *LoggingAuditor) AuditGetContainerLogs(ctx context.Context, namespace, podName, containerName string) error {
	request := mongoDbAuditorContainerLogsRequest{
		Node:          a.nodeName,
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		Timestamp:     time.Now(),
	}

	fmt.Printf("%v+\n", request)
	return nil
}
