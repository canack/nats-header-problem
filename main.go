package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	streamName  = "MY_STREAM"
	subjectName = "MY_SUBJECT"
)

func main() {
	// Establish a connection to the NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	// Create a JetStream context
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("Failed to create JetStream context:", err)
	}

	// Create or update a stream with the given subject
	_, err = js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{subjectName},
	})
	if err != nil {
		log.Fatal("Failed to create or update stream:", err)
	}

	// Start producer and consumer as separate goroutines
	go startProducer(js)
	go startConsumer(js)

	// Keep the program running
	select {}
}

// startProducer publishes messages periodically
func startProducer(js jetstream.JetStream) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	messageCount := 0

	for range ticker.C {
		messageCount++
		msg := nats.NewMsg(subjectName)

		// Add a header to the first three messages
		if messageCount <= 3 {
			msg.Header.Add("Publish-Count", fmt.Sprintf("%d", messageCount))
		}

		msg.Data = []byte(fmt.Sprintf("Message #%d", messageCount))

		// Publish the message
		if _, err := js.PublishMsg(context.Background(), msg); err != nil {
			log.Printf("Error publishing message #%d: %v", messageCount, err)
		}
	}
}

// startConsumer processes incoming messages
func startConsumer(js jetstream.JetStream) {
	subscription, err := js.CreateOrUpdateConsumer(context.Background(), streamName, jetstream.ConsumerConfig{
		Durable:        "test-durable",
		FilterSubjects: []string{subjectName},
		MaxDeliver:     2,
		AckWait:        time.Second,
	})
	if err != nil {
		log.Fatal("Failed to create or update consumer:", err)
	}

	// Consume messages and attempt to set a header, expecting a panic on messages without headers
	_, err = subscription.Consume(func(msg jetstream.Msg) {
		fmt.Printf("Received message: %s, Publish Count: %s\n", msg.Data(), msg.Headers().Get("Publish-Count"))

		// This line will cause a panic if the message has no initial headers
		msg.Headers().Set("No-Panic", "true")

		msg.Ack()
	})
	if err != nil {
		log.Printf("Error consuming messages: %v", err)
	}
}
