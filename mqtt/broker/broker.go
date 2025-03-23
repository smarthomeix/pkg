package broker

import (
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NewBroker initializes an MQTT client with resilience features.
func NewBroker(host, clientID string) *Client {
	opts := mqtt.NewClientOptions().
		AddBroker(host).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetOrderMatters(true) // Ensures ordered message processing

	var subMu sync.Mutex
	subscriptions := make(map[string]mqtt.MessageHandler)

	// Connection lost handler
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	})

	// OnConnect handler ensures re-subscription
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		subMu.Lock()
		// Copy subscriptions first before unlocking to avoid deadlock
		toSubscribe := make(map[string]mqtt.MessageHandler)

		for topic, handler := range subscriptions {
			toSubscribe[topic] = handler
		}

		subMu.Unlock() // 🔓 Unlock here to avoid deadlock

		// Now subscribe to topics safely outside of the lock
		for topic, handler := range toSubscribe {
			if token := client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
				log.Printf("Subscription to %s failed: %v", topic, token.Error())
			} else {
				log.Printf("Subscribed to %s", topic)
			}
		}
	})

	client := mqtt.NewClient(opts)

	// Blocking connection loop with exponential backoff
	baseDelay := time.Second
	maxDelay := 30 * time.Second

	for {
		token := client.Connect()
		if token.Wait() && token.Error() == nil {
			break
		}

		log.Printf("Failed to connect to broker: %v. Retrying in %s...", token.Error(), baseDelay)
		time.Sleep(baseDelay)
		if baseDelay < maxDelay {
			baseDelay *= 2
		}
	}

	// Return wrapped client with subscription tracking
	return &Client{Client: client, subs: subscriptions, mu: &subMu}
}

// Client wraps mqtt.Client to track subscriptions
type Client struct {
	mqtt.Client
	subs map[string]mqtt.MessageHandler
	mu   *sync.Mutex
}

// SubscribeWithTracking ensures subscriptions are re-added on reconnect
func (tc *Client) SubscribeWithTracking(topic string, qos byte, callback mqtt.MessageHandler) {
	// Add subscription first before attempting connection
	tc.mu.Lock()
	tc.subs[topic] = callback
	tc.mu.Unlock()

	// Subscribe outside of lock
	if token := tc.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		log.Printf("Subscription to %s failed: %v", topic, token.Error())
	} else {
		log.Printf("Subscribed to %s", topic)
	}
}

// Disconnect gracefully shuts down the MQTT client
func (tc *Client) Disconnect() {
	tc.Client.Disconnect(250) // 250ms grace period

	log.Println("MQTT client disconnected gracefully")
}
