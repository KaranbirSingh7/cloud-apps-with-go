package messaging_test

import (
	"canvas/integrationtest"
	"canvas/model"
	"context"
	"testing"

	"github.com/matryer/is"
)

func TestQueue(t *testing.T) {
	integrationtest.SkipIfShort(t) // -short flag meaning skip this test

	t.Run("sends a message to the queue, receives it, and deletes it", func(t *testing.T) {
		is := is.New(t)

		queue, cleanup := integrationtest.CreateQueue()
		defer cleanup()

		// send a message
		err := queue.Send(context.Background(), model.Message{
			"foo": "bar",
		})
		is.NoErr(err)

		// receive our message
		m, err := queue.Receive(context.Background())
		is.NoErr(err)
		is.Equal(model.Message{"foo": "bar"}, *m)

		// check if last message is out of our queue
		m, _ = queue.Receive(context.Background())
		is.NoErr(err)
		is.Equal(nil, m)

	})

}
