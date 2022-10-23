package integrationtest

import (
	"canvas/infra"
	"canvas/messaging"

	"github.com/maragudk/env"
)

func CreateQueue() (*messaging.Queue, func()) {
	env.MustLoad("../.env-test")

	queueOpts, err, cleanup := infra.CreateAzureServiceBusResource(env.GetStringOrDefault("AZURE_SUBSCRIPTION_ID", "xxx"), env.GetStringOrDefault("AZURE_RESOURCE_GROUP_NAME", "myTestRGWooho"))
	if err != nil {
		panic(err)
	}

	queue := messaging.NewQueue(*queueOpts)

	return queue, cleanup
}
