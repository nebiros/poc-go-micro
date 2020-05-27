package azqueue

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/config/cmd"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
)

const (
	defaultNumWorkers                          = 16
	defaultPoisonMessageDequeueThreshold int64 = 4
)

func init() {
	log.Init(log.WithLevel(log.Level(log.TraceLevel)))

	cmd.DefaultBrokers["azqueue"] = NewBroker
}

type azqueueBroker struct {
	options    broker.Options
	serviceURL azqueue.ServiceURL
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

func (b *azqueueBroker) azqueueServiceURL() azqueue.ServiceURL {
	if v := b.options.Context.Value(serviceURLOptionKey{}); v != nil {
		return v.(azqueue.ServiceURL)
	}

	return azqueue.ServiceURL{}
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
	if serviceURL := b.azqueueServiceURL(); (azqueue.ServiceURL{}) != serviceURL {
		b.serviceURL = serviceURL
		return nil
	}

	storageAccountName := b.storageAccountName()
	if storageAccountName == "" {
		return errors.New("storage account name cannot be empty")
	}

	storageAccountKey := b.storageAccountKey()
	if storageAccountKey == "" {
		return errors.New("storage account key cannot be empty")
	}

	credential, err := azqueue.NewSharedKeyCredential(storageAccountName, storageAccountKey)
	if err != nil {
		return err
	}

	primaryURL, err := url.Parse(fmt.Sprintf("https://%s.queue.core.windows.net", storageAccountName))
	if err != nil {
		return err
	}

	// TODO: add AZ pipeline options as broker options
	b.serviceURL = azqueue.NewServiceURL(*primaryURL, azqueue.NewPipeline(credential, azqueue.PipelineOptions{}))

	return nil
}

func (b *azqueueBroker) Disconnect() error {
	panic("implement me")
}

func (b *azqueueBroker) Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	panic("implement me")
}

func (b *azqueueBroker) Subscribe(queueName string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	options := broker.SubscribeOptions{
		AutoAck: true,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	queueURL := b.serviceURL.NewQueueURL(strings.ToLower(strings.TrimSpace(queueName)))
	messagesURL := queueURL.NewMessagesURL()

	subscriber := &subscriber{
		options:     options,
		queueName:   queueName,
		serviceURL:  b.serviceURL,
		queueURL:    queueURL,
		messagesURL: messagesURL,
		exit:        make(chan bool),
	}

	go subscriber.run(h)

	return subscriber, nil
}

func (b *azqueueBroker) String() string {
	return "azqueue"
}

type subscriber struct {
	options     broker.SubscribeOptions
	queueName   string
	serviceURL  azqueue.ServiceURL
	queueURL    azqueue.QueueURL
	messagesURL azqueue.MessagesURL
	exit        chan bool
}

func (s *subscriber) numWorkers() int {
	if v := s.options.Context.Value(numWorkers{}); v != nil {
		return v.(int)
	}

	return -1
}

func (s *subscriber) poisonMessageDequeueThreshold() int64 {
	if v := s.options.Context.Value(poisonMessageDequeueThreshold{}); v != nil {
		return v.(int64)
	}

	return -1
}

// run initialize a service that wishes to process messages
// See: https://pkg.go.dev/github.com/Azure/azure-storage-queue-go/azqueue?tab=doc#example-package
func (s *subscriber) run(h broker.Handler) {
	log.Debugf("AZQueue: subscription started. Queue: %s, URL: %s", s.queueName, s.queueURL.String())

	for {
		select {
		case <-s.exit:
			return

		default:
			_, err := s.queueURL.GetProperties(s.options.Context)
			if err != nil {
				if errorType := err.(azqueue.StorageError).ServiceCode(); errorType == azqueue.ServiceCodeQueueNotFound {
					log.Info("AZQueue: queue not found, creating")

					if _, err := s.queueURL.Create(s.options.Context, azqueue.Metadata{}); err != nil {
						log.Fatalf("AZQueue: error creating queue: %s", err)
					}

					if _, err := s.queueURL.GetProperties(s.options.Context); err != nil {
						log.Fatalf("AZQueue: error parsing URL: %s", err)
					}
				} else {
					log.Fatalf("AZQueue: error getting queue properties: %s", err)
				}
			}

			numWorkers := s.numWorkers()
			if numWorkers == -1 {
				numWorkers = defaultNumWorkers
			}

			poisonMessageDequeueThreshold := s.poisonMessageDequeueThreshold()
			if poisonMessageDequeueThreshold == -1 {
				poisonMessageDequeueThreshold = defaultPoisonMessageDequeueThreshold
			}

			messagesChan := make(chan *azqueue.DequeuedMessage, numWorkers)

			// Allocate go routines
			for w := 0; w < numWorkers; w++ {
				go s.worker(messagesChan, poisonMessageDequeueThreshold, h)
			}

			// Shows the service's infinite loop that dequeues messages and dispatches them in batches for processsing
			for {
				dequeue, err := s.messagesURL.Dequeue(s.options.Context, azqueue.QueueMaxMessagesDequeue, 10*time.Second)
				if err != nil {
					log.Fatal(err)
				}

				if dequeue.NumMessages() == 0 {
					// The queue was empty; sleep a bit and try again
					// Shorter time means higher costs & less latency to dequeue a brokerMessage
					// Higher time means lower costs & more latency to dequeue a brokerMessage
					time.Sleep(time.Second * 10)
				} else {
					// We got some messages, put them in the channel so that many can be processed in parallel:
					// NOTE: The queue does not guarantee FIFO ordering & processing messages in parallel also does
					// not preserve FIFO ordering. So, the "Output:" order below is not guaranteed but usually works.
					for m := int32(0); m < dequeue.NumMessages(); m++ {
						messagesChan <- dequeue.Message(m)
					}
				}

				// This batch of dequeued messages are in the channel, dequeue another batch
				break // NOTE: For this example only, break out of the infinite loop so this example terminates
			}
		}
	}
}

func (s *subscriber) worker(messagesChan <-chan *azqueue.DequeuedMessage, poisonMessageDequeueThreshold int64, h broker.Handler) {
	for {
		message := <-messagesChan

		messageIDURL := s.messagesURL.NewMessageIDURL(message.ID)
		popReceipt := message.PopReceipt

		if message.DequeueCount > poisonMessageDequeueThreshold {
			// This brokerMessage has attempted to be processed too many times; treat it as a poison brokerMessage
			// DO NOT attempt to process this brokerMessage
			// Log this brokerMessage as a poison brokerMessage somewhere (code not shown)
			// Delete this poison brokerMessage from the queue so it will never be dequeued again
			_, err := messageIDURL.Delete(s.options.Context, popReceipt)
			if err != nil {
				log.Fatal(err)
			}

			continue
		}

		s.handleDequeuedMessage(message, messageIDURL, h)
	}
}

func (s *subscriber) handleDequeuedMessage(dequeuedMessage *azqueue.DequeuedMessage, messageIDURL azqueue.MessageIDURL, h broker.Handler) {
	body, err := base64.StdEncoding.DecodeString(dequeuedMessage.Text)
	if err != nil {
		log.Errorf("AZQueue: failed to decode message body: %s", err)
		return
	}

	log.Debugf("AZQueue: received dequeued message: %d bytes", len(body))

	bm := &broker.Message{
		Header: s.makeMessageHeader(dequeuedMessage),
		Body:   body,
	}

	e := &azqueueEvent{
		queueName:       s.queueName,
		context:         s.options.Context,
		messageIDURL:    messageIDURL,
		brokerMessage:   bm,
		dequeuedMessage: dequeuedMessage,
		err:             nil,
	}

	if e.err = h(e); e.err != nil {
		log.Error(e.err)
	}

	if s.options.AutoAck {
		err := e.Ack()
		if err != nil {
			log.Errorf("AZQueue: failed auto-acknowledge of dequeued message: %s", err)
		}
	}
}

func (s *subscriber) makeMessageHeader(dequeuedMessage *azqueue.DequeuedMessage) map[string]string {
	result := make(map[string]string)
	result["InsertionTime"] = dequeuedMessage.InsertionTime.String()
	result["ExpirationTime"] = dequeuedMessage.ExpirationTime.String()
	result["NextVisibleTime"] = dequeuedMessage.NextVisibleTime.String()

	return result
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.options
}

// Topic returns the name of the queue at azure
func (s *subscriber) Topic() string {
	return s.queueName
}

func (s *subscriber) Unsubscribe() error {
	select {
	case <-s.exit:
		return nil

	default:
		close(s.exit)
		return nil
	}
}

type azqueueEvent struct {
	queueName       string
	context         context.Context
	messageIDURL    azqueue.MessageIDURL
	brokerMessage   *broker.Message
	dequeuedMessage *azqueue.DequeuedMessage
	err             error
}

// Topic returns the name of the queue at azure
func (e *azqueueEvent) Topic() string {
	return e.queueName
}

func (e *azqueueEvent) Message() *broker.Message {
	return e.brokerMessage
}

func (e *azqueueEvent) Ack() error {
	popReceipt := e.dequeuedMessage.PopReceipt

	// OPTIONAL: while processing a brokerMessage, you can update the brokerMessage's visibility timeout
	// (to prevent other servers from dequeuing the same brokerMessage simultaneously) and update the
	// brokerMessage's text (to prevent some successfully-completed processing from re-executing the
	// next time this brokerMessage is dequeued):
	update, err := e.messageIDURL.Update(e.context, popReceipt, time.Second*20, "updated msg")
	if err != nil {
		return err
	}

	// Performing any operation on a brokerMessage ID always requires the most recent pop receipt
	popReceipt = update.PopReceipt

	// After processing the brokerMessage, delete it from the queue so it won't be dequeued ever again:
	_, err = e.messageIDURL.Delete(e.context, popReceipt)
	if err != nil {
		return err
	}

	return nil
}

func (e *azqueueEvent) Error() error {
	return e.err
}
