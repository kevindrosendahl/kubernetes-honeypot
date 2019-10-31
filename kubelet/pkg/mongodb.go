package kubelet

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	corev1 "k8s.io/api/core/v1"
)

const (
	kubeletDatabase string = "kubelet"
	auditCollection string = "audit"
)

type MongoDbAuditor struct {
	name   string
	client *mongo.Client
}

type MongoDbAuditorEventType string

const (
	CreateEvent MongoDbAuditorEventType = "create"
	UpdateEvent MongoDbAuditorEventType = "update"
	RemoveEvent MongoDbAuditorEventType = "remove"
)

type MongoDbAuditorCreatePodEvent struct {
	EventType MongoDbAuditorEventType `json:"eventType"`
	Kubelet   string                  `json:"kubelet"`
	Pod       *corev1.Pod             `json:"pod"`
	Timestamp time.Time               `json:"timestamp"`
}

func NewMongoDbAuditor(connectionString, name string) (*MongoDbAuditor, error) {
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
		name:   name,
		client: client,
	}, nil
}

func (a *MongoDbAuditor) AuditCreatePod(ctx context.Context, pod *corev1.Pod) error {
	event := MongoDbAuditorCreatePodEvent{
		EventType: CreateEvent,
		Kubelet:   a.name,
		Pod:       pod,
		Timestamp: time.Now(),
	}

	_, err := a.collection().InsertOne(ctx, &event)
	return err
}

func (a *MongoDbAuditor) AuditUpdatePod(ctx context.Context, pod *corev1.Pod) error {
	event := MongoDbAuditorCreatePodEvent{
		EventType: UpdateEvent,
		Kubelet:   a.name,
		Pod:       pod,
		Timestamp: time.Now(),
	}

	_, err := a.collection().InsertOne(ctx, &event)
	return err
}

func (a *MongoDbAuditor) AuditRemovePod(ctx context.Context, pod *corev1.Pod) error {
	event := MongoDbAuditorCreatePodEvent{
		EventType: RemoveEvent,
		Kubelet:   a.name,
		Pod:       pod,
		Timestamp: time.Now(),
	}

	_, err := a.collection().InsertOne(ctx, &event)
	return err
}

func (a *MongoDbAuditor) collection() *mongo.Collection {
	return a.client.Database(kubeletDatabase).Collection(auditCollection)
}
