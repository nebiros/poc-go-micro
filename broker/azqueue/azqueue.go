package azqueue

import (
	"context"
	"net/url"

	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/config/cmd"
)

type azqueueBroker struct {
	queueURL azqueue.QueueURL
	options  broker.Options
}

func (b *azqueueBroker) storageAccountName() string {
	if v := b.options.Context.Value(storageAccountNameOptionKey{}); v != nil {
		return v.(string)
	}

	return ""
}

func (b *azqueueBroker) storageAccountKey() string {
	if v := b.options.Context.Value(storageAccountKeyOptionKey{}); v != nil {
		return v.(string)
	}

	return ""
}

func (b *azqueueBroker) storageQueueName() string {
	if v := b.options.Context.Value(storageQueueNameOptionKey{}); v != nil {
		return v.(string)
	}

	return ""
}

func (b *azqueueBroker) azqueueQueueURL() azqueue.QueueURL {
	if v := b.options.Context.Value(azqueueQueueURLOptionKey{}); v != nil {
		return v.(azqueue.QueueURL)
	}

	return azqueue.QueueURL{}
}

func (b *azqueueBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&b.options)
	}

	return nil
}

func (b *azqueueBroker) Options() broker.Options {
	return b.options
}

func (b *azqueueBroker) Address() string {
	return ""
}

func (b *azqueueBroker) Connect() error {
	if queueURL := b.azqueueQueueURL(); (azqueue.QueueURL{}) != queueURL {
		b.queueURL = queueURL
		return nil
	}

	storageAccountName := b.storageAccountName()
	storageAccountKey := b.storageAccountKey()
	storageQueueName := b.storageQueueName()

	credential, err := azqueue.NewSharedKeyCredential(storageAccountName, storageAccountKey)
	if err != nil {
		return err
	}

	qURL, err := url.Parse(fmt.Sprintf("https://%s.queue.core.windows.net/%s", storageAccountName, storageQueueName))
	if err != nil {
		return err
	}

	b.queueURL = azqueue.NewQueueURL(*qURL, azqueue.NewPipeline(credential, azqueue.PipelineOptions{}))

	return nil
}

func (b *azqueueBroker) Disconnect() error {
	panic("implement me")
}

func (b *azqueueBroker) Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	panic("implement me")
}

func (b *azqueueBroker) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	panic("implement me")
}

func (b *azqueueBroker) String() string {
	panic("implement me")
}

func init() {
	cmd.DefaultBrokers["azqueue"] = NewBroker
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &azqueueBroker{
		options: options,
	}
}
