package utils

import "testing"

func TestNatsInit(t *testing.T) {
	NatsInit()
	message := []byte("Hello from Go!")
	err := NATSPublish(message)
	if err != nil {
		return
	}
}
