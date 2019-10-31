package kubelet

import (
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"

	corev1 "k8s.io/api/core/v1"
)

type HoneypotProvider struct {
	store   PodStore
	auditor Auditor
}

func NewHoneypotProviderFromConfig(cfg *HoneypotConfig, nodeName string) (*HoneypotProvider, error) {
	store, err := NewFileSystemPodStore(cfg.PodStorePath)
	if err != nil {
		return nil, err
	}

	auditor, err := NewMongoDbAuditor(cfg.PodStorePath, nodeName)
	if err != nil {
		return nil, err
	}

	return NewHoneypotProvider(store, auditor), nil
}

func NewHoneypotProvider(store PodStore, auditor Auditor) *HoneypotProvider {
	return &HoneypotProvider{
		store:   store,
		auditor: auditor,
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
		return nil, handleError(err)
	}

	return pod, nil
}

func (p *HoneypotProvider) GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error) {
	pod, err := p.store.GetPod(namespace, name)
	if err != nil {
		return nil, handleError(err)
	}
	if pod != nil {
		return nil, nil
	}

	return &pod.Status, nil
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

func (p *HoneypotProvider) ConfigureNode(context.Context, *corev1.Node) {
	// no-op
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
