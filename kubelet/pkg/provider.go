package kubelet

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

type HoneypotProvider struct {
	podStore PodStore
	auditor  Auditor
}

func (p *HoneypotProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.auditor.CreatePod(ctx, pod); err != nil {
		return err
	}

	if err := p.podStore.AddPod(pod); err != nil {
		return err
	}

	return nil
}

func (p *HoneypotProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.auditor.UpdatePod(ctx, pod); err != nil {
		return err
	}

	if err := p.podStore.UpdatePod(pod); err != nil {
		return err
	}

	return nil
}

func (p *HoneypotProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	return p.podStore.RemovePod(pod)
}

func (p *HoneypotProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return p.podStore.GetPod(namespace, name)
}

func (p *HoneypotProvider) GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error) {
	pod, err := p.podStore.GetPod(namespace, name)
	if err != nil {
		return nil, err
	}
	if pod != nil {
		return nil, nil
	}

	return &pod.Status, nil
}

func (p *HoneypotProvider) GetPods(context.Context) ([]*corev1.Pod, error) {
	return p.podStore.GetPods()
}
