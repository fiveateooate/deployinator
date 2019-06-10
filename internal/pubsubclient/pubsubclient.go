package pubsubclient

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"cloud.google.com/go/pubsub"
	pb "github.com/fiveateooate/deployinator/deployproto"
	"github.com/golang/protobuf/proto"
	google "golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// PubSubClient holds queue related stuff
type PubSubClient struct {
	TopicName      string
	SubName        string
	ProjectID      string
	Cancel         context.CancelFunc
	requestTimeout time.Duration
	CTX            context.Context
	cli            *pubsub.Client
	MyTopic        *pubsub.Topic
	MySub          *pubsub.Subscription
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

// Connect - connect the client
func (qcli *PubSubClient) Connect() {
	var (
		err error
	)
	qcli.CTX, qcli.Cancel = context.WithCancel(context.Background())
	qcli.SubName = qcli.subName()

	creds, err := google.FindDefaultCredentials(qcli.CTX, pubsub.ScopePubSub)
	if err != nil {
		log.Fatalf("Failed to find credentials: %v", err)
	}
	qcli.cli, err = pubsub.NewClient(qcli.CTX, qcli.ProjectID, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	qcli.getTopic()
}

// List - idk ... things
func (qcli *PubSubClient) List() {
	it := qcli.cli.Topics(qcli.CTX)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(topic)
	}
}

func (qcli *PubSubClient) exists() bool {
	qcli.MyTopic = qcli.cli.Topic(qcli.TopicName)
	exists, _ := qcli.MyTopic.Exists(qcli.CTX)
	return exists
}

// GetTopic - create a topic if not exists
func (qcli *PubSubClient) getTopic() error {
	var (
		err    error
		exists bool
	)
	qcli.MyTopic = qcli.cli.Topic(qcli.TopicName)
	exists, err = qcli.MyTopic.Exists(qcli.CTX)
	if err != nil {
		return err
	}
	if !exists {
		qcli.MyTopic, err = qcli.cli.CreateTopic(qcli.CTX, qcli.TopicName)
		return err
	}
	log.Println(qcli.MyTopic)
	return nil
}

// Publish - add some something
func (qcli *PubSubClient) Publish(alert *pb.DeployMessage) {
	data, _ := proto.Marshal(alert)
	msg := &pubsub.Message{Data: data}

	qcli.MyTopic.PublishSettings = pubsub.PublishSettings{
		ByteThreshold:  5000,
		CountThreshold: 10,
		DelayThreshold: 100 * time.Millisecond,
	}
	result := qcli.MyTopic.Publish(qcli.CTX, msg)
	log.Println(result)
}

func (qcli *PubSubClient) subName() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	n := 12
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(s1.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return fmt.Sprintf("%s-sub-%s", qcli.TopicName, string(b))
}

// Subscribe - do that to a topic
func (qcli *PubSubClient) Subscribe() error {
	log.Printf("topic: %s\n", qcli.TopicName)
	qcli.MySub = qcli.cli.Subscription(qcli.SubName)
	exists, err := qcli.MySub.Exists(qcli.CTX)
	if err != nil {
		log.Printf("Error checking for subscription: %v\n", err)
		return err
	}
	if !exists {
		if _, err = qcli.cli.CreateSubscription(qcli.CTX, qcli.SubName, pubsub.SubscriptionConfig{Topic: qcli.MyTopic}); err != nil {
			log.Fatalf("Failed to create subscription: %v", err)
		}
	}
	log.Println(qcli.MySub)
	return nil
}

func processMessage(ctx context.Context, msg *pubsub.Message) {
	var message pb.DeployMessage
	err := proto.Unmarshal(msg.Data, &message)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	log.Printf("Name: %s, Namespace: %s.\n", message.Name, message.Namespace)
	msg.Ack()
	return
}

// GetMessage - get a message?
func (qcli *PubSubClient) GetMessage() {
	err := qcli.MySub.Receive(qcli.CTX, processMessage)
	if err != nil {
		log.Println(err)
		return
	}
}

// Disconnect delete subscription
func (qcli *PubSubClient) Disconnect() {
	if err := qcli.MySub.Delete(qcli.CTX); err != nil {
		log.Println(err)
		return
	}
	qcli.Cancel()
	log.Printf("Subscription %s deleted", qcli.SubName)
}
