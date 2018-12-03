package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"os"
	"time"
)

/*
*   参考: https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/custom_store/main.go
 */

// NoOpStore 实现 go-mqtt/Store 接口
// opts.SetStore(myNoOpStore)
type NoOpStore struct {
	// Contain nothing
}

func (store *NoOpStore) Open() {
	// Do nothing
}

func (store *NoOpStore) Put(string, packets.ControlPacket) {
	// Do nothing
}

func (store *NoOpStore) Get(string) packets.ControlPacket {
	// Do nothing
	return nil
}

func (store *NoOpStore) Del(string) {
	// Do nothing
}

func (store *NoOpStore) All() []string {
	return nil
}
func (store *NoOpStore) Close() {
	// Do nothing
}

func (store *NoOpStore) Reset() {
	// Do nothing
}

func main() {
	// myNoOpStore := &NoOpStore{}
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID("custom-store")
	// opts.SetStore(myNoOpStore)

	var handler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", message.Topic())
		fmt.Printf("MSG: %s\n", message.Payload())
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// c.Subscribe("/go-mqtt/sample", 0, handler)
	if token := c.Subscribe("/go-mqtt/sample", 0, handler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("/go-mqtt/sample", 0, false, text)
		token.Wait()
	}

	for i := 1; i < 5; i++ {
		time.Sleep(1 * time.Second)
	}

	time.Sleep(5 * time.Second)
	c.Disconnect(250)
}
