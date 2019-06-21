package pubsubclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	pb "github.com/fiveateooate/deployinator/deployproto"
	sharedfuncs "github.com/fiveateooate/deployinator/internal/common"
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

type messageHandler func(context.Context, *pubsub.Message)

// NewClient - connect the client
func (qcli *PubSubClient) NewClient() {
	var (
		err error
	)
	qcli.CTX, qcli.Cancel = context.WithCancel(context.Background())

	creds, err := google.FindDefaultCredentials(qcli.CTX, pubsub.ScopePubSub)
	if err != nil {
		log.Fatalf("Failed to find credentials: %v", err)
	}
	qcli.cli, err = pubsub.NewClient(qcli.CTX, qcli.ProjectID, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
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

// Exists - returns true if topic exists
func (qcli *PubSubClient) Exists() bool {
	qcli.MyTopic = qcli.cli.Topic(qcli.TopicName)
	exists, err := qcli.MyTopic.Exists(qcli.CTX)
	if err != nil {
		log.Printf("exists error: %v", err)
	}
	return exists
}

// SetTopic - set or create a topic if not exists
func (qcli *PubSubClient) SetTopic() error {
	var (
		err    error
		exists bool
	)
	qcli.MyTopic = qcli.cli.Topic(qcli.TopicName)
	exists, err = qcli.MyTopic.Exists(qcli.CTX)
	if !exists {
		log.Printf("creating topic: %s", qcli.TopicName)
		qcli.MyTopic, err = qcli.cli.CreateTopic(qcli.CTX, qcli.TopicName)
		return err
	}
	// log.Printf("getTopic: %v\n", qcli.MyTopic)
	return nil
}

// Publish - add some something
func (qcli *PubSubClient) Publish(alert *pb.DeployMessage) (string, error) {
	var results []*pubsub.PublishResult
	data, _ := proto.Marshal(alert)
	msgid := ""
	msg := &pubsub.Message{Data: data}

	qcli.MyTopic.PublishSettings = pubsub.PublishSettings{
		ByteThreshold:  5000,
		CountThreshold: 10,
		DelayThreshold: 100 * time.Millisecond,
	}
	result := qcli.MyTopic.Publish(qcli.CTX, msg)
	// log.Println(result)
	results = append(results, result)
	for _, r := range results {
		id, err := r.Get(qcli.CTX)
		if err != nil {
			log.Println(err)
			return "", err
		}
		msgid = id
		// fmt.Printf("Published a message with a message ID: %s\n", id)
	}
	return msgid, nil
}

// PublishResponse - add some something
func (qcli *PubSubClient) PublishResponse(deploystatus *pb.DeployStatusMessage) (string, error) {
	var results []*pubsub.PublishResult
	msgid := ""
	data, _ := proto.Marshal(deploystatus)
	msg := &pubsub.Message{Data: data}

	log.Printf("topic: %v\n", qcli.MyTopic)

	qcli.MyTopic.PublishSettings = pubsub.PublishSettings{
		ByteThreshold:  5000,
		CountThreshold: 10,
		DelayThreshold: 100 * time.Millisecond,
	}

	result := qcli.MyTopic.Publish(qcli.CTX, msg)
	results = append(results, result)
	for _, r := range results {
		id, err := r.Get(qcli.CTX)
		if err != nil {
			log.Println(err)
			return "", err
		}
		msgid = id
	}
	return msgid, nil
}

// Subscribe - do that to a topic
func (qcli *PubSubClient) Subscribe() error {
	qcli.SubName = fmt.Sprintf("%s-%s", qcli.TopicName, sharedfuncs.RandString(12))
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
	return nil
}

// GetMessage - get messages
// blocking
func (qcli *PubSubClient) GetMessage(fn messageHandler) {
	err := qcli.MySub.Receive(qcli.CTX, fn)
	if err != nil {
		log.Println(err)
		return
	}
	return
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

// Stop - close topic/flush messages
func (qcli *PubSubClient) Stop() {
	qcli.MyTopic.Stop()
}

//Delete get rid of topic
func (qcli *PubSubClient) Delete() {
	qcli.MyTopic.Stop()
	qcli.MyTopic.Delete(qcli.CTX)
}

//GetAll - get messages but filter them
func (qcli *PubSubClient) GetAll() []pb.DeployStatusMessage {
	var messages []pb.DeployStatusMessage
	log.Println("here in getall yall")
	if err := qcli.MySub.SeekToTime(qcli.CTX, time.Now().Add(-time.Minute*10)); err != nil {
		log.Fatalln(err)
	}
	log.Println("guess I seeked back yo")
	err := qcli.MySub.Receive(qcli.CTX, func(ctx context.Context, msg *pubsub.Message) {
		var message pb.DeployStatusMessage
		log.Printf("got a message: %v\n", message)
		err := proto.Unmarshal(msg.Data, &message)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		msg.Ack()
		log.Println(message.Status)
		messages = append(messages, message)
	})
	if err != nil {
		log.Println(err)
		return messages
	}
	log.Println("leaving getall yall")
	return messages
}

//GetMsgIDMessages - get messages but filter them
func (qcli *PubSubClient) GetMsgIDMessages(msgid string) pb.DeployStatusMessage {
	var message pb.DeployStatusMessage
	message.Success = false
	if err := qcli.MySub.SeekToTime(qcli.CTX, time.Now().Add(-time.Minute*10)); err != nil {
		log.Fatalln(err)
	}
	err := qcli.MySub.Receive(qcli.CTX, func(ctx context.Context, msg *pubsub.Message) {
		err := proto.Unmarshal(msg.Data, &message)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		if message.MsgID == msgid {
			msg.Ack()
			ctx.Done()
			qcli.Cancel()
		} else {
			msg.Nack()
		}
	})
	if err != nil {
		log.Println(err)
		return message
	}
	return message
}
