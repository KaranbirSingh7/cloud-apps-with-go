// Package messaging is for components that enable messaging to other systems.
package messaging

import (
	"canvas/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"go.uber.org/zap"
)

type Queue struct {
	Client   *azservicebus.Client
	log      *zap.Logger
	mutex    sync.Mutex
	name     string
	url      *string
	waitTime time.Duration
}

type NewQueueOptions struct {
	Namespace string
	Log       *zap.Logger
	Name      string
	WaitTime  time.Duration
	Client    *azservicebus.Client
}

func NewQueue(opts NewQueueOptions) *Queue {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	return &Queue{
		Client:   opts.Client,
		log:      opts.Log,
		name:     opts.Name,
		waitTime: opts.WaitTime,
	}
}

func GetAzureServiceBusClient(hostname string) *azservicebus.Client {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println("ERROR: Cannot get azure credentials")
		return nil
	}
	client, err := azservicebus.NewClient(hostname, cred, nil)
	if err != nil {
		log.Println("ERROR: unable to create new Azure Service Bus client", err)
		return nil
	}
	return client

}

// Send a message to the queue as JSON.
func (q *Queue) Send(ctx context.Context, m model.Message) error {
	// boilerplate, json marshalling
	messageAsBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	// messageAsString := string(messageAsBytes)

	sender, err := q.Client.NewSender(q.name, nil)
	if err != nil {
		return err
	}

	sbMessage := &azservicebus.Message{
		Body: messageAsBytes,
	}

	err = sender.SendMessage(ctx, sbMessage, nil)
	if err != nil {
		return err
	}
	return nil
}

func (q *Queue) Receive(ctx context.Context) (model.Message, error) {
	var m model.Message

	receiver, err := q.Client.NewReceiverForQueue(q.name, nil)
	if err != nil {
		return nil, err
	}
	defer receiver.Close(ctx)

	messages, err := receiver.ReceiveMessages(ctx, '1', nil)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		body := message.Body
		fmt.Println("Message received: %s", string(body))

		if err := json.Unmarshal(body, &m); err != nil {
			return nil, err
		}

		// mark message as complete
		err = receiver.CompleteMessage(ctx, message, nil)
		if err != nil {
			return nil, err
		}

	}

	return m, nil
}
