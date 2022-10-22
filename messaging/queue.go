// Package messaging is for components that enable messaging to other systems.
package messaging

import (
	"canvas/model"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"go.uber.org/zap"
)

type Queue struct {
	Client *azservicebus.Client
	log    *zap.Logger
	name   string
}

type NewQueueOptions struct {
	Namespace   string
	Log         *zap.Logger
	Name        string
	KeyName     string
	KeyPassword string
}

func NewQueue(opts NewQueueOptions) *Queue {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	return &Queue{
		Client: opts.NewAzureServiceBusClient(),
		log:    opts.Log,
		name:   opts.Name,
	}
}

func (opts *NewQueueOptions) NewAzureServiceBusClient() *azservicebus.Client {
	connectionString := fmt.Sprintf(
		"Endpoint=sb://%s.servicebus.windows.net/;SharedAccessKeyName=%s;SharedAccessKey=%s", opts.Namespace, opts.KeyName, opts.KeyPassword,
	)

	client, err := azservicebus.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil
	}

	return client

}

// Send a message to the queue as JSON.
func (q *Queue) Send(ctx context.Context, m model.Message) error {
	q.log.Sugar().Debugf("sending message %v to queue %q", m, q.name)
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

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
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
