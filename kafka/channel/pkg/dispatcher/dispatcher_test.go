package dispatcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"knative.dev/eventing-contrib/kafka/channel/pkg/utils"

	cluster "github.com/bsm/sarama-cluster"

	"github.com/Shopify/sarama"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"
	"knative.dev/eventing/pkg/channel"
	_ "knative.dev/pkg/system/testing"
)

type mockConsumer struct {
	message    chan *sarama.ConsumerMessage
	partitions chan cluster.PartitionConsumer
}

func (c *mockConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return c.message
}

func (c *mockConsumer) Partitions() <-chan cluster.PartitionConsumer {
	return c.partitions
}

func (c *mockConsumer) Close() error {
	return nil
}

func (c *mockConsumer) MarkOffset(msg *sarama.ConsumerMessage, metadata string) {
	return
}

type mockPartitionConsumer struct {
	highWaterMarkOffset int64
	l                   sync.Mutex
	topic               string
	partition           int32
	offset              int64
	messages            chan *sarama.ConsumerMessage
	errors              chan *sarama.ConsumerError
	singleClose         sync.Once
	consumed            bool
}

// AsyncClose implements the AsyncClose method from the sarama.PartitionConsumer interface.
func (pc *mockPartitionConsumer) AsyncClose() {
	pc.singleClose.Do(func() {
		close(pc.messages)
		close(pc.errors)
	})
}

// Close implements the Close method from the sarama.PartitionConsumer interface. It will
// verify whether the partition consumer was actually started.
func (pc *mockPartitionConsumer) Close() error {
	if !pc.consumed {
		return fmt.Errorf("Expectations set on %s/%d, but no partition consumer was started.", pc.topic, pc.partition)
	}
	pc.AsyncClose()
	var (
		closeErr error
		wg       sync.WaitGroup
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var errs = make(sarama.ConsumerErrors, 0)
		for err := range pc.errors {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			closeErr = errs
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range pc.messages {
			// drain
		}
	}()
	wg.Wait()
	return closeErr
}

// Errors implements the Errors method from the sarama.PartitionConsumer interface.
func (pc *mockPartitionConsumer) Errors() <-chan *sarama.ConsumerError {
	return pc.errors
}

// Messages implements the Messages method from the sarama.PartitionConsumer interface.
func (pc *mockPartitionConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return pc.messages
}

func (pc *mockPartitionConsumer) HighWaterMarkOffset() int64 {
	return atomic.LoadInt64(&pc.highWaterMarkOffset) + 1
}

func (pc *mockPartitionConsumer) Topic() string {
	return pc.topic
}

// Partition returns the consumed partition
func (pc *mockPartitionConsumer) Partition() int32 {
	return pc.partition
}

// InitialOffset returns the offset used for creating the PartitionConsumer instance.
// The returned offset can be a literal offset, or OffsetNewest, or OffsetOldest
func (pc *mockPartitionConsumer) InitialOffset() int64 {
	return 0
}

// MarkOffset marks the offset of a message as preocessed.
func (pc *mockPartitionConsumer) MarkOffset(offset int64, metadata string) {
}

// ResetOffset resets the offset to a previously processed message.
func (pc *mockPartitionConsumer) ResetOffset(offset int64, metadata string) {
	pc.offset = 0
}

type mockSaramaCluster struct {
	// closed closes the message channel so that it doesn't block during the test
	closed bool
	// Handle to the latest created consumer, useful to access underlying message chan
	consumerChannel chan *sarama.ConsumerMessage
	// Handle to the latest created partition consumer, useful to access underlying message chan
	partitionConsumerChannel chan cluster.PartitionConsumer
	// createErr will return an error when creating a consumer
	createErr bool
	// consumer mode
	consumerMode cluster.ConsumerMode
}

func (c *mockSaramaCluster) NewConsumer(groupID string, topics []string) (KafkaConsumer, error) {
	if c.createErr {
		return nil, errors.New("error creating consumer")
	}

	var consumer *mockConsumer
	if c.consumerMode != cluster.ConsumerModePartitions {
		consumer = &mockConsumer{
			message: make(chan *sarama.ConsumerMessage),
		}
		if c.closed {
			close(consumer.message)
		}
		c.consumerChannel = consumer.message
	} else {
		consumer = &mockConsumer{
			partitions: make(chan cluster.PartitionConsumer),
		}
		if c.closed {
			close(consumer.partitions)
		}
		c.partitionConsumerChannel = consumer.partitions
	}
	return consumer, nil
}

func (c *mockSaramaCluster) GetConsumerMode() cluster.ConsumerMode {
	return c.consumerMode
}

func TestFromKafkaMessage(t *testing.T) {
	data := []byte("data")
	kafkaMessage := &sarama.ConsumerMessage{
		Headers: []*sarama.RecordHeader{
			{
				Key:   []byte("k1"),
				Value: []byte("v1"),
			},
		},
		Value: data,
	}
	want := &channel.Message{
		Headers: map[string]string{
			"k1": "v1",
		},
		Payload: data,
	}
	got := fromKafkaMessage(kafkaMessage)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected message (-want, +got) = %s", diff)
	}
}

func TestToKafkaMessage(t *testing.T) {
	data := []byte("data")
	channelRef := channel.ChannelReference{
		Name:      "test-channel",
		Namespace: "test-ns",
	}
	msg := &channel.Message{
		Headers: map[string]string{
			"k1": "v1",
		},
		Payload: data,
	}
	want := &sarama.ProducerMessage{
		Topic: "knative-messaging-kafka.test-ns.test-channel",
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("k1"),
				Value: []byte("v1"),
			},
		},
		Value: sarama.ByteEncoder(data),
	}
	got := toKafkaMessage(channelRef, msg, utils.TopicName)
	if diff := cmp.Diff(want, got, cmpopts.IgnoreUnexported(sarama.ProducerMessage{})); diff != "" {
		t.Errorf("unexpected message (-want, +got) = %s", diff)
	}
}

type dispatchTestHandler struct {
	t       *testing.T
	payload []byte
	done    chan bool
}

func (h *dispatchTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.t.Error("Failed to read the request body")
	}
	if diff := cmp.Diff(h.payload, body); diff != "" {
		h.t.Errorf("unexpected body (-want, +got) = %s", diff)
	}
	h.done <- true
}

func TestSubscribe(t *testing.T) {
	sc := &mockSaramaCluster{}
	data := []byte("data")
	d := &KafkaDispatcher{
		kafkaCluster:   sc,
		kafkaConsumers: make(map[channel.ChannelReference]map[subscription]KafkaConsumer),
		dispatcher:     channel.NewMessageDispatcher(zap.NewNop().Sugar()),
		logger:         zap.NewNop(),
		topicFunc:      utils.TopicName,
	}

	testHandler := &dispatchTestHandler{
		t:       t,
		payload: data,
		done:    make(chan bool)}

	server := httptest.NewServer(testHandler)
	defer server.Close()

	channelRef := channel.ChannelReference{
		Name:      "test-channel",
		Namespace: "test-ns",
	}

	subRef := subscription{
		UID:           "test-sub",
		SubscriberURI: server.URL[7:],
	}
	err := d.subscribe(channelRef, subRef)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	defer close(sc.consumerChannel)
	sc.consumerChannel <- &sarama.ConsumerMessage{
		Headers: []*sarama.RecordHeader{
			{
				Key:   []byte("k1"),
				Value: []byte("v1"),
			},
		},
		Value: data,
	}

	<-testHandler.done

}

func TestPartitionConsumer(t *testing.T) {
	sc := &mockSaramaCluster{consumerMode: cluster.ConsumerModePartitions}
	data := []byte("data")
	d := &KafkaDispatcher{
		kafkaCluster:   sc,
		kafkaConsumers: make(map[channel.ChannelReference]map[subscription]KafkaConsumer),
		dispatcher:     channel.NewMessageDispatcher(zap.NewNop().Sugar()),
		logger:         zap.NewNop(),
		topicFunc:      utils.TopicName,
	}
	testHandler := &dispatchTestHandler{
		t:       t,
		payload: data,
		done:    make(chan bool)}
	server := httptest.NewServer(testHandler)
	defer server.Close()
	channelRef := channel.ChannelReference{
		Name:      "test-channel",
		Namespace: "test-ns",
	}
	subRef := subscription{
		UID:           "test-sub",
		SubscriberURI: server.URL[7:],
	}
	err := d.subscribe(channelRef, subRef)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	defer close(sc.partitionConsumerChannel)
	pc := &mockPartitionConsumer{
		topic:     channelRef.Name,
		partition: 1,
		messages:  make(chan *sarama.ConsumerMessage, 1),
	}
	pc.messages <- &sarama.ConsumerMessage{
		Topic:     channelRef.Name,
		Partition: 1,
		Headers: []*sarama.RecordHeader{
			{
				Key:   []byte("k1"),
				Value: []byte("v1"),
			},
		},
		Value: data,
	}
	sc.partitionConsumerChannel <- pc
	<-testHandler.done
}

func TestSubscribeError(t *testing.T) {
	sc := &mockSaramaCluster{
		createErr: true,
	}
	d := &KafkaDispatcher{
		kafkaCluster: sc,
		logger:       zap.NewNop(),
		topicFunc:    utils.TopicName,
	}

	channelRef := channel.ChannelReference{
		Name:      "test-channel",
		Namespace: "test-ns",
	}

	subRef := subscription{
		UID: "test-sub",
	}
	err := d.subscribe(channelRef, subRef)
	if err == nil {
		t.Errorf("Expected error want %s, got %s", "error creating consumer", err)
	}
}

func TestUnsubscribeUnknownSub(t *testing.T) {
	sc := &mockSaramaCluster{
		createErr: true,
	}
	d := &KafkaDispatcher{
		kafkaCluster: sc,
		logger:       zap.NewNop(),
	}

	channelRef := channel.ChannelReference{
		Name:      "test-channel",
		Namespace: "test-ns",
	}

	subRef := subscription{
		UID: "test-sub",
	}
	if err := d.unsubscribe(channelRef, subRef); err != nil {
		t.Errorf("Unsubscribe error: %v", err)
	}
}

func TestKafkaDispatcher_Start(t *testing.T) {
	d := &KafkaDispatcher{}
	err := d.Start(make(chan struct{}))
	if err == nil {
		t.Errorf("Expected error want %s, got %s", "message receiver is not set", err)
	}

	receiver, err := channel.NewMessageReceiver(func(channel channel.ChannelReference, message *channel.Message) error {
		return nil
	}, zap.NewNop().Sugar())
	if err != nil {
		t.Fatalf("Error creating new message receiver. Error:%s", err)
	}
	d.receiver = receiver
	err = d.Start(make(chan struct{}))
	if err == nil {
		t.Errorf("Expected error want %s, got %s", "kafkaAsyncProducer is not set", err)
	}
}

var sortStrings = cmpopts.SortSlices(func(x, y string) bool {
	return x < y
})
