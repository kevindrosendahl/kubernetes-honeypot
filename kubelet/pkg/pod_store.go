package kubelet

import corev1 "k8s.io/api/core/v1"

type PodStore interface {
	AddPod(pod *corev1.Pod) error

	UpdatePod(pod *corev1.Pod) error

	RemovePod(pod *corev1.Pod) error

	GetPod(namespace, name string) (*corev1.Pod, error)

	GetPods() ([]*corev1.Pod, error)
}
