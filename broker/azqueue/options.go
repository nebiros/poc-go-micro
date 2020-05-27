package azqueue

import (
	"context"

	"github.com/Azure/azure-storage-queue-go/azqueue"

	"github.com/micro/go-micro/v2/broker"
)

type (
	storageAccountNameOptionKey   struct{}
	storageAccountKeyOptionKey    struct{}
	serviceURLOptionKey           struct{}
	numWorkers                    struct{}
	poisonMessageDequeueThreshold struct{}
)

func setBrokerOption(k, v interface{}) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, k, v)
	}
}

func setSubscribeOption(k, v interface{}) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, k, v)
	}
}

func StorageAccountName(name string) broker.Option {
	return setBrokerOption(storageAccountNameOptionKey{}, name)
}

func StorageAccountKey(key string) broker.Option {
	return setBrokerOption(storageAccountKeyOptionKey{}, key)
}

func ServiceURL(serviceURL azqueue.ServiceURL) broker.Option {
	return setBrokerOption(serviceURLOptionKey{}, serviceURL)
}

func NumWorkers(number int) broker.SubscribeOption {
	return setSubscribeOption(numWorkers{}, number)
}

// PoisonMessageDequeueThreshold indicates how many times a brokerMessage is attempted to be processed
// before considering it a poison brokerMessage
func PoisonMessageDequeueThreshold(number int64) broker.SubscribeOption {
	return setSubscribeOption(poisonMessageDequeueThreshold{}, number)
}
