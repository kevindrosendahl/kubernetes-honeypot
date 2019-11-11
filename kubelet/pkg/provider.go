package kubelet

import (
	"context"
	"fmt"
	"github.com/virtual-kubelet/virtual-kubelet/errdefs"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/virtual-kubelet/node-cli/provider"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HoneypotProvider struct {
	store          PodStore
	auditor        Auditor
	honeypotConfig *HoneypotConfig
	kubeletConfig  *provider.InitConfig
}

func NewHoneypotProviderFromConfig(honeypotConfig *HoneypotConfig, kubeletConfig *provider.InitConfig) (*HoneypotProvider, error) {
	store, err := NewFileSystemPodStore(honeypotConfig.PodStorePath)
	if err != nil {
		return nil, err
	}

	auditor, err := NewAuditorFromConfig(honeypotConfig, kubeletConfig)
	if err != nil {
		return nil, err
	}

	return NewHoneypotProvider(store, auditor, honeypotConfig, kubeletConfig), nil
}

func NewHoneypotProvider(
	store PodStore,
	auditor Auditor,
	honeypotConfig *HoneypotConfig,
	kubeletConfig *provider.InitConfig,
) *HoneypotProvider {
	return &HoneypotProvider{
		store:          store,
		auditor:        auditor,
		honeypotConfig: honeypotConfig,
		kubeletConfig:  kubeletConfig,
	}
}

func (p *HoneypotProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.auditor.AuditCreatePod(ctx, pod); err != nil {
		return handleError(err)
	}

	if err := p.store.AddPod(pod); err != nil {
		return handleError(err)
	}

	return nil
}

func (p *HoneypotProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.auditor.AuditUpdatePod(ctx, pod); err != nil {
		return handleError(err)
	}

	if err := p.store.UpdatePod(pod); err != nil {
		return handleError(err)
	}

	return nil
}

func (p *HoneypotProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.auditor.AuditRemovePod(ctx, pod); err != nil {
		return handleError(err)
	}

	if err := p.store.RemovePod(pod); err != nil {
		return handleError(err)
	}

	return nil
}

func (p *HoneypotProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	pod, err := p.store.GetPod(namespace, name)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return nil, err
		}

		return nil, handleError(err)
	}

	return pod, nil
}

func (p *HoneypotProvider) GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error) {
	pod, err := p.store.GetPod(namespace, name)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return nil, err
		}

		return nil, handleError(err)
	}
	if pod == nil {
		return nil, nil
	}

	containerStatuses := make([]corev1.ContainerStatus, len(pod.Spec.Containers))
	for _, container := range pod.Spec.Containers {
		status := corev1.ContainerStatus{
			Name: container.Name,
			State: corev1.ContainerState{
				Running: &corev1.ContainerStateRunning{StartedAt: metav1.Now()},
			},
			Ready: true,
			Image: container.Image,
		}

		containerStatuses = append(containerStatuses, status)
	}

	return &corev1.PodStatus{
		Phase:             corev1.PodRunning,
		ContainerStatuses: containerStatuses,
	}, nil
}

func (p *HoneypotProvider) GetPods(context.Context) ([]*corev1.Pod, error) {
	pods, err := p.store.GetPods()
	if err != nil {
		return nil, handleError(err)
	}

	return pods, nil
}

func (p *HoneypotProvider) GetContainerLogs(ctx context.Context, namespace, podName, containerName string, opts api.ContainerLogOpts) (io.ReadCloser, error) {
	if err := p.auditor.AuditGetContainerLogs(ctx, namespace, podName, containerName); err != nil {
		return nil, handleError(err)
	}

	return nil, internalError()
}

func (p *HoneypotProvider) RunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string, attach api.AttachIO) error {
	if err := p.auditor.AuditRunInContainer(ctx, namespace, podName, containerName, cmd); err != nil {
		return handleError(err)
	}

	return internalError()
}

func (p *HoneypotProvider) ConfigureNode(ctx context.Context, node *corev1.Node) {
	node.Status.Capacity = p.capacity()
	node.Status.Allocatable = p.capacity()
	node.Status.Conditions = p.nodeConditions()
	node.Status.Addresses = p.nodeAddresses()
	node.Status.DaemonEndpoints = p.nodeDaemonEndpoints()
	os := p.kubeletConfig.OperatingSystem
	if os == "" {
		os = "Linux"
	}
	node.Status.NodeInfo.OperatingSystem = os
	node.Status.NodeInfo.Architecture = "amd64"
	node.ObjectMeta.Labels = make(map[string]string)
}

func (p *HoneypotProvider) capacity() corev1.ResourceList {
	return corev1.ResourceList{
		"cpu":    resource.MustParse(p.honeypotConfig.Capacity.Cpu),
		"memory": resource.MustParse(p.honeypotConfig.Capacity.Memory),
		"pods":   resource.MustParse(p.honeypotConfig.Capacity.Pods),
	}
}

func (p *HoneypotProvider) nodeConditions() []corev1.NodeCondition {
	return []corev1.NodeCondition{
		{
			Type:               "Ready",
			Status:             corev1.ConditionTrue,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletReady",
			Message:            "kubelet is ready.",
		},
		{
			Type:               "OutOfDisk",
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientDisk",
			Message:            "kubelet has sufficient disk space available",
		},
		{
			Type:               "MemoryPressure",
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientMemory",
			Message:            "kubelet has sufficient memory available",
		},
		{
			Type:               "DiskPressure",
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasNoDiskPressure",
			Message:            "kubelet has no disk pressure",
		},
		{
			Type:               "NetworkUnavailable",
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "RouteCreated",
			Message:            "RouteController created a route",
		},
	}
}

func (p *HoneypotProvider) nodeAddresses() []corev1.NodeAddress {
	return []corev1.NodeAddress{
		{
			Type:    "InternalIP",
			Address: p.kubeletConfig.InternalIP,
		},
	}
}

func (p *HoneypotProvider) nodeDaemonEndpoints() corev1.NodeDaemonEndpoints {
	return corev1.NodeDaemonEndpoints{
		KubeletEndpoint: corev1.DaemonEndpoint{
			Port: p.kubeletConfig.DaemonPort,
		},
	}
}

func handleError(err error) error {
	// We don't want to leak any info about this kubelet, but do want to know when errors occur.
	// So we'll log the error but then just return a generic internal error.
	log.Errorf("encountered error: %s", err.Error())
	return internalError()
}

func internalError() error {
	// Return an error message that looks like there was an internal error so we don't leak any
	// info that may reveal to the attacker that this is a honeypot.
	return fmt.Errorf("internal error")
}
