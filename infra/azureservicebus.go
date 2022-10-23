package infra

import (
	"canvas/messaging"
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus"
)

func CreateAzureServiceBusResource(subscriptionID string, resourceGroupName string) (*messaging.NewQueueOptions, error, func()) {
	ctx := context.Background()

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err, nil
	}

	nsclient, err := armservicebus.NewNamespacesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err, nil
	}

	poller, err := nsclient.BeginCreateOrUpdate(ctx, resourceGroupName, "test-asb-canvas-ns", armservicebus.SBNamespace{
		Location: to.Ptr("East US"),
	}, nil)
	if err != nil {
		return nil, err, nil
	}

	resNs, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err, nil
	}

	// keep polling until resource is created

	client, err := armservicebus.NewQueuesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err, nil
	}

	resQueue, err := client.CreateOrUpdate(ctx, resourceGroupName, *resNs.Name, "test-asb-canvas-q", armservicebus.SBQueue{}, nil)
	if err != nil {
		return nil, err, nil
	}

	resKeys, err := nsclient.ListKeys(
		ctx,
		resourceGroupName,
		*resNs.Name,
		"RootManageSharedAccessKey",
		nil,
	)
	if err != nil {
		return nil, err, nil
	}

	opts := messaging.NewQueueOptions{
		Namespace:   *resNs.Name,
		Name:        *resQueue.Name,
		KeyName:     *resKeys.KeyName,
		KeyPassword: *resKeys.PrimaryConnectionString,
	}

	cleanup := func() {
		nsclient.BeginDelete(ctx, resourceGroupName, *resNs.Name, nil)
	}

	return &opts, nil, cleanup
}
