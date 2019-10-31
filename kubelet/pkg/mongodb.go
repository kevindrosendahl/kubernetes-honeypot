package kubelet

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	corev1 "k8s.io/api/core/v1"
)

const (
	kubeletDatabase                 string = "kubelet"
	containerLogsRequestsCollection string = "container-logs-requests"
	execRequestsCollection          string = "exec-requests"
	podEventsCollection             string = "pod-events"
)

type MongoDbAuditor struct {
	nodeName string
	client   *mongo.Client
}

type MongoDbAuditorEventType string

const (
	CreateEvent MongoDbAuditorEventType = "create"
	UpdateEvent MongoDbAuditorEventType = "update"
	RemoveEvent MongoDbAuditorEventType = "remove"
)

type mongoDbAuditorPodEvent struct {
	EventType MongoDbAuditorEventType `json:"eventType"`
	Node      string                  `json:"node"`
	Pod       *corev1.Pod             `json:"pod"`
	Timestamp time.Time               `json:"timestamp"`
}

type mongoDbAuditorExecRequest struct {
	Node          string    `json:"node"`
	Namespace     string    `json:"namespace"`
	PodName       string    `json:"podName"`
	ContainerName string    `json:"containerName"`
	Cmd           []string  `json:"cmd"`
	Timestamp     time.Time `json:"timestamp"`
}

type mongoDbAuditorContainerLogsRequest struct {
	Node          string    `json:"node"`
	Namespace     string    `json:"namespace"`
	PodName       string    `json:"podName"`
	ContainerName string    `json:"containerName"`
	Timestamp     time.Time `json:"timestamp"`
}

func NewMongoDbAuditor(connectionString, nodeName string) (*MongoDbAuditor, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return &MongoDbAuditor{
		nodeName: nodeName,
		client:   client,
	}, nil
}

func (a *MongoDbAuditor) AuditCreatePod(ctx context.Context, pod *corev1.Pod) error {
	event := mongoDbAuditorPodEvent{
		EventType: CreateEvent,
		Node:      a.nodeName,
		Pod:       pod,
		Timestamp: time.Now(),
	}

	_, err := a.podEventsCollection().InsertOne(ctx, &event)
	return err
}

func (a *MongoDbAuditor) AuditUpdatePod(ctx context.Context, pod *corev1.Pod) error {
	event := mongoDbAuditorPodEvent{
		EventType: UpdateEvent,
		Node:      a.nodeName,
		Pod:       pod,
		Timestamp: time.Now(),
	}

	_, err := a.podEventsCollection().InsertOne(ctx, &event)
	return err
}

func (a *MongoDbAuditor) AuditRemovePod(ctx context.Context, pod *corev1.Pod) error {
	event := mongoDbAuditorPodEvent{
		EventType: RemoveEvent,
		Node:      a.nodeName,
		Pod:       pod,
		Timestamp: time.Now(),
	}

	_, err := a.podEventsCollection().InsertOne(ctx, &event)
	return err
}

func (a *MongoDbAuditor) AuditRunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string) error {
	request := mongoDbAuditorExecRequest{
		Node:          a.nodeName,
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		Cmd:           cmd,
		Timestamp:     time.Now(),
	}

	_, err := a.execRequestsCollection().InsertOne(ctx, &request)
	return err
}

func (a *MongoDbAuditor) AuditGetContainerLogs(ctx context.Context, namespace, podName, containerName string) error {
	request := mongoDbAuditorContainerLogsRequest{
		Node:          a.nodeName,
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		Timestamp:     time.Now(),
	}

	_, err := a.containerLogsRequestsCollection().InsertOne(ctx, &request)
	return err
}

func (a *MongoDbAuditor) containerLogsRequestsCollection() *mongo.Collection {
	return a.client.Database(kubeletDatabase).Collection(containerLogsRequestsCollection)
}

func (a *MongoDbAuditor) execRequestsCollection() *mongo.Collection {
	return a.client.Database(kubeletDatabase).Collection(execRequestsCollection)
}

func (a *MongoDbAuditor) podEventsCollection() *mongo.Collection {
	return a.client.Database(kubeletDatabase).Collection(podEventsCollection)
}
