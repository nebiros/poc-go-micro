package azqueue

import (
	"context"

	"github.com/Azure/azure-storage-queue-go/azqueue"

	"github.com/micro/go-micro/v2/broker"
)

type (
	storageAccountNameOptionKey struct{}
	storageAccountKeyOptionKey  struct{}
	storageQueueNameOptionKey   struct{}
	azqueueQueueURLOptionKey    struct{}
)

func StorageAccountName(name string) broker.Option {
	return func(options *broker.Options) {
		options.Context = context.WithValue(options.Context, storageAccountNameOptionKey{}, name)
	}
}

func StorageAccountKey(key string) broker.Option {
	return func(options *broker.Options) {
		options.Context = context.WithValue(options.Context, storageAccountKeyOptionKey{}, key)
	}
}

func StorageQueueName(name string) broker.Option {
	return func(options *broker.Options) {
		options.Context = context.WithValue(options.Context, storageQueueNameOptionKey{}, name)
	}
}

func AZQueueQueueURL(queueURL azqueue.QueueURL) broker.Option {
	return func(options *broker.Options) {
		options.Context = context.WithValue(options.Context, azqueueQueueURLOptionKey{}, queueURL)
	}
}
