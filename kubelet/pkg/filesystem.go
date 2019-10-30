package kubelet

import (
	"io/ioutil"
	"os"
	"sync"

	corev1 "k8s.io/api/core/v1"
)

type FileSystemPodStore struct {
	filename string
	store    *InMemoryPodStore
	lock     sync.Mutex
}

func NewFileSystemPodStore(filename string) (*FileSystemPodStore, error) {
	store := &FileSystemPodStore{
		filename: filename,
		store:    NewInMemoryPodStore(),
		lock:     sync.Mutex{},
	}

	// Check to see if the supplied file exists.
	// If it doesn't, then don't try to load the store from it.
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}

		return nil, err
	}

	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *FileSystemPodStore) AddPod(pod *corev1.Pod) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.store.AddPod(pod); err != nil {
		return err
	}

	if err := s.flush(); err != nil {
		return err
	}

	return nil
}

func (s *FileSystemPodStore) UpdatePod(pod *corev1.Pod) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.store.AddPod(pod); err != nil {
		return err
	}

	if err := s.flush(); err != nil {
		return err
	}

	return nil
}

func (s *FileSystemPodStore) RemovePod(pod *corev1.Pod) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.store.RemovePod(pod); err != nil {
		return err
	}

	if err := s.flush(); err != nil {
		return err
	}

	return nil
}

func (s *FileSystemPodStore) GetPod(namespace, name string) (*corev1.Pod, error) {
	return s.store.GetPod(namespace, name)
}

func (s *FileSystemPodStore) GetPods() ([]*corev1.Pod, error) {
	return s.store.GetPods()
}

func (s *FileSystemPodStore) flush() error {
	json, err := s.store.ToJson()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.filename, json, 0644)
}

func (s *FileSystemPodStore) load() error {
	json, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return err
	}

	store, err := NewInMemoryPodStoreFromJson(json)
	if err != nil {
		return err
	}

	s.store = store
	return nil
}
