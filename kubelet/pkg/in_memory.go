package kubelet

import (
	"encoding/json"
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
)

type InMemoryPodStore struct {
	pods map[string]*corev1.Pod
	lock sync.Mutex
}

func NewInMemoryPodStore() *InMemoryPodStore {
	pods := make(map[string]*corev1.Pod)
	return &InMemoryPodStore{pods: pods, lock: sync.Mutex{}}
}

func NewInMemoryPodStoreWithPods(pods map[string]*corev1.Pod) *InMemoryPodStore {
	store := NewInMemoryPodStore()
	store.pods = pods
	return store
}

func NewInMemoryPodStoreFromJson(jsonBytes []byte) (*InMemoryPodStore, error) {
	pods := make(map[string]*corev1.Pod)
	if err := json.Unmarshal(jsonBytes, &pods); err != nil {
		return nil, err
	}

	return NewInMemoryPodStoreWithPods(pods), nil
}

func (s *InMemoryPodStore) AddPod(pod *corev1.Pod) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.pods[podKey(pod)] = pod
	return nil
}

func (s *InMemoryPodStore) UpdatePod(pod *corev1.Pod) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.pods[podKey(pod)] = pod
	return nil
}

func (s *InMemoryPodStore) RemovePod(pod *corev1.Pod) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.pods, podKey(pod))
	return nil
}

func (s *InMemoryPodStore) GetPod(namespace, name string) (*corev1.Pod, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.pods[key(namespace, name)], nil
}

func (s *InMemoryPodStore) GetPods() ([]*corev1.Pod, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	pods := make([]*corev1.Pod, 0, len(s.pods))
	for _, pod := range s.pods {
		pods = append(pods, pod)
	}

	return pods, nil
}

func (s *InMemoryPodStore) ToJson() ([]byte, error) {
	return json.Marshal(s.pods)
}

func podKey(pod *corev1.Pod) string {
	return key(pod.Namespace, pod.Name)
}

func key(namespace, name string) string {
	return fmt.Sprintf("%s.%s", namespace, name)
}
