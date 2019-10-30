package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubelet "github.com/kevindrosendahl/kubernetes-honeypot/kubelet/pkg"
)

func main() {
	store, err := kubelet.NewFileSystemPodStore(os.Args[1])
	if err != nil {
		panic(err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	auditor, err := kubelet.NewMongoDbAuditor(client, "local")
	if err != nil {
		panic(err)
	}

	provider := kubelet.NewHoneypotProvider(store, auditor)

	err = provider.CreatePod(context.Background(), &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "foo", Name: "bar"}})
	if err != nil {
		panic(err)
	}

	pods, err := provider.GetPods(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(pods)
}
