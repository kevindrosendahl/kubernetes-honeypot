package main

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubelet "github.com/kevindrosendahl/kubernetes-honeypot/kubelet/pkg"
)

func main() {
	// Create a new store and add a pod.
	store, err := kubelet.NewFileSystemPodStore(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = store.AddPod(&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "foo", Name: "bar"}})
	if err != nil {
		panic(err)
	}

	// Re-open the store and print out the pods.
	store, err = kubelet.NewFileSystemPodStore(os.Args[1])
	if err != nil {
		panic(err)
	}

	pods, err := store.GetPods()
	if err != nil {
		panic(err)
	}

	fmt.Println(pods)
}
