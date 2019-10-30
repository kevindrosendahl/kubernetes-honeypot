package kubelet

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

type HoneypotProvider struct {
	store   PodStore
	auditor Auditor
}

func NewHoneypotProvider(store PodStore, auditor Auditor) *HoneypotProvider {
	return &HoneypotProvider{
		store:   store,
		auditor: auditor,
	}
}

func (p *HoneypotProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.auditor.AuditCreatePod(ctx, pod); err != nil {
		return err
	}

	if err := p.store.AddPod(pod); err != nil {
		return err
	}

	return nil
}

func (p *HoneypotProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.auditor.AuditUpdatePod(ctx, pod); err != nil {
		return err
	}

	if err := p.store.UpdatePod(pod); err != nil {
		return err
	}

	return nil
}

func (p *HoneypotProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	return p.store.RemovePod(pod)
}

func (p *HoneypotProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return p.store.GetPod(namespace, name)
}

func (p *HoneypotProvider) GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error) {
	pod, err := p.store.GetPod(namespace, name)
	if err != nil {
		return nil, err
	}
	if pod != nil {
		return nil, nil
	}

	return &pod.Status, nil
}

func (p *HoneypotProvider) GetPods(context.Context) ([]*corev1.Pod, error) {
	return p.store.GetPods()
}
